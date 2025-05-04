package project

import (
	"context"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	pkgerr "loan-ranger/internal/pkg/error"
	"loan-ranger/internal/pkg/types"
)

func (s Service) CreateProject(ctx context.Context, data payload.CreateProjectPayload) error {

	projectID, err := s.Project.Insert(ctx, data)
	if err != nil {
		return pkgerr.Err500("internal server error")
	}

	history := db.CreateProjectHistory{
		ProjectID: projectID,
		Status:    types.StatusProposed,
		PICMail:   data.ActorMail,
		PICName:   data.ActorName,
	}
	if err = s.ProjectHistory.Insert(ctx, history); err != nil {
		return pkgerr.Err500("internal server error")
	}

	return nil
}
