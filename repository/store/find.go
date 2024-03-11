package store

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) FindByID(ctx context.Context, id string) (model.Store, error) {
	query := `
		SELECT 
			id, nation, region, district, sub_district, currency, created_at
		FROM 
			stores
		WHERE
			id = ?;
	`

	var store model.Store
	err := r.generalDB.GetContext(ctx, &store, r.generalDB.Rebind(query), id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Store{}, errors.ErrRecordNotFound
	} else if err != nil {
		return model.Store{}, err
	}

	return store, nil
}
