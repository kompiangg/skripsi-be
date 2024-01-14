package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v9"
)

type Order struct {
	ID         uuid.UUID     `db:"id"`
	ItemID     null.String   `db:"item_id"`
	StoreID    null.String   `db:"store_id"`
	PaymentID  null.String   `db:"payment_id"`
	CustomerID null.String   `db:"customer_id"`
	CashierID  uuid.NullUUID `db:"cashier_id"`
	Quantity   null.Int      `db:"quantity"`
	Unit       null.String   `db:"unit"`
	Price      null.Float64  `db:"price"`
	TotalPrice null.Float64  `db:"total_price"`
	CreatedAt  time.Time     `db:"created_at"`
}
