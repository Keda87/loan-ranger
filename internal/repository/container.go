package repository

import (
	"loan-ranger/internal/repository/project"
	"loan-ranger/internal/repository/project_history"
	"loan-ranger/internal/repository/project_investment"
)

type Container struct {
	Project           project.Repository
	ProjectHistory    project_history.Repository
	ProjectInvestment project_investment.Repository
}
