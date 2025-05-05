package project

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"loan-ranger/internal/model/db"
	pkgerr "loan-ranger/internal/pkg/error"
	"log/slog"
)

func (s Service) DetailProject(ctx context.Context, id uuid.UUID) (db.ProjectDetail, error) {
	project, err := s.Project.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.ProjectDetail{}, pkgerr.Err404("project not found or invalid project id")
		}
		slog.Error("error on get project detail", slog.String("err", err.Error()))
		return db.ProjectDetail{}, pkgerr.Err500("internal server error")
	}

	agreementSignedURL, _ := s.Bucket.GetSignURL(ctx, project.BorrowerAgreementURL.String)
	project.BorrowerAgreementURL.SetValid(agreementSignedURL)

	return project, nil
}
