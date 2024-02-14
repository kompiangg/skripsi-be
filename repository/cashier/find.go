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

func (r repository) FindCashierAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error) {
	q := `
	SELECT
		accounts.id as id,
		accounts.username as username,
		accounts.password as password,
		accounts.created_at as created_at
	FROM
		accounts
	INNER JOIN
		cashiers ON cashiers.account_id = accounts.id
	WHERE
		accounts.id = $1
	;
	`

	var account model.Account
	err := r.generalDB.GetContext(ctx, &account, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Account{}, errors.New(errors.ErrRecordNotFound)
	} else if err != nil {
		return model.Account{}, errors.New(err)
	}

	return account, nil
}
