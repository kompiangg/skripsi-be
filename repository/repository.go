package repository

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/order"
	"skripsi-be/repository/shard"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Sharding shard.Repository
	Order    order.Repository

	longTermDB *sqlx.DB
	shardDB    []*sqlx.DB
}

func (r Repository) BeginLongTermDBTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.longTermDB.BeginTxx(ctx, nil)
}

func (r Repository) BeginShardDBTx(ctx context.Context, dbIndex int) (*sqlx.Tx, error) {
	return r.shardDB[dbIndex].BeginTxx(ctx, nil)
}

func New(
	config config.Config,
	longTermDatabase *sqlx.DB,
	shardingDatabase []*sqlx.DB,
	redis *redis.Client,
) (Repository, error) {
	sharding := shard.New(
		shard.Config{
			Shards: config.ShardingDatabase,
		},
		redis,
	)

	order := order.New(
		order.Config{
			ShardingDatabase: config.ShardingDatabase,
		},
		shardingDatabase,
	)

	return Repository{
		Sharding:   sharding,
		Order:      order,
		longTermDB: longTermDatabase,
		shardDB:    shardingDatabase,
	}, nil
}
