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
			id, cashier_id, item_id, store_id, payment_id, customer_id, 
			quantity, unit, price, total_price, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, tx.Rebind(q))
	if err != nil {
		tx.Rollback()
		return errors.New(err)
	}
	defer stmt.Close()

	for _, order := range orders {
		_, err = stmt.ExecContext(ctx,
			order.ID, order.CashierID, order.ItemID, order.StoreID, order.PaymentID,
			order.CustomerID, order.Quantity, order.Unit, order.Price,
			order.TotalPrice, order.CreatedAt)
		if err != nil {
			tx.Rollback()
			return errors.New(err)
		}
	}

	return nil
}

func (r repository) InsertToLongTermDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error {
	q := `
	INSERT INTO orders (
		id, cashier_id, item_id, store_id, payment_id, customer_id, 
		quantity, unit, price, total_price, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	stmt, err := tx.PrepareContext(ctx, tx.Rebind(q))
	if err != nil {
		tx.Rollback()
		return errors.New(err)
	}
	defer stmt.Close()

	for _, order := range orders {
		_, err = stmt.ExecContext(ctx,
			order.ID, order.CashierID, order.ItemID, order.StoreID, order.PaymentID,
			order.CustomerID, order.Quantity, order.Unit, order.Price,
			order.TotalPrice, order.CreatedAt)
		if err != nil {
			tx.Rollback()
			return errors.New(err)
		}
	}

	return nil
}
