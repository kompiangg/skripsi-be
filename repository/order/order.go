package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/type/model"
	"skripsi-be/type/params"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	InsertToShardDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	InsertDetailsToShardDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error
	DeleteAllData(ctx context.Context, tx *sqlx.Tx) error
	DeleteOrderDetails(ctx context.Context, tx *sqlx.Tx) error
	FindAllOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.Order, error)
	FindOrderDetailsOnShardDB(ctx context.Context, param params.FindOrderDetailsOnShardRepo) ([]model.OrderDetail, error)
	FindAllOrderAndDetailsOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.OrderWithOrderDetails, error)
	FindOrderByIDOnShardDB(ctx context.Context, shardIdx int, id string) (model.Order, error)
	FindOrderDetailsByOrderIDOnShardDB(ctx context.Context, shardIdx int, orderID string) ([]model.OrderDetail, error)

	InsertToLongTermDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error
	InsertDetailsToLongTermDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error
	FindAllOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.Order, error)
	FindOrderDetailsOnLongTermDB(ctx context.Context, param params.FindOrderDetailsOnLongTermRepo) ([]model.OrderDetail, error)
	FindAllOrderAndDetailsOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.OrderWithOrderDetails, error)
	FindOrderByIDOnLongTermDB(ctx context.Context, id string) (model.Order, error)
	FindOrderDetailsByOrderIDOnLongTermDB(ctx context.Context, orderID string) ([]model.OrderDetail, error)

	GetAggregateTopSellingProductOnLongTermDB(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.GetAggregateTopSellingProductResultRepo, error)
	GetAggregateOrderOnLongTermDB(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.GetAggregateOrderResultRepo, error)
	GetAggregateTopSellingProductOnShardDB(ctx context.Context, dbIdx int, startDate time.Time, endDate time.Time) ([]model.GetAggregateTopSellingProductResultRepo, error)
	GetAggregateOrderOnShardDB(ctx context.Context, dbIdx int, startDate time.Time, endDate time.Time) ([]model.GetAggregateOrderResultRepo, error)
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
