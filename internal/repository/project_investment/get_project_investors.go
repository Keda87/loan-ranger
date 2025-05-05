package project_investment

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/pkg/dbase"
)

func (r Repository) GetProjectInvestors(ctx context.Context, projectID uuid.UUID) ([]db.ProjectInvestorItem, error) {
	items := make([]db.ProjectInvestorItem, 0)

	columns := []string{
		"id", "project_id", "investor_id",
		"investor_name", "investor_mail", "investment_amount",
	}
	query, args := sq.
		Select(columns...).
		From("project_investments").
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var conn dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		conn = tx
	}

	if err := conn.SelectContext(ctx, &items, query, args...); err != nil {
		return items, err
	}

	return items, nil
}
