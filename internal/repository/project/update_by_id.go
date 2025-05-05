package project

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/pkg/dbase"
	pkgerr "loan-ranger/internal/pkg/error"
)

func (r Repository) UpdateByID(ctx context.Context, data db.UpdateProject, id uuid.UUID) error {
	params := sq.Eq{
		"current_status": data.CurrentStatus,
		"updated_at":     sq.Expr("NOW()"),
	}

	if data.CurrentPICName.Valid && data.CurrentPICName.String != "" {
		params["current_pic_name"] = data.CurrentPICName
	}

	if data.CurrentPICMail.Valid && data.CurrentPICMail.String != "" {
		params["current_pic_mail"] = data.CurrentPICMail
	}

	if data.BorrowerAgreementURL.Valid && data.BorrowerAgreementURL.String != "" {
		params["borrower_agreement_url"] = data.BorrowerAgreementURL
	}

	if data.TotalInvestedAmount.Valid && data.TotalInvestedAmount.Float64 > 0 {
		params["total_invested_amount"] = data.TotalInvestedAmount
	}

	if data.ApprovedAt.Valid && !data.ApprovedAt.IsZero() {
		params["approved_at"] = data.ApprovedAt
	}

	if data.DisbursedAt.Valid && !data.DisbursedAt.IsZero() {
		params["disbursed_at"] = data.DisbursedAt
	}

	query, args := sq.
		Update("projects").
		SetMap(params).
		Where("id = ? AND updated_at = ?", id, data.LastUpdatedAt). // optimistic locking mechanism.
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var conn dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		conn = tx
	}

	row, err := conn.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	updated, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if updated == 0 {
		return pkgerr.Err422("please try again")
	}

	return nil
}
