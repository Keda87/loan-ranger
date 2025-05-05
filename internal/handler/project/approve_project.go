package project

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loan-ranger/internal/model/payload"
	pkgerr "loan-ranger/internal/pkg/error"
	"net/http"
)

func (h Handler) ApproveProject(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		projectIDStr = c.Param("project_id")
		body         payload.ApproveProject
	)

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return pkgerr.Err400("invalid project id format")
	}

	body.ProjectID = projectID
	if err = c.Bind(&body); err != nil {
		return err
	}

	if err = c.Validate(&body); err != nil {
		return err
	}

	if err = h.Project.ApproveProject(ctx, body); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
