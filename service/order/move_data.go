package order

import (
	"context"
	"skripsi-be/type/constant"

	"github.com/go-errors/errors"
)

func (s service) MoveDataThroughShard(ctx context.Context) error {
	shardCount, err := s.shardRepo.GetShardCountScheduler(ctx)
	if errors.Is(err, constant.ErrRedisNil) {
		shardCount = 0
		err = nil
	} else if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	shardCount++
	for idx := len(s.config.Shards) - 1; idx >= 0; idx-- {
		if shardCount%s.config.Shards[idx].RangeInDay != 0 {
			continue
		}

		currentData, err := s.orderRepo.FindAllOnShardDB(ctx, idx)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		currentIndexTx, err := s.beginShardTx(ctx, idx)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		err = s.orderRepo.DeleteAllDataFromOneDB(ctx, currentIndexTx)
		if err != nil {
			currentIndexTx.Rollback()
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		if idx != len(s.config.Shards)-1 {
			afterIndexTx, err := s.beginShardTx(ctx, idx+1)
			if err != nil {
				currentIndexTx.Rollback()
				return errors.Wrap(err, constant.SkipErrorParameter)
			}

			err = s.orderRepo.InsertToShardDB(ctx, afterIndexTx, currentData)
			if err != nil {
				afterIndexTx.Rollback()
				currentIndexTx.Rollback()
				return errors.Wrap(err, constant.SkipErrorParameter)
			}

			err = afterIndexTx.Commit()
			if err != nil {
				currentIndexTx.Rollback()
				return errors.Wrap(err, constant.SkipErrorParameter)
			}
		}

		err = currentIndexTx.Commit()
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}
	}

	if shardCount == s.config.Shards[len(s.config.Shards)-1].RangeInDay {
		shardCount = 0
	}

	err = s.shardRepo.SetShardCountScheduler(ctx, shardCount)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	return nil
}
