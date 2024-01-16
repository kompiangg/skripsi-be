package admin

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) InsertNewAdmin(ctx context.Context, admin model.Admin) error {
	q := `
	INSERT INTO admins (
		id,
		account_id,
		name
	) VALUES (
		:id,
		:account_id,
		:name
	);
	`

	_, err := r.generalDB.NamedExecContext(ctx, q, admin)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
