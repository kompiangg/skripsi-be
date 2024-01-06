package order

import (
	"context"
	"skripsi-be/type/model"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

func (r repository) InsertDataToShardDB(ctx context.Context, tx *sqlx.Tx, orders []model.Order) error {
	q := `
		INSERT INTO orders (
				id, payment_id, customer_id, item_id, store_id, 
				quantity, unit, price, total_price, discount, created_at
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
			order.ID, order.PaymentID, order.CustomerID, order.ItemID,
			order.StoreID, order.Quantity, order.Unit, order.Price,
			order.TotalPrice, order.Discount, order.CreatedAt)
		if err != nil {
			tx.Rollback()
			return errors.New(err)
		}
	}

	return nil
}
