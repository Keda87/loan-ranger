package project

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/guregu/null"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
	pkgerr "loan-ranger/internal/pkg/error"
	"loan-ranger/internal/pkg/types"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

func (s Service) DisburseProject(ctx context.Context, data payload.DisburseProject) error {

	err := dbase.BeginTransaction(ctx, s.DB, func(ctx context.Context) error {
		project, err := s.Project.GetByID(ctx, data.ProjectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return pkgerr.Err404("project not found or invalid project id")
			}
			slog.Error("error on get project detail", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		if project.CurrentStatus != types.StatusInvested {
			return pkgerr.Err422("disbursement is not allowed for non invested project")
		}

		fileName := filepath.Join("borrower", "agreement", fmt.Sprintf("project-%s%s", project.ID.String(), data.DocumentExtension))
		if err != nil {
			slog.Error("error on generate pdf", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		bucketURL, err := s.Bucket.Upload(ctx, fileName, *data.SignedAgreementDocument)
		if err != nil {
			slog.Error("error on upload to bucket", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		upd := db.UpdateProject{
			CurrentStatus:        project.CurrentStatus.Next(),
			LastUpdatedAt:        project.UpdatedAt,
			BorrowerAgreementURL: null.StringFrom(bucketURL),
			DisbursedAt:          null.TimeFrom(time.Now().UTC()),
			CurrentPICName:       null.StringFrom(data.ActorName),
			CurrentPICMail:       null.StringFrom(data.ActorMail),
		}
		if err = s.Project.UpdateByID(ctx, upd, project.ID); err != nil {
			slog.Error("error on update project", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		hist := db.CreateProjectHistory{
			ProjectID: project.ID,
			Status:    upd.CurrentStatus,
			PICName:   upd.CurrentPICName.String,
			PICMail:   upd.CurrentPICMail.String,
			Extra: map[string]string{
				"field_visit_pic_id":   data.FieldVisitPICID.String(),
				"field_visit_pic_name": data.FieldVisitPICName,
				"field_visit_pic_mail": data.FieldVisitPICMail,
			},
		}
		if err = s.ProjectHistory.Insert(ctx, hist); err != nil {
			slog.Error("error on insert disbursed project history", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func getFileExtensionFromBuffer(buf []byte) (string, error) {
	// Detect content type (only looks at first 512 bytes)
	contentType := http.DetectContentType(buf)

	// Get file extensions associated with the MIME type
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || len(exts) == 0 {
		return "", fmt.Errorf("no extension found for MIME type: %s", contentType)
	}

	return exts[0], nil // return the first extension (e.g. ".jpg")
}
