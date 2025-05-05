package project

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loan-ranger/internal/model/payload"
	pkgerr "loan-ranger/internal/pkg/error"
	"net/http"
)

func (h Handler) InvestProject(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		projectIDStr = c.Param("project_id")
		body         payload.InvestProject
	)

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return pkgerr.Err422("invalid project id format")
	}

	if err = c.Bind(&body); err != nil {
		return err
	}

	body.ProjectID = projectID
	if err = c.Validate(&body); err != nil {
		return err
	}

	if err = h.Project.InvestProject(ctx, body); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
