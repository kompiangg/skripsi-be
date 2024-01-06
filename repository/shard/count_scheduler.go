package shard

import (
	"context"
	"skripsi-be/type/constant"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/redis/go-redis/v9"
)

func (s repository) GetShardCountScheduler(ctx context.Context) (int, error) {
	res, err := s.redis.Get(ctx, constant.ShardSchedulerCount).Result()
	if errors.Is(err, redis.Nil) {
		return 0, constant.ErrRedisNil
	} else if err != nil {
		return 0, errors.New(err)
	}

	count, err := strconv.Atoi(res)
	if err != nil {
		return 0, errors.New(err)
	}

	return count, nil
}

func (s repository) SetShardCountScheduler(ctx context.Context, count int) error {
	err := s.redis.Set(ctx, constant.ShardSchedulerCount, count, 0).Err()
	if err != nil {
		return errors.New(err)
	}

	return nil
}
