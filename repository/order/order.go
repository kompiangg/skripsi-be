package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/type/model"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	InsertDataToShardDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	DeleteAllDataFromOneDB(ctx context.Context, tx *sqlx.Tx) error
	FindAllOnShardDB(ctx context.Context, dbIndex int) ([]model.Order, error)
}

type Config struct {
	ShardingDatabase config.ShardingDatabase
}

type repository struct {
	config  Config
	shardDB []*sqlx.DB
}

func New(
	config Config,
	shardDB []*sqlx.DB,
) Repository {
	return repository{
		config:  config,
		shardDB: shardDB,
	}
}
