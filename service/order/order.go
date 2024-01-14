package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/order"
	"skripsi-be/repository/shard"
	"skripsi-be/type/params"

	"github.com/jmoiron/sqlx"
)

type Order interface {
	MoveDataThroughShard(ctx context.Context) error
	InsertToShard(ctx context.Context, param params.ServiceInsertOrdersToShardParam) error
	InsertToLongTerm(ctx context.Context, param params.ServiceInsertOrdersToLongTermParam) error
}

type Config struct {
	Shards config.Shards
	Date   config.Date
}

type service struct {
	config    Config
	shardRepo shard.Repository
	orderRepo order.Repository

	beginShardTx    func(ctx context.Context, dbIdx int) (*sqlx.Tx, error)
	beginLongTermTx func(ctx context.Context) (*sqlx.Tx, error)
}

func New(
	config Config,
	shardRepo shard.Repository,
	orderRepo order.Repository,
	beginShardTx func(ctx context.Context, dbIdx int) (*sqlx.Tx, error),
	beginLongTermTx func(ctx context.Context) (*sqlx.Tx, error),
) Order {
	return service{
		config:          config,
		shardRepo:       shardRepo,
		orderRepo:       orderRepo,
		beginShardTx:    beginShardTx,
		beginLongTermTx: beginLongTermTx,
	}
}
