package project

import (
	"loan-ranger/internal/pkg"
	"loan-ranger/internal/repository"
)

type Service struct {
	*pkg.Options
	*repository.Container
}
