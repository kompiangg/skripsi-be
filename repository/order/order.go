package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/type/model"
	"skripsi-be/type/params"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	InsertToShardDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	InsertDetailsToShardDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error
	InsertToLongTermDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	InsertDetailsToLongTermDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error
	DeleteAllData(ctx context.Context, tx *sqlx.Tx) error
	DeleteOrderDetails(ctx context.Context, tx *sqlx.Tx) error
	FindAllOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.Order, error)
	FindOrderDetailsOnShardDB(ctx context.Context, param params.FindOrderDetailsOnShardRepo) ([]model.OrderDetail, error)
	FindAllOrderAndDetailsOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.OrderWithOrderDetails, error)
	FindAllOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.Order, error)
	FindOrderDetailsOnLongTermDB(ctx context.Context, param params.FindOrderDetailsOnLongTermRepo) ([]model.OrderDetail, error)
	FindAllOrderAndDetailsOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.OrderWithOrderDetails, error)
}

type Config struct {
	ShardingDatabase config.ShardingDatabase
}

type repository struct {
	config     Config
	shardDB    []*sqlx.DB
	longTermDB *sqlx.DB
}

func New(
	config Config,
	shardDB []*sqlx.DB,
	longTermDB *sqlx.DB,
) Repository {
	return repository{
		config:     config,
		shardDB:    shardDB,
		longTermDB: longTermDB,
	}
}
