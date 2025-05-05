package project

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loan-ranger/internal/model/payload"
	pkgerr "loan-ranger/internal/pkg/error"
	"net/http"
	"path/filepath"
)

func (h Handler) DisburseProject(c echo.Context) error {
	var (
		ctx          = c.Request().Context()
		projectIDStr = c.Param("project_id")
		body         payload.DisburseProject
	)

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		return pkgerr.Err400("invalid project id format")
	}

	file, err := c.FormFile("signed_agreement_document")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	body.ProjectID = projectID
	body.SignedAgreementDocument = &src
	body.DocumentExtension = filepath.Ext(file.Filename)
	if err = c.Bind(&body); err != nil {
		return err
	}

	if err = c.Validate(&body); err != nil {
		return err
	}

	if err = h.Project.DisburseProject(ctx, body); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
