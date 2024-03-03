package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

type Order struct {
	ID              string          `db:"id"`
	CashierID       uuid.UUID       `db:"cashier_id"`
	StoreID         null.String     `db:"store_id"`
	PaymentID       null.String     `db:"payment_id"`
	CustomerID      null.String     `db:"customer_id"`
	TotalQuantity   null.Int64      `db:"total_quantity"`
	TotalPrice      decimal.Decimal `db:"total_price"`
	TotalPriceInUSD decimal.Decimal `db:"total_price_in_usd"`
	Currency        null.String     `db:"currency"`
	UsdRate         decimal.Decimal `db:"usd_rate"`
	CreatedAt       time.Time       `db:"created_at"`
	OrderDetails    []OrderDetail   `db:"-"`
}

type OrderDetail struct {
	ID       uuid.UUID       `db:"id"`
	OrderID  string          `db:"order_id"`
	ItemID   null.String     `db:"item_id"`
	Quantity null.Int64      `db:"quantity"`
	Unit     null.String     `db:"unit"`
	Price    decimal.Decimal `db:"price"`
}

type OrderWithOrderDetails struct {
	OrderID         string          `db:"order_id"`
	OrderDetailID   uuid.UUID       `db:"order_detail_id"`
	CashierID       uuid.UUID       `db:"cashier_id"`
	StoreID         null.String     `db:"store_id"`
	PaymentID       null.String     `db:"payment_id"`
	CustomerID      null.String     `db:"customer_id"`
	TotalQuantity   null.Int64      `db:"total_quantity"`
	TotalPrice      decimal.Decimal `db:"total_price"`
	TotalPriceInUSD decimal.Decimal `db:"total_price_in_usd"`
	Currency        null.String     `db:"currency"`
	UsdRate         decimal.Decimal `db:"usd_rate"`
	CreatedAt       time.Time       `db:"created_at"`
	ItemID          null.String     `db:"item_id"`
	Quantity        null.Int64      `db:"quantity"`
	Unit            null.String     `db:"unit"`
	Price           decimal.Decimal `db:"price"`
}

type GetAggregateTopSellingProductResultRepo struct {
	ItemID                string `db:"item_id"`
	ItemSoldTotalQuantity int64  `db:"item_sold_total_quantity"`
}

type GetAggregateOrderResultRepo struct {
	CustomerID             string          `db:"customer_id"`
	OrderQuantity          int64           `db:"order_quantity"`
	NotMemberOrderQuantity int64           `db:"not_member_order_quantity"`
	ItemSoldTotalQuantity  int64           `db:"item_sold_total_quantity"`
	ItemSoldTotalPrice     decimal.Decimal `db:"item_sold_total_price"`
}
