package project

import (
	"github.com/labstack/echo/v4"
	"loan-ranger/internal/model/payload"
	"net/http"
)

func (h Handler) CreateProject(c echo.Context) error {
	var (
		ctx  = c.Request().Context()
		body payload.CreateProjectPayload
	)

	if err := c.Bind(&body); err != nil {
		return err
	}

	if err := c.Validate(&body); err != nil {
		return err
	}

	resp, err := h.Project.CreateProject(ctx, body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, payload.ResponseData{Data: resp})
}
