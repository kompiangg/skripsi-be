package cashier

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"

	"github.com/google/uuid"
)

func (r repository) FindByID(ctx context.Context, id uuid.UUID) (model.Cashier, error) {
	query := `
		SELECT
			id, account_id, store_id, "name", created_at
		FROM 
			cashiers
		WHERE
			id = ?;
	`

	var cashier model.Cashier
	err := r.generalDB.GetContext(ctx, &cashier, r.generalDB.Rebind(query), id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Cashier{}, errors.ErrRecordNotFound
	} else if err != nil {
		return model.Cashier{}, err
	}

	return cashier, nil
}
