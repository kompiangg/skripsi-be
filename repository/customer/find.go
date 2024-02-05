package customer

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) FindByID(ctx context.Context, id string) (model.Customer, error) {
	query := `
		SELECT 
			id, "name", contact, created_at
		FROM 
			customers
		WHERE
			id = ?;
	`

	var customer model.Customer
	err := r.generalDB.GetContext(ctx, &customer, r.generalDB.Rebind(query), id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Customer{}, errors.ErrRecordNotFound
	} else if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}
