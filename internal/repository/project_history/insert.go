package project_history

import (
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"loan-ranger/internal/model/db"
	"loan-ranger/internal/pkg/dbase"
)

func (r Repository) Insert(ctx context.Context, data db.CreateProjectHistory) error {
	params := sq.Eq{
		"project_id": data.ProjectID,
		"status":     data.Status,
		"pic_name":   data.PICName,
		"pic_mail":   data.PICMail,
	}
	if data.Extra != nil {
		jsonBytes, _ := json.Marshal(data.Extra)
		params["extra"] = jsonBytes
	}

	query, args := sq.
		Insert("project_histories").
		SetMap(params).
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
