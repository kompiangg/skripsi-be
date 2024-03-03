package order

import (
	"context"
	"skripsi-be/type/model"
	"time"
)

// Long Term DB
func (r repository) GetAggregateTopSellingProductOnLongTermDB(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.GetAggregateTopSellingProductResultRepo, error) {
	var result []model.GetAggregateTopSellingProductResultRepo
	var err error

	query := `
		SELECT
			item_id,
			SUM(quantity) as item_sold_total_quantity
		FROM
			order_details
		INNER JOIN
			orders on orders.id = order_details.order_id
		WHERE
			orders.created_at >= ? AND orders.created_at <= ?
		GROUP BY
			item_id
		ORDER BY
			item_sold_total_quantity DESC;
	`

	err = r.longTermDB.SelectContext(ctx, &result, r.longTermDB.Rebind(query), startDate, endDate)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r repository) GetAggregateOrderOnLongTermDB(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.GetAggregateOrderResultRepo, error) {
	var err error
	result := []model.GetAggregateOrderResultRepo{}

	query := `
		SELECT
			coalesce(customer_id, 'null') as customer_id,
			count(customer_id) as order_quantity,
			sum(total_quantity) as item_sold_total_quantity,
			sum(total_price_in_usd) as item_sold_total_price,
			count(case when orders.customer_id is null then '1' else orders.customer_id end) as not_member_order_quantity
		FROM
			orders
		WHERE
			created_at >= ? AND created_at <= ?
		GROUP BY
			customer_id;
	`

	err = r.longTermDB.SelectContext(ctx, &result, r.longTermDB.Rebind(query), startDate, endDate)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Shard DB
func (r repository) GetAggregateTopSellingProductOnShardDB(ctx context.Context, dbIdx int, startDate time.Time, endDate time.Time) ([]model.GetAggregateTopSellingProductResultRepo, error) {
	var result []model.GetAggregateTopSellingProductResultRepo
	var err error

	query := `
		SELECT
			item_id,
			SUM(quantity) as item_sold_total_quantity
		FROM
			order_details
		INNER JOIN
			orders on orders.id = order_details.order_id
		WHERE
			orders.created_at >= ? AND orders.created_at <= ?
		GROUP BY
			item_id
		ORDER BY
			item_sold_total_quantity DESC;
	`

	err = r.shardDB[dbIdx].SelectContext(ctx, &result, r.shardDB[dbIdx].Rebind(query), startDate, endDate)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r repository) GetAggregateOrderOnShardDB(ctx context.Context, dbIdx int, startDate time.Time, endDate time.Time) ([]model.GetAggregateOrderResultRepo, error) {
	var err error
	result := []model.GetAggregateOrderResultRepo{}

	query := `
		SELECT
			coalesce(customer_id, 'null') as customer_id,
			count(customer_id) as order_quantity,
			sum(total_quantity) as item_sold_total_quantity,
			sum(total_price_in_usd) as item_sold_total_price,
			count(case when orders.customer_id is null then '1' else orders.customer_id end) as not_member_order_quantity
		FROM
			orders
		WHERE
			created_at >= ? AND created_at <= ?
		GROUP BY
			customer_id;
	`

	err = r.shardDB[dbIdx].SelectContext(ctx, &result, r.shardDB[dbIdx].Rebind(query), startDate, endDate)
	if err != nil {
		return nil, err
	}

	return result, nil
}
