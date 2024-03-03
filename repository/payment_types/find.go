package payment_types

import (
	"context"
	"skripsi-be/type/model"
)

func (r repository) FindByLikeOneOfAllColumn(ctx context.Context, req string) ([]model.PaymentType, error) {
	query := `
		SELECT 
			id, "type", "bank", created_at
		FROM 
			payment_types
		WHERE
			id LIKE ? OR "type" LIKE ? OR "bank" LIKE ?;
	`

	var items []model.PaymentType
	err := r.generalDB.SelectContext(ctx, &items, r.generalDB.Rebind(query), "%"+req+"%", "%"+req+"%", "%"+req+"%")
	if err != nil {
		return nil, err
	}

	return items, nil
}
