package order

import (
	"context"
	"skripsi-be/type/constant"
	"skripsi-be/type/params"

	"github.com/go-errors/errors"
)

func (s service) MoveDataThroughShard(ctx context.Context) error {
	shardScheduler, err := s.scheduler.FindByName(ctx, constant.ShardSchedulerCount)
	if err != nil {
		errors.Wrap(err, constant.SkipErrorParameter)
	}

	shardScheduler.RunCount++
	for idx := len(s.config.Shards) - 1; idx >= 0; idx-- {
		if shardScheduler.RunCount%s.config.Shards[idx].DataRetention != 0 {
			continue
		}

		currentData, err := s.orderRepo.FindAllOnShardDB(ctx, params.ShardTimeSeriesWhereQuery{
			ShardIndex: idx,
		})
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		currentOrderDetails, err := s.orderRepo.FindOrderDetailsOnShardDB(ctx, params.FindOrderDetailsOnShardRepo{
			ShardIndex: idx,
		})
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		currentIndexTx, err := s.beginShardTx(ctx, idx)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		err = s.orderRepo.DeleteOrderDetails(ctx, currentIndexTx)
		if err != nil {
			currentIndexTx.Rollback()
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		err = s.orderRepo.DeleteAllData(ctx, currentIndexTx)
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

			err = s.orderRepo.InsertDetailsToShardDB(ctx, afterIndexTx, currentOrderDetails)
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

	if shardScheduler.RunCount == s.config.Shards[len(s.config.Shards)-1].DataRetention {
		shardScheduler.RunCount = 0
	}

	tx, err := s.beginGeneralTx(ctx)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	defer tx.Rollback()

	err = s.scheduler.IncrementRunCount(ctx, tx, constant.ShardSchedulerCount)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	return nil
}
