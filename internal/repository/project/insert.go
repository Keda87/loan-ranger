package project

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"loan-ranger/internal/model/payload"
	"loan-ranger/internal/pkg/dbase"
	"loan-ranger/internal/pkg/types"
)

func (r Repository) Insert(ctx context.Context, data payload.CreateProjectPayload) (insertedID uuid.UUID, err error) {

	query, args := sq.Insert("projects").SetMap(sq.Eq{
		"name":                  data.Name,
		"borrower_id":           data.BorrowerID,
		"borrower_name":         data.BorrowerName,
		"borrower_mail":         data.BorrowerMail,
		"borrower_rate":         data.BorrowerRate,
		"loan_principal_amount": data.LoanPrincipalAmount,
		"roi_rate":              data.ROIRate,
		"current_status":        types.StatusProposed,
		"current_pic_name":      data.ActorName,
		"current_pic_mail":      data.ActorMail,
	}).PlaceholderFormat(sq.Dollar).Suffix("RETURNING id").MustSql()

	var conn dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		conn = tx
	}

	if err = conn.GetContext(ctx, &insertedID, query, args...); err != nil {
		return uuid.Nil, err
	}

	return insertedID, nil
}
