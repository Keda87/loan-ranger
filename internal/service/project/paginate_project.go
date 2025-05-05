package project

import (
	"context"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
)

func (s Service) PaginateProject(ctx context.Context, param payload.ProjectPaginationFilter) (items []db.ProjectItem, total int64, err error) {
	return s.Project.Paginate(ctx, param)
}
