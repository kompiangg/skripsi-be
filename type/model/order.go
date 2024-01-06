package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v9"
)

type Order struct {
	ID         uuid.UUID     `db:"id"`
	PaymentID  uuid.NullUUID `db:"payment_id"`
	CustomerID null.String   `db:"customer_id"`
	ItemID     uuid.NullUUID `db:"item_id"`
	StoreID    uuid.NullUUID `db:"store_id"`
	Quantity   null.Int      `db:"quantity"`
	Unit       null.String   `db:"unit"`
	Price      null.Int64    `db:"price"`
	TotalPrice null.Int64    `db:"total_price"`
	CreatedAt  time.Time     `db:"created_at"`
}
