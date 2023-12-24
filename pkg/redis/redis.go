package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisConfig struct {
	Hostname string
	Username string
	Password string
	DB       int
}

const maxRetry = 5

func New(config RedisConfig) (*redis.Client, error) {
	var rdb *redis.Client
	sleepDuration := 1
	retry := 0

	for {
		rdb = redis.NewClient(&redis.Options{
			Username: config.Username,
			Addr:     config.Hostname,
			Password: config.Password,
			DB:       config.DB,
		})

		ctx := context.Background()

		err := rdb.Ping(ctx).Err()
		if err == nil {
			break
		}

		if retry == maxRetry {
			log.Fatal().Msgf("failed to connect to redis, retry: %d", retry)
			return nil, err
		}

		log.Info().Msgf("retrying to connect to redis, retry: %d", retry)
		log.Info().Msgf("cause: %v", err)
		time.Sleep(time.Duration(sleepDuration * int(time.Second)))
		sleepDuration++
	}

	return rdb, nil
}
