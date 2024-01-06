package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `db:"id"`
	PaymentID  uuid.UUID `db:"payment_id"`
	CustomerID uuid.UUID `db:"customer_id"`
	ItemID     uuid.UUID `db:"item_id"`
	StoreID    uuid.UUID `db:"store_id"`
	Quantity   int       `db:"quantity"`
	Unit       string    `db:"unit"`
	Price      int64     `db:"price"`
	TotalPrice int64     `db:"total_price"`
	Discount   int64     `db:"discount"`
	CreatedAt  time.Time `db:"created_at"`
}
