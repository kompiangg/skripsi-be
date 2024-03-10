package order

import (
	"context"
	"skripsi-be/type/constant"
	"skripsi-be/type/model"
	"skripsi-be/type/params"

	"github.com/go-errors/errors"
)

func (s service) InsertToLongTerm(ctx context.Context, param params.ServiceInsertOrdersToLongTermParam) error {
	err := param.Validate(ctx)
	if err != nil {
		return errors.New(err)
	}

	orders := make([]model.Order, len(param))
	for idx, order := range param {
		orders[idx] = order.ToOrderModel()
	}

	tx, err := s.beginLongTermTx(ctx)
	if err != nil {
		return errors.New(err)
	}

	err = s.orderRepo.InsertToLongTermDB(ctx, tx, orders)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	for idx := range orders {
		err = s.orderRepo.InsertDetailsToLongTermDB(ctx, tx, orders[idx].OrderDetails)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, constant.SkipErrorParameter)
		}
	}

	tx.Commit()

	return nil
}

func (s service) InsertToLongTermSeeder(ctx context.Context, param params.ServiceInsertOrdersToLongTermParam) error {
	err := param.Validate(ctx)
	if err != nil {
		return errors.New(err)
	}

	tx, err := s.beginLongTermTx(ctx)
	if err != nil {
		return errors.New(err)
	}

	orders := make([]model.Order, len(param))
	for idx, order := range param {
		orders[idx] = order.ToOrderModelInSeeder()
	}

	err = s.orderRepo.InsertToLongTermDB(ctx, tx, orders)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	for idx := range orders {
		err = s.orderRepo.InsertDetailsToLongTermDB(ctx, tx, orders[idx].OrderDetails)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, constant.SkipErrorParameter)
		}
	}

	tx.Commit()

	return nil
}
