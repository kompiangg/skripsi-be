package payment_types

import (
	"context"
	"skripsi-be/type/model"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByLikeOneOfAllColumn(ctx context.Context, req string) ([]model.PaymentType, error)
}

type Config struct{}

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
