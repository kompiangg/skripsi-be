package repository

import (
	"skripsi-be/config"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
}

func New(
	config config.Config,
	longTermDatabase *sqlx.DB,
	shardingDatabase []*sqlx.DB,
	redis *redis.Client,
) (Repository, error) {
	return Repository{}, nil
}
