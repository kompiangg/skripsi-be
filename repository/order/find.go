package order

import (
	"context"
	"skripsi-be/type/model"
	"skripsi-be/type/params"

	"github.com/Masterminds/squirrel"
	"github.com/go-errors/errors"
)

// Shard
func (r repository) FindAllOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.Order, error) {
	qBuilder := squirrel.Select("id, cashier_id, store_id, payment_id, customer_id, total_quantity, total_unit, total_price, total_price_in_usd, currency, usd_rate, created_at").
		From("orders")

	if param.StartDate.Valid && param.EndDate.Valid {
		qBuilder = qBuilder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": param.StartDate.Time},
			squirrel.LtOrEq{"created_at": param.EndDate.Time},
		})
	}

	q, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, errors.New(err)
	}

	var orders []model.Order
	err = r.shardDB[param.ShardIndex].SelectContext(ctx, &orders, r.shardDB[param.ShardIndex].Rebind(q), args...)
	if err != nil {
		return nil, errors.New(err)
	}

	return orders, nil
}

func (r repository) FindOrderDetailsOnShardDB(ctx context.Context, param params.FindOrderDetailsOnShardRepo) ([]model.OrderDetail, error) {
	qBuilder := squirrel.Select("id, order_id, item_id, quantity, unit, price").
		From("order_details")

	if param.OrderID.Valid {
		qBuilder = qBuilder.Where(squirrel.Eq{"order_id": param.OrderID})
	}

	q, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, errors.New(err)
	}

	var orderDetails []model.OrderDetail
	err = r.shardDB[param.ShardIndex].SelectContext(ctx, &orderDetails, r.shardDB[param.ShardIndex].Rebind(q), args...)
	if err != nil {
		return nil, errors.New(err)
	}

	return orderDetails, nil
}

// Long Term
func (r repository) FindAllOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.Order, error) {
	qBuilder := squirrel.Select("id, cashier_id, store_id, payment_id, customer_id, total_quantity, total_unit, total_price, total_price_in_usd, currency, usd_rate, created_at").
		From("orders")

	if param.StartDate.Valid && param.EndDate.Valid {
		qBuilder = qBuilder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": param.StartDate.Time},
			squirrel.LtOrEq{"created_at": param.EndDate.Time},
		})
	}

	q, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, errors.New(err)
	}

	var orders []model.Order
	err = r.longTermDB.SelectContext(ctx, &orders, r.longTermDB.Rebind(q), args...)
	if err != nil {
		return nil, errors.New(err)
	}

	return orders, nil
}

func (r repository) FindOrderDetailsOnLongTermDB(ctx context.Context, param params.FindOrderDetailsOnLongTermRepo) ([]model.OrderDetail, error) {
	qBuilder := squirrel.Select("id, order_id, item_id, quantity, unit, price").
		From("order_details")

	if param.OrderID.Valid {
		qBuilder = qBuilder.Where(squirrel.Eq{"order_id": param.OrderID})
	}

	q, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, errors.New(err)
	}

	var orderDetails []model.OrderDetail
	err = r.longTermDB.SelectContext(ctx, &orderDetails, r.longTermDB.Rebind(q), args...)
	if err != nil {
		return nil, errors.New(err)
	}

	return orderDetails, nil
}
