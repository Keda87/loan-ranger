package project

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/pkg/dbase"
)

func (r Repository) GetByID(ctx context.Context, id uuid.UUID) (res db.ProjectDetail, err error) {
	columns := []string{
		"id", "name", "borrower_id", "borrower_name", "borrower_mail", "borrower_agreement_url",
		"current_status", "current_pic_name", "current_pic_mail", "loan_principal_amount", "total_invested_amount",
		"borrower_rate", "roi_rate", "approved_at", "disbursed_at", "created_at", "updated_at",
	}
	query, args := sq.
		Select(columns...).
		From("projects").
		Where("id = ? AND deleted_at IS NULL", id).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var conn dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		conn = tx
	}

	if err = conn.GetContext(ctx, &res, query, args...); err != nil {
		return res, err
	}

	return res, nil
}
