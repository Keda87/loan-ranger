package repository

import (
	"loan-ranger/internal/repository/project"
	"loan-ranger/internal/repository/project_history"
)

type Container struct {
	Project        project.Repository
	ProjectHistory project_history.Repository
}
