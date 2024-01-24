package order

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

func (r repository) DeleteAllData(ctx context.Context, tx *sqlx.Tx) error {
	deleteQuery := `
		DELETE FROM orders
	`

	_, err := tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

func (r repository) DeleteOrderDetails(ctx context.Context, tx *sqlx.Tx) error {
	deleteQuery := `
		DELETE FROM order_details
	`

	_, err := tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
