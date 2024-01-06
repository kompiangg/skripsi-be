package order

import (
	"context"
	"skripsi-be/type/model"

	"github.com/go-errors/errors"
)

func (r repository) FindAllOnShardDB(ctx context.Context, dbIndex int) ([]model.Order, error) {
	q := `
		SELECT
			id,
			payment_id,
			customer_id,
			item_id,
			store_id,
			quantity,
			unit,
			price,
			total_price,
			created_at
		FROM
			orders
	`

	var orders []model.Order
	err := r.shardDB[dbIndex].SelectContext(ctx, &orders, q)
	if err != nil {
		return nil, errors.New(err)
	}

	return orders, nil
}
