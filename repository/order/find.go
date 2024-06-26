package order

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
	"skripsi-be/type/params"

	"github.com/Masterminds/squirrel"
)

// Shard
func (r repository) FindAllOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.Order, error) {
	qBuilder := squirrel.
		Select("id, cashier_id, store_id, payment_id, customer_id, total_quantity, total_price, total_price_in_usd, currency, usd_rate, created_at").
		From("orders").
		OrderBy("created_at DESC")

	if param.StartDate.Valid && param.EndDate.Valid {
		qBuilder = qBuilder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": param.StartDate.Time},
			squirrel.LtOrEq{"created_at": param.EndDate.Time},
		})
	}

	if param.CashierID.Valid {
		qBuilder = qBuilder.Where(squirrel.Eq{"cashier_id": param.CashierID})
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

func (r repository) FindAllOrderAndDetailsOnShardDB(ctx context.Context, param params.ShardTimeSeriesWhereQuery) ([]model.OrderWithOrderDetails, error) {
	qBuilder := squirrel.
		Select(
			"o.id as order_id",
			"od.id as order_detail_id",
			"o.cashier_id as cashier_id",
			"o.store_id as store_id",
			"o.payment_id as payment_id",
			"o.customer_id as customer_id",
			"o.total_quantity as total_quantity",
			"o.total_price as total_price",
			"o.total_price_in_usd as total_price_in_usd",
			"o.currency as currency",
			"o.usd_rate as usd_rate",
			"o.created_at as created_at",
			"od.item_id as item_id",
			"od.quantity as quantity",
			"od.unit as unit",
			"od.price as price",
		).
		From("orders o").
		InnerJoin("order_details od ON o.id = od.order_id")

	if param.StartDate.Valid && param.EndDate.Valid {
		qBuilder = qBuilder.Where(squirrel.And{
			squirrel.GtOrEq{"o.created_at": param.StartDate.Time},
			squirrel.LtOrEq{"o.created_at": param.EndDate.Time},
		})
	}

	q, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, errors.New(err)
	}

	var orders []model.OrderWithOrderDetails
	err = r.shardDB[param.ShardIndex].SelectContext(ctx, &orders, r.shardDB[param.ShardIndex].Rebind(q), args...)
	if err != nil {
		return nil, errors.New(err)
	}

	return orders, nil
}

func (r repository) FindOrderByIDOnShardDB(ctx context.Context, shardIdx int, id string) (model.Order, error) {
	q := `
		SELECT 
			id, 
			cashier_id, 
			store_id, 
			payment_id, 
			customer_id, 
			total_quantity, 
			total_price, 
			total_price_in_usd, 
			currency, 
			usd_rate, 
			created_at
		FROM 
			orders 
		WHERE id = ?`

	var order model.Order
	err := r.longTermDB.GetContext(ctx, &order, r.longTermDB.Rebind(q), id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, errors.ErrNotFound
	} else if err != nil {
		return model.Order{}, errors.New(err)
	}

	return order, nil
}

func (r repository) FindOrderDetailsByOrderIDOnShardDB(ctx context.Context, shardIdx int, orderID string) ([]model.OrderDetail, error) {
	q := `
		SELECT 
			id, 
			order_id, 
			item_id, 
			quantity, 
			unit, 
			price
		FROM 
			order_details 
		WHERE order_id = ?
	`

	var orderDetails []model.OrderDetail
	err := r.longTermDB.SelectContext(ctx, &orderDetails, r.longTermDB.Rebind(q), orderID)
	if err != nil {
		return nil, errors.New(err)
	}

	return orderDetails, nil
}

// Long Term
func (r repository) FindAllOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.Order, error) {
	qBuilder := squirrel.Select("id, cashier_id, store_id, payment_id, customer_id, total_quantity, total_price, total_price_in_usd, currency, usd_rate, created_at").
		From("orders").
		OrderBy("created_at DESC")

	if param.StartDate.Valid && param.EndDate.Valid {
		qBuilder = qBuilder.Where(squirrel.And{
			squirrel.GtOrEq{"created_at": param.StartDate.Time},
			squirrel.LtOrEq{"created_at": param.EndDate.Time},
		})
	}

	if param.CashierID.Valid {
		qBuilder = qBuilder.Where(squirrel.Eq{"cashier_id": param.CashierID})
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

func (r repository) FindAllOrderAndDetailsOnLongTermDB(ctx context.Context, param params.LongTermWhereQuery) ([]model.OrderWithOrderDetails, error) {
	qBuilder := squirrel.
		Select(
			"o.id as order_id",
			"od.id as order_detail_id",
			"o.cashier_id as cashier_id",
			"o.store_id as store_id",
			"o.payment_id as payment_id",
			"o.customer_id as customer_id",
			"o.total_quantity as total_quantity",
			"o.total_price as total_price",
			"o.total_price_in_usd as total_price_in_usd",
			"o.currency as currency",
			"o.usd_rate as usd_rate",
			"o.created_at as created_at",
			"od.item_id as item_id",
			"od.quantity as quantity",
			"od.unit as unit",
			"od.price as price",
		).
		From("orders o").
		InnerJoin("order_details od ON o.id = od.order_id")

	if param.StartDate.Valid && param.EndDate.Valid {
		qBuilder = qBuilder.Where(squirrel.And{
			squirrel.GtOrEq{"o.created_at": param.StartDate.Time},
			squirrel.LtOrEq{"o.created_at": param.EndDate.Time},
		})
	}

	q, args, err := qBuilder.ToSql()
	if err != nil {
		return nil, errors.New(err)
	}

	var orders []model.OrderWithOrderDetails
	err = r.longTermDB.SelectContext(ctx, &orders, r.longTermDB.Rebind(q), args...)
	if err != nil {
		return nil, errors.New(err)
	}

	return orders, nil
}

func (r repository) FindOrderByIDOnLongTermDB(ctx context.Context, id string) (model.Order, error) {
	q := `
		SELECT 
			id, 
			cashier_id, 
			store_id, 
			payment_id, 
			customer_id, 
			total_quantity, 
			total_price, 
			total_price_in_usd, 
			currency, 
			usd_rate, 
			created_at
		FROM 
			orders 
		WHERE id = ?`

	var order model.Order
	err := r.longTermDB.GetContext(ctx, &order, r.longTermDB.Rebind(q), id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Order{}, errors.ErrItemNotFound
	} else if err != nil {
		return model.Order{}, errors.New(err)
	}

	return order, nil
}

func (r repository) FindOrderDetailsByOrderIDOnLongTermDB(ctx context.Context, orderID string) ([]model.OrderDetail, error) {
	q := `
		SELECT 
			id, 
			order_id, 
			item_id, 
			quantity, 
			unit, 
			price
		FROM 
			order_details 
		WHERE order_id = ?
	`

	var orderDetails []model.OrderDetail
	err := r.longTermDB.SelectContext(ctx, &orderDetails, r.longTermDB.Rebind(q), orderID)
	if err != nil {
		return nil, errors.New(err)
	}

	return orderDetails, nil
}
