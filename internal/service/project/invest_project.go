package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guregu/null"

	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
	pkgerr "loan-ranger/internal/pkg/error"
	"loan-ranger/internal/pkg/types"
)

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

			// TODO:
			// - insert history.                                        [V]
			// - publish background event generate pdf each lenders.    [ ]
			// - generate pdf for all lenders.                          [ ]
			// - send email to lenders.                                 [ ]
			// - update project_investments.investment_agreement_url    [ ]
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

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
