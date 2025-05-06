package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/guregu/null"
	"golang.org/x/sync/errgroup"

	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
	pkgerr "loan-ranger/internal/pkg/error"
	"loan-ranger/internal/pkg/files"
	"loan-ranger/internal/pkg/mailer"
	"loan-ranger/internal/pkg/types"
)

const goroutineLimit = 50

func (s Service) InvestProject(ctx context.Context, data payload.InvestProject) error {
	var totalInvestment float64

	err := dbase.BeginTransaction(ctx, s.DB, func(ctx context.Context) error {
		project, err := s.Project.GetByID(ctx, data.ProjectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return pkgerr.Err404("project not found or invalid project id")
			}
			slog.Error("error on get project detail", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		if data.InvestmentAmount > project.LoanPrincipalAmount {
			return pkgerr.Err422("out of quota for this project")
		}

		if project.CurrentStatus == types.StatusInvested {
			return pkgerr.Err422("this project is fully funded")
		}

		if project.CurrentStatus != types.StatusApproved {
			return pkgerr.Err422("investment only for the approved project")
		}

		totalInvestment = project.TotalInvestedAmount + data.InvestmentAmount
		if totalInvestment > project.LoanPrincipalAmount {
			remainingAmount := project.LoanPrincipalAmount - project.TotalInvestedAmount
			msg := fmt.Sprintf("out of quota, only %.2f remaining for this project", remainingAmount)
			return pkgerr.Err422(msg)
		}

		upd := db.UpdateProject{
			CurrentStatus:       project.CurrentStatus,
			LastUpdatedAt:       project.UpdatedAt,
			TotalInvestedAmount: null.NewFloat(totalInvestment, true),
		}
		if totalInvestment == project.LoanPrincipalAmount {
			upd.CurrentStatus = project.CurrentStatus.Next()
		}
		if err = s.Project.UpdateByID(ctx, upd, data.ProjectID); err != nil {
			return err
		}

		invst := db.CreateProjectInvestment{
			ProjectID:        data.ProjectID,
			InvestorID:       data.InvestorID,
			InvestorName:     data.InvestorName,
			InvestorMail:     data.InvestorMail,
			InvestmentAmount: data.InvestmentAmount,
		}
		if err = s.ProjectInvestment.Insert(ctx, invst); err != nil {
			slog.Error("error on insert investment", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		if upd.CurrentStatus == types.StatusInvested {
			invested := db.CreateProjectHistory{
				ProjectID: project.ID,
				Status:    upd.CurrentStatus,
				PICName:   "", // project is auto invested if fully funded without trigger from the PIC.
				PICMail:   "", // project is auto invested if fully funded without trigger from the PIC.
			}
			if err = s.ProjectHistory.Insert(ctx, invested); err != nil {
				slog.Error("error on create history invested", slog.String("err", err.Error()))
				return pkgerr.Err500("internal server error")
			}

			investors, err := s.ProjectInvestment.GetProjectInvestors(ctx, project.ID)
			if err != nil {
				slog.Error("error on get project investors", slog.String("err", err.Error()))
				return pkgerr.Err500("internal server error")
			}

			eg, ctx := errgroup.WithContext(ctx)
			eg.SetLimit(goroutineLimit)
			for _, investor := range investors {
				i := investor
				eg.TryGo(func() error {
					lenderAgreementURL, err := s.generateLenderPDF(ctx, project, i)
					if err != nil {
						return err
					}

					if err = s.sendLendingAgreement(ctx, project, i, lenderAgreementURL); err != nil {
						return err
					}

					return nil
				})
			}

			if err := eg.Wait(); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s Service) generateLenderPDF(
	ctx context.Context,
	project db.ProjectDetail,
	investor db.ProjectInvestorItem,
) (signedURL string, err error) {
	content := `
Surat Perjanjian

Berikut adalah surat perjanjian investasi pada project:
- Nama Lender: %s
- Nama Project: %s
- Nilai Investasi Anda: %.2f
- Bunga Bagi Hasil: %.2f persen
`
	content = fmt.Sprintf(content, investor.InvestorName, project.Name, investor.InvestmentAmount, project.ROIRate)
	pdfBuff, err := files.GeneratePDFBuffer(content)
	if err != nil {
		return "", err
	}

	agreementPathKey := filepath.Join("lenders", "projects", project.ID.String(), "investors", investor.InvestorID.String(), "agreement.pdf")
	_, err = s.Bucket.Upload(ctx, agreementPathKey, pdfBuff)
	if err != nil {
		return "", err
	}

	if err = s.ProjectInvestment.SetAgreementURL(ctx, agreementPathKey, investor.ID); err != nil {
		return "", err
	}

	signedURL, err = s.Bucket.GetSignURL(ctx, agreementPathKey)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}

func (s Service) sendLendingAgreement(
	ctx context.Context,
	project db.ProjectDetail,
	investor db.ProjectInvestorItem,
	agreementURL string,
) error {

	var sb strings.Builder
	sb.WriteString("<p>Surat Perjanjian Investasi<p><br>")
	sb.WriteString(fmt.Sprintf("<p>Name Proyek: %s</p></br>", project.Name))
	sb.WriteString(fmt.Sprintf("<p>Name Investor: %s</p></br>", investor.InvestorName))
	sb.WriteString(fmt.Sprintf("<p>Jumlah Investasi: %.2f</p></br>", project.TotalInvestedAmount))
	sb.WriteString(fmt.Sprintf("<p>Bunga: %.2f</p></br>", project.ROIRate))
	sb.WriteString(fmt.Sprintf("<p>Dokumen: %s</p></br>", agreementURL))

	return s.EmailClient.SendEmail(ctx, mailer.SendEmail{
		Subject: "Surat Perjanjian Investasi",
		ToEmail: investor.InvestorMail,
		Body:    sb.String(),
	})
}
