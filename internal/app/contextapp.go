package app

import (
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"loan-ranger/internal/pkg/config"
)

type ContextApp struct {
	Config config.Config
}

func (c ContextApp) GetDB() *sqlx.DB {
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Config.DBUser,
		c.Config.DBPass,
		c.Config.DBHost,
		c.Config.DBPort,
		c.Config.DBName,
	)

	db, err := sqlx.Open("pgx", dbDSN)
	if err != nil {
		slog.Error("error on open connection sqlx",
			slog.String("err", err.Error()),
			slog.String("dsn", dbDSN),
		)
		panic(err)
	}
	db.SetMaxIdleConns(c.Config.DBMaxIdle)
	db.SetMaxOpenConns(c.Config.DBMaxOpen)
	db.SetConnMaxLifetime(time.Minute * 5)

	if err = db.Ping(); err != nil {
		slog.Error("error on ping", slog.String("err", err.Error()))
	}

	return db
}
