package constant

import (
	"github.com/go-errors/errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisNil = errors.New(redis.Nil)
)

var (
	ShardSchedulerCount = "shard_scheduler_count"
)
