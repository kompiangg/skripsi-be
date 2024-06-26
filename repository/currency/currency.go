package currency

import (
	"context"
	"skripsi-be/type/model"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByBaseAndQuote(ctx context.Context, base string, quote string) (model.Currency, error)
}

type Config struct {
}

type repository struct {
	config    Config
	generalDB *sqlx.DB
}

func New(
	config Config,
	generalDB *sqlx.DB,
) Repository {
	return repository{
		config:    config,
		generalDB: generalDB,
	}
}
