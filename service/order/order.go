package order

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/order"
	"skripsi-be/repository/shard"

	"github.com/jmoiron/sqlx"
)

type Order interface {
	MoveDataThroughShard(ctx context.Context) error
}

type Config struct {
	Shards []config.Shard
}

type service struct {
	config    Config
	shardRepo shard.Repository
	orderRepo order.Repository

	beginShardTx func(ctx context.Context, dbIdx int) (*sqlx.Tx, error)
}

func New(
	config Config,
	shardRepo shard.Repository,
	orderRepo order.Repository,
	beginShardTx func(ctx context.Context, dbIdx int) (*sqlx.Tx, error),
) Order {
	return service{
		config:       config,
		shardRepo:    shardRepo,
		orderRepo:    orderRepo,
		beginShardTx: beginShardTx,
	}
}
