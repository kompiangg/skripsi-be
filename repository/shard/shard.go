package shard

import (
	"context"
	"skripsi-be/config"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	GetShardCountScheduler(ctx context.Context) (int, error)
	SetShardCountScheduler(ctx context.Context, count int) error
}

type Config struct {
	Shards config.ShardingDatabase
}

type repository struct {
	config Config
	redis  *redis.Client
}

func New(
	config Config,
	redis *redis.Client,
) Repository {
	return repository{
		config: config,
		redis:  redis,
	}
}
