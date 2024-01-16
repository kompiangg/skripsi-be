package account

import (
	"context"
	"skripsi-be/type/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	InsertNewAccount(ctx context.Context, account model.Account) error
	FindAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error)
	FindAccountByUsername(ctx context.Context, username string) (model.Account, error)
}

type Config struct{}

type repository struct {
	config    Config
	generalDB *sqlx.DB
}

func New(config Config, generalDB *sqlx.DB) Repository {
	return repository{
		config:    config,
		generalDB: generalDB,
	}
}
