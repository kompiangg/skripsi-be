package admin

import (
	"context"
	"skripsi-be/type/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindAdminByID(ctx context.Context, id uuid.UUID) (model.Admin, error)
	FindAdminAccountByID(ctx context.Context, id uuid.UUID) (model.Admin, error)
	InsertNewAdmin(ctx context.Context, admin model.Admin) error
	InsertNewCashier(ctx context.Context, cashier model.Cashier) error
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
