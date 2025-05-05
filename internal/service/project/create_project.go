package project

import (
	"context"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
	pkgerr "loan-ranger/internal/pkg/error"
	"loan-ranger/internal/pkg/types"
	"log/slog"
)

func (s Service) CreateProject(ctx context.Context, data payload.CreateProjectPayload) (db.ProjectDetail, error) {
	var detail db.ProjectDetail

	err := dbase.BeginTransaction(ctx, s.DB, func(ctx context.Context) error {
		projectID, err := s.Project.Insert(ctx, data)
		if err != nil {
			slog.Error("error on insert project", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		history := db.CreateProjectHistory{
			ProjectID: projectID,
			Status:    types.StatusProposed,
			PICMail:   data.ActorMail,
			PICName:   data.ActorName,
		}
		if err = s.ProjectHistory.Insert(ctx, history); err != nil {
			slog.Error("error on insert project history", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		detail, err = s.Project.GetByID(ctx, projectID)
		if err != nil {
			slog.Error("error on get detail project", slog.String("err", err.Error()))
			return pkgerr.Err500("internal server error")
		}

		return nil
	})
	if err != nil {
		return detail, err
	}

	return detail, nil
}
