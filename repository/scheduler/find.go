package scheduler

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) FindByName(ctx context.Context, name string) (model.Scheduler, error) {
	query := `
		SELECT 
			id, name, run_count, last_run_at, created_at
		FROM 
			schedulers
		WHERE
			name = ?;
	`

	var scheduler model.Scheduler
	err := r.generalDB.GetContext(ctx, &scheduler, r.generalDB.Rebind(query), name)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Scheduler{}, errors.ErrNotFound
	} else if err != nil {
		return model.Scheduler{}, err
	}

	return scheduler, nil
}
