package item

import (
	"context"
	"database/sql"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
)

func (r repository) FindByID(ctx context.Context, id string) (model.Item, error) {
	query := `
		SELECT 
			id, "name", "desc", price, origin_country, supplier, unit, created_at
		FROM 
			items
		WHERE
			id = ?;
	`

	var item model.Item
	err := r.generalDB.GetContext(ctx, &item, r.generalDB.Rebind(query), id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Item{}, errors.ErrRecordNotFound
	} else if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (r repository) FindLikeNameOrID(ctx context.Context, nameOrID string) ([]model.Item, error) {
	query := `
		SELECT 
			id, "name", "desc", price, origin_country, supplier, unit, created_at
		FROM 
			items
		WHERE
			"name" LIKE ? OR id LIKE ?;
	`

	var items []model.Item
	err := r.generalDB.SelectContext(ctx, &items, r.generalDB.Rebind(query), "%"+nameOrID+"%", "%"+nameOrID+"%")
	if err != nil {
		return nil, err
	}

	return items, nil
}
