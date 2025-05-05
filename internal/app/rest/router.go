package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"loan-ranger/internal/handler/project"
	"loan-ranger/internal/pkg"
	"loan-ranger/internal/pkg/custom"
	"loan-ranger/internal/repository"
	projectrepo "loan-ranger/internal/repository/project"
	projecthistoryrepo "loan-ranger/internal/repository/project_history"
	projectinvestmentrepo "loan-ranger/internal/repository/project_investment"
	"loan-ranger/internal/service"
	projectsvc "loan-ranger/internal/service/project"
)

func (s *Server) initRouter() {
	s.E.Validator = custom.NewValidator()
	s.E.HTTPErrorHandler = custom.NewErrorHandler

	s.E.Use(middleware.RequestID())
	s.E.Use(middleware.Logger())
	s.E.Use(middleware.Recover())
	s.E.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
	}))

	s.E.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "UP!"})
	})

	var (
		repositoryContainer = repository.Container{
			Project:           projectrepo.Repository{Options: &pkg.Options{Config: s.opt.Config, DB: s.opt.DB}},
			ProjectHistory:    projecthistoryrepo.Repository{Options: &pkg.Options{Config: s.opt.Config, DB: s.opt.DB}},
			ProjectInvestment: projectinvestmentrepo.Repository{Options: &pkg.Options{Config: s.opt.Config, DB: s.opt.DB}},
		}

		serviceContainer = service.Container{
			Project: projectsvc.Service{Options: s.opt, Container: &repositoryContainer},
		}

		projectHandler = project.Handler{Container: &serviceContainer}
	)

	v1 := s.E.Group("/v1")

	projectroute := v1.Group("/projects")
	projectroute.POST("", projectHandler.CreateProject)
	projectroute.GET("", projectHandler.ListProject)
	projectroute.GET("/:project_id", projectHandler.DetailProject)
	projectroute.PATCH("/:project_id/approval", projectHandler.ApproveProject)
	projectroute.PATCH("/:project_id/disbursement", projectHandler.DisburseProject)
	projectroute.POST("/:project_id/investment", projectHandler.InvestProject)

}
