package pkg

import (
	"github.com/jmoiron/sqlx"
	"loan-ranger/internal/pkg/config"
)

type Options struct {
	Config config.Config
	DB     *sqlx.DB
}
