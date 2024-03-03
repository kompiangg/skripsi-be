package item

import (
	"context"
	"skripsi-be/type/model"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (model.Item, error)
	FindLikeNameOrID(ctx context.Context, nameOrID string) ([]model.Item, error)
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
