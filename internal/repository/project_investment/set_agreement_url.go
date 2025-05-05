package project_investment

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"loan-ranger/internal/pkg/dbase"
)

func (r Repository) SetAgreementURL(ctx context.Context, agreementPath string, id uuid.UUID) error {
	query, args := sq.
		Update("project_investments").
		SetMap(sq.Eq{
			"investment_agreement_url": agreementPath,
			"updated_at":               sq.Expr("NOW()"),
		}).
		Where("id = ? AND deleted_at IS NULL", id).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var conn dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		conn = tx
	}

	_, err := conn.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
