package project

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loan-ranger/internal/model/payload"
	pkgerr "loan-ranger/internal/pkg/error"
	"net/http"
)

func (h Handler) DetailProject(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		projectIDStr = c.Param("project_id")
	)

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return pkgerr.Err400("invalid project id format")
	}

	detail, err := h.Project.DetailProject(ctx, projectID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, payload.ResponseData{Data: detail})
}
