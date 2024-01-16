package repository

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/account"
	"skripsi-be/repository/admin"
	"skripsi-be/repository/order"
	"skripsi-be/repository/shard"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Sharding     shard.Repository
	Order        order.Repository
	Account      account.Repository
	Admin        admin.Repository
	LongTermDBTx func(ctx context.Context) (*sqlx.Tx, error)
	ShardDBTx    func(ctx context.Context, dbIndex int) (*sqlx.Tx, error)
}

func beginLongTermDBTx(longtermDB *sqlx.DB) func(ctx context.Context) (*sqlx.Tx, error) {
	return func(ctx context.Context) (*sqlx.Tx, error) {
		return longtermDB.BeginTxx(context.Background(), nil)
	}
}

func beginShardDBTx(shardingDatabase []*sqlx.DB) func(ctx context.Context, dbIndex int) (*sqlx.Tx, error) {
	return func(ctx context.Context, dbIndex int) (*sqlx.Tx, error) {
		return shardingDatabase[dbIndex].BeginTxx(ctx, nil)
	}
}

func New(
	config config.Config,
	longTermDatabase *sqlx.DB,
	generalDatabase *sqlx.DB,
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

	account := account.New(
		account.Config{},
		generalDatabase,
	)

	admin := admin.New(
		admin.Config{},
		generalDatabase,
	)

	return Repository{
		Sharding:     sharding,
		Order:        order,
		Account:      account,
		Admin:        admin,
		LongTermDBTx: beginLongTermDBTx(longTermDatabase),
		ShardDBTx:    beginShardDBTx(shardingDatabase),
	}, nil
}
