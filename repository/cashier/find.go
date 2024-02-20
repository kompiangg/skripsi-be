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

func (r repository) FindCashierAccountByID(ctx context.Context, id uuid.UUID) (model.Cashier, error) {
	q := `
	SELECT
		cashiers.id as id,
		cashiers.account_id as account_id,
		cashiers.store_id as store_id,
		cashiers.name as name,
		cashiers.created_at as created_at
	FROM 
		cashiers
	INNER JOIN
		accounts ON cashiers.account_id = accounts.id
	WHERE
		accounts.id = $1
	;
	`

	var cashier model.Cashier
	err := r.generalDB.GetContext(ctx, &cashier, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Cashier{}, errors.New(errors.ErrRecordNotFound)
	} else if err != nil {
		return model.Cashier{}, errors.New(err)
	}

	return cashier, nil
}
