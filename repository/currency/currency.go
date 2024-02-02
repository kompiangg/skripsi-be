package currency

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type Repository interface {
	FindUSDCurrencyRate(ctx context.Context, currency string) (decimal.Decimal, error)
}

type Config struct {
}

type repository struct {
	config      Config
	redisClient *redis.Client
}

func New(
	config Config,
	redisClient *redis.Client,
) Repository {
	return repository{
		config:      config,
		redisClient: redisClient,
	}
}
