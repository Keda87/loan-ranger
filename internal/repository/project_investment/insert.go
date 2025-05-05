package project_investment

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/pkg/dbase"
)

func (r Repository) Insert(ctx context.Context, data db.CreateProjectInvestment) error {
	params := sq.Eq{
		"project_id":        data.ProjectID,
		"investor_id":       data.InvestorID,
		"investor_name":     data.InvestorName,
		"investor_mail":     data.InvestorMail,
		"investment_amount": data.InvestmentAmount,
	}
	if data.InvestmentAgreementURL.Valid && data.InvestmentAgreementURL.String != "" {
		params["investment_agreement_url"] = data.InvestmentAgreementURL
	}

	query, args := sq.
		Insert("project_investments").
		SetMap(params).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var conn dbase.SQLQueryExec = r.DB
	if tx := dbase.GetTxFromContext(ctx); tx != nil {
		conn = tx
	}

	if _, err := conn.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
