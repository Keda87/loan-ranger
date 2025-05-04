package project

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) CreateProject(c echo.Context) error {
	return c.NoContent(http.StatusCreated)
}
