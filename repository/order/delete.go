package order

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

func (r repository) DeleteAllDataFromOneDB(ctx context.Context, tx *sqlx.Tx) error {
	deleteQuery := `
		DELETE FROM orders
	`

	_, err := tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		tx.Rollback()
		return errors.New(err)
	}

	return nil
}
