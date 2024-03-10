package scheduler

import (
	"context"
	"skripsi-be/pkg/errors"

	"github.com/jmoiron/sqlx"
)

func (r repository) IncrementRunCount(ctx context.Context, tx *sqlx.Tx, jobName string) error {
	q := `
		UPDATE schedulers
		SET 
			run_count = run_count + 1,
			last_run_at = NOW()
		WHERE name = ?;
	`

	_, err := tx.ExecContext(ctx, tx.Rebind(q), jobName)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
