package cashier

import (
	"context"
	"skripsi-be/type/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (model.Cashier, error)
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
