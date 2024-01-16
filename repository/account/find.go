package account

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"

	"github.com/google/uuid"
)

func (r repository) FindAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error) {
	q := `
	SELECT
		id,
		username,
		password,
		created_at
	FROM
		accounts
	WHERE
		id = $1
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

func (r repository) FindAccountByUsername(ctx context.Context, username string) (model.Account, error) {
	q := `
	SELECT
		id,
		username,
		password,
		created_at
	FROM
		accounts
	WHERE
		username = $1;
	`

	var account model.Account
	err := r.generalDB.GetContext(ctx, &account, q, username)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Account{}, errors.New(errors.ErrRecordNotFound)
	} else if err != nil {
		return model.Account{}, errors.New(err)
	}

	return account, nil
}
