package currency

import (
	"context"
	"fmt"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

func (r repository) FindUSDCurrencyRate(ctx context.Context, currency string) (decimal.Decimal, error) {
	key := fmt.Sprintf("USD_%s", currency)

	rate, err := r.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return decimal.Zero, constant.ErrRedisNil
	} else if err != nil {
		return decimal.Zero, errors.Wrap(err)
	}

	res, err := decimal.NewFromString(rate)
	if err != nil {
		return decimal.Zero, errors.Wrap(err)
	}

	return res, nil
}
