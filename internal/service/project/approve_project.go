package project

import (
	"context"
	"database/sql"
	"errors"
	"github.com/guregu/null"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
	pkgerr "loan-ranger/internal/pkg/error"
	"loan-ranger/internal/pkg/types"
	"log/slog"
	"time"
)

func (s Service) ApproveProject(ctx context.Context, data payload.ApproveProject) error {

	err := dbase.BeginTransaction(ctx, s.DB, func(ctx context.Context) error {
		project, err := s.Project.GetByID(ctx, data.ProjectID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return pkgerr.Err404("project not found or invalid project id")
			}
			slog.Error("error on get project detail", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		if project.CurrentStatus != types.StatusProposed {
			return pkgerr.Err422("only proposed project allowed to approve")
		}

		upd := db.UpdateProject{
			CurrentStatus:  project.CurrentStatus.Next(),
			CurrentPICMail: data.ActorMail,
			CurrentPICName: data.ActorName,
			LastUpdatedAt:  project.UpdatedAt,
			ApprovedAt:     null.NewTime(time.Now().UTC(), true),
		}
		if err = s.Project.UpdateByID(ctx, upd, data.ProjectID); err != nil {
			slog.Error("error on update project detail", slog.String("err", err.Error()))
			return err
		}

		nextHistory := db.CreateProjectHistory{
			ProjectID: data.ProjectID,
			Status:    upd.CurrentStatus,
			PICName:   upd.CurrentPICName,
			PICMail:   upd.CurrentPICMail,
			Extra: map[string]string{
				"field_visit_pic_id":    data.FieldVisitPICID.String(),
				"field_visit_pic_name":  data.FieldVisitPICName,
				"field_visit_pic_mail":  data.FieldVisitPICMail,
				"field_visit_proof_url": data.FieldVisitProofURL,
			},
		}
		if err = s.ProjectHistory.Insert(ctx, nextHistory); err != nil {
			slog.Error("error on insert next history", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
