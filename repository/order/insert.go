package order

import (
	"context"
	"skripsi-be/type/model"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

func (r repository) InsertToShardDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error {
	q := `
		INSERT INTO orders (
			id, cashier_id, store_id, payment_id, customer_id, 
			total_quantity, total_unit, total_price, total_price_in_usd, currency, usd_rate, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, tx.Rebind(q))
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()

	for _, order := range orders {
		_, err = stmt.ExecContext(ctx,
			order.ID, order.CashierID, order.StoreID, order.PaymentID,
			order.CustomerID, order.TotalQuantity, order.TotalUnit, order.TotalPrice, order.TotalPriceInUSD, order.Currency,
			order.UsdRate, order.CreatedAt)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func (r repository) InsertDetailsToShardDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error {
	q := `
		INSERT INTO order_details (
			id, order_id, item_id, quantity, unit, price
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, tx.Rebind(q))
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()

	for _, orderDetail := range orderDetails {
		_, err = stmt.ExecContext(ctx,
			orderDetail.ID, orderDetail.OrderID, orderDetail.ItemID,
			orderDetail.Quantity, orderDetail.Unit, orderDetail.Price)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func (r repository) InsertToLongTermDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error {
	q := `
	INSERT INTO orders (
		id, cashier_id, store_id, payment_id, customer_id, 
		total_quantity, total_unit, total_price, total_price_in_usd, currency, usd_rate, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	stmt, err := tx.PrepareContext(ctx, tx.Rebind(q))
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()

	for _, order := range orders {
		_, err = stmt.ExecContext(ctx,
			order.ID, order.CashierID, order.StoreID, order.PaymentID,
			order.CustomerID, order.TotalQuantity, order.TotalUnit, order.TotalPrice, order.TotalPriceInUSD, order.Currency,
			order.UsdRate, order.CreatedAt)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func (r repository) InsertDetailsToLongTermDB(ctx context.Context, tx *sqlx.Tx, orderDetails []model.OrderDetail) error {
	q := `
	INSERT INTO order_details (
		id, order_id, item_id, quantity, unit, price
	) VALUES (?, ?, ?, ?, ?, ?)
`

	stmt, err := tx.PrepareContext(ctx, tx.Rebind(q))
	if err != nil {
		return errors.New(err)
	}
	defer stmt.Close()

	for _, orderDetail := range orderDetails {
		_, err = stmt.ExecContext(ctx,
			orderDetail.ID, orderDetail.OrderID, orderDetail.ItemID,
			orderDetail.Quantity, orderDetail.Unit, orderDetail.Price)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}
