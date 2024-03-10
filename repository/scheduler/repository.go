package scheduler

import (
	"context"
	"skripsi-be/type/model"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	FindByName(ctx context.Context, jobName string) (model.Scheduler, error)
	IncrementRunCount(ctx context.Context, tx *sqlx.Tx, jobName string) error
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
