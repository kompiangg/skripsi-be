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

func (r repository) FindAdminAccountByID(ctx context.Context, id uuid.UUID) (model.Admin, error) {
	q := `
	SELECT
		admins.id as id,
		admins.account_id as account_id,
		admins.name as name,
		admins.created_at as created_at
	FROM
		admins
	INNER JOIN
		accounts ON admins.account_id = accounts.id
	WHERE
		accounts.id = $1
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
