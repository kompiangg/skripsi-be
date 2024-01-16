package admin

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"

	"github.com/google/uuid"
)

func (r repository) FindAdminByID(ctx context.Context, id uuid.UUID) (model.Admin, error) {
	q := `
	SELECT
		id,
		account_id,
		name,
		created_at
	FROM
		admins
	WHERE
		id = $1
	;
	`

	var admin model.Admin
	err := r.generalDB.GetContext(ctx, &admin, q, id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Admin{}, errors.New(errors.ErrRecordNotFound)
	} else if err != nil {
		return model.Admin{}, errors.New(err)
	}

	return admin, nil
}

func (r repository) FindAdminAccountByID(ctx context.Context, id uuid.UUID) (model.Account, error) {
	q := `
	SELECT
		accounts.id as id,
		accounts.username as username,
		accounts.password as password,
		accounts.created_at as created_at
	FROM
		accounts
	INNER JOIN
		admins ON admins.account_id = accounts.id
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
