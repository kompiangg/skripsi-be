package account

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) InsertNewAccount(ctx context.Context, account model.Account) error {
	q := `
	INSERT INTO accounts (
		id,
		username,
		password
	) VALUES (
		:id,
		:username,
		:password
	);
	`

	_, err := r.generalDB.NamedExecContext(ctx, q, account)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
