package order

import (
	"context"
	"skripsi-be/type/constant"
	"skripsi-be/type/model"
	"skripsi-be/type/params"

	"github.com/go-errors/errors"
)

func (s service) InsertToShard(ctx context.Context, param params.ServiceInsertOrdersToShardParam) error {
	err := param.Validate(ctx)
	if err != nil {
		return errors.New(err)
	}

	orderDBIndex := make([][]model.Order, len(s.config.Shards))
	for _, order := range param {
		dbIdx, err := s.getShardIndexByDateTime(order.CreatedAt)
		if err != nil {
			err = nil
			continue
		}

		orderDBIndex[dbIdx] = append(orderDBIndex[dbIdx], order.ToOrderModel())
	}

	for dbIdx, orders := range orderDBIndex {
		if len(orders) == 0 {
			continue
		}

		tx, err := s.beginShardTx(ctx, dbIdx)
		if err != nil {
			return errors.New(err)
		}

		err = s.orderRepo.InsertToShardDB(ctx, tx, orders)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		for idx := range orders {
			err = s.orderRepo.InsertDetailsToShardDB(ctx, tx, orders[idx].OrderDetails)
			if err != nil {
				tx.Rollback()
				return errors.Wrap(err, constant.SkipErrorParameter)
			}
		}

		tx.Commit()
	}

	return nil
}
