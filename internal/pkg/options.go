package pkg

import (
	"github.com/jmoiron/sqlx"
	"loan-ranger/internal/pkg/config"
	"loan-ranger/internal/pkg/files"
	"loan-ranger/internal/pkg/mailer"
)

type Options struct {
	Config      config.Config
	DB          *sqlx.DB
	Bucket      files.BucketInterface
	EmailClient mailer.EmailClientInterface
}
