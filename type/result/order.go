package result

import (
	"skripsi-be/type/model"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

type Order struct {
	ID              uuid.UUID       `json:"id"`
	CashierID       uuid.UUID       `json:"cashier_id"`
	StoreID         null.String     `json:"store_id"`
	PaymentID       null.String     `json:"payment_id"`
	CustomerID      null.String     `json:"customer_id"`
	TotalQuantity   null.Int64      `json:"total_quantity"`
	TotalUnit       null.Int64      `json:"total_unit"`
	TotalPrice      decimal.Decimal `json:"total_price"`
	TotalPriceInUSD decimal.Decimal `json:"total_price_in_usd"`
	Currency        null.String     `json:"currency"`
	UsdRate         decimal.Decimal `json:"usd_rate"`
	CreatedAt       time.Time       `json:"created_at"`
	OrderDetails    []OrderDetail   `json:"-"`
}

type OrderDetail struct {
	ID       uuid.UUID       `json:"id"`
	OrderID  uuid.UUID       `json:"order_id"`
	ItemID   null.String     `json:"item_id"`
	Quantity null.Int64      `json:"quantity"`
	Unit     null.String     `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (o *Order) FromModel(orderModel model.Order) {
	o.ID = orderModel.ID
	o.CashierID = orderModel.CashierID
	o.StoreID = orderModel.StoreID
	o.PaymentID = orderModel.PaymentID
	o.CustomerID = orderModel.CustomerID
	o.TotalQuantity = orderModel.TotalQuantity
	o.TotalUnit = orderModel.TotalUnit
	o.TotalPrice = orderModel.TotalPrice
	o.TotalPriceInUSD = orderModel.TotalPriceInUSD
	o.Currency = orderModel.Currency
	o.UsdRate = orderModel.UsdRate
	o.CreatedAt = orderModel.CreatedAt
	o.OrderDetails = make([]OrderDetail, len(orderModel.OrderDetails))

	for idx := range o.OrderDetails {
		o.OrderDetails[idx] = OrderDetail{
			ID:       orderModel.OrderDetails[idx].ID,
			OrderID:  orderModel.OrderDetails[idx].OrderID,
			ItemID:   orderModel.OrderDetails[idx].ItemID,
			Quantity: orderModel.OrderDetails[idx].Quantity,
			Unit:     orderModel.OrderDetails[idx].Unit,
			Price:    orderModel.OrderDetails[idx].Price,
		}
	}
}