package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/type/model"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	InsertToShardDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	InsertDetailsToShardDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error
	InsertToLongTermDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	InsertDetailsToLongTermDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error
	DeleteAllData(ctx context.Context, tx *sqlx.Tx) error
	DeleteOrderDetails(ctx context.Context, tx *sqlx.Tx) error
	FindAllOnShardDB(ctx context.Context, dbIndex int) ([]model.Order, error)
	FindOrderDetailsOnShardDB(ctx context.Context, dbIndex int) ([]model.OrderDetail, error)
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
