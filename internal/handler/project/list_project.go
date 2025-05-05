package project

import (
	"github.com/labstack/echo/v4"
	"loan-ranger/internal/model/payload"
	"net/http"
)

func (h Handler) ListProject(c echo.Context) error {
	var (
		ctx   = c.Request().Context()
		param payload.ProjectPaginationFilter
	)

	if err := c.Bind(&param); err != nil {
		return err
	}
	param.Normalize()

	items, count, err := h.Project.PaginateProject(ctx, param)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, param.Paginate(items, count))
}
