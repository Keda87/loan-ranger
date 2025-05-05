package project

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
)

func (r Repository) Paginate(ctx context.Context, param payload.ProjectPaginationFilter) (items []db.ProjectItem, total int64, err error) {
	items = make([]db.ProjectItem, 0)

	query, args := r.
		queryAll(param, []string{"id", "name", "current_status", "loan_principal_amount", "total_invested_amount", "roi_rate", "created_at"}...).
		Limit(uint64(param.Limit)).
		Offset(uint64((param.Page - 1) * param.Limit)).
		OrderBy("created_at DESC").
		MustSql()

	queryCount, argsCount := r.queryAll(param, "COUNT(id) AS total").MustSql()

	var con dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		return items, total, err
	}

	if err = con.SelectContext(ctx, &items, query, args...); err != nil {
		return items, total, err
	}

	if err = con.GetContext(ctx, &total, queryCount, argsCount...); err != nil {
		return items, total, err
	}

	return items, total, nil
}

func (r Repository) queryAll(param payload.ProjectPaginationFilter, columns ...string) sq.SelectBuilder {

	query := sq.
		Select(columns...).
		From("projects").
		PlaceholderFormat(sq.Dollar)

	if param.Search != "" {
		query = query.Where(sq.ILike{"name": param.Search + "%"})
	}

	if param.Status.String() != "" {
		query = query.Where("current_status = ?", param.Status)
	}

	return query
}
