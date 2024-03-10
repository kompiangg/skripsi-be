package currency

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) FindByBaseAndQuote(ctx context.Context, base string, quote string) (model.Currency, error) {
	q := `
		SELECT
			id,
			base,
			quote,
			rate,
			created_at
		FROM
			currencies
		WHERE
			base = ?
			AND quote = ?;
	`

	var res model.Currency
	err := r.generalDB.GetContext(ctx, &res, r.generalDB.Rebind(q), base, quote)
	if errors.Is(err, sql.ErrNoRows) {
		return res, errors.ErrNotFound
	} else if err != nil {
		return res, errors.Wrap(err)
	}

	return res, nil
}
