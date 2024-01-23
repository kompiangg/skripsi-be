package params

import (
	"context"
	"skripsi-be/type/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

type ServiceInsertOrderToShard struct {
	ID           uuid.UUID           `json:"id"`
	CashierID    uuid.UUID           `json:"cashier_id"`
	StoreID      string              `json:"store_id"`
	PaymentID    string              `json:"payment_id"`
	CustomerID   string              `json:"customer_id"`
	Currency     string              `json:"currency"`
	UsdRate      decimal.Decimal     `json:"usd_rate"`
	CreatedAt    time.Time           `json:"created_at"`
	OrderDetails []OrderDetailsShard `json:"order_details"`
}

type OrderDetailsShard struct {
	ID       uuid.UUID       `json:"id"`
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s ServiceInsertOrderToShard) ToOrderModel() model.Order {
	totalQuantity := null.Int64From(0)
	totalUnit := null.Int64From(0)
	totalPrice := decimal.NewFromInt(0)

	for _, orderDetail := range s.OrderDetails {
		totalQuantity = null.NewInt64(totalQuantity.Int64+orderDetail.Quantity, true)
		totalUnit = null.NewInt64(totalUnit.Int64+orderDetail.Quantity, true)
		totalPrice = totalPrice.Add(orderDetail.Price)
	}

	return model.Order{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         null.StringFrom(s.StoreID),
		PaymentID:       null.StringFrom(s.PaymentID),
		TotalQuantity:   totalQuantity,
		TotalUnit:       totalUnit,
		TotalPrice:      totalPrice,
		TotalPriceInUSD: totalPrice.Mul(s.UsdRate),
		CustomerID:      null.StringFrom(s.CustomerID),
		Currency:        null.StringFrom(s.Currency),
		UsdRate:         s.UsdRate,
		CreatedAt:       s.CreatedAt,
		OrderDetails:    s.ToOrderDetailModel(),
	}
}

func (s ServiceInsertOrderToShard) ToOrderDetailModel() []model.OrderDetail {
	orderDetails := make([]model.OrderDetail, 0)

	for _, orderDetail := range s.OrderDetails {
		orderDetails = append(orderDetails, model.OrderDetail{
			ID:       orderDetail.ID,
			OrderID:  s.ID,
			ItemID:   null.StringFrom(orderDetail.ItemID),
			Quantity: null.NewInt64(orderDetail.Quantity, true),
			Unit:     null.StringFrom(orderDetail.Unit),
			Price:    orderDetail.Price,
		})
	}

	return orderDetails
}

type ServiceInsertOrdersToShardParam []ServiceInsertOrderToShard

func (s ServiceInsertOrdersToShardParam) Validate(ctx context.Context) error {
	for _, order := range s {
		err := validation.ValidateStructWithContext(ctx, &order,
			validation.Field(&order.ID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.StoreID, validation.Required, validation.NotNil),
			validation.Field(&order.PaymentID, validation.Required, validation.NotNil),
			validation.Field(&order.CustomerID, validation.Required, validation.NotNil),
			validation.Field(&order.Currency, validation.Required, validation.NotNil),
			validation.Field(&order.UsdRate, validation.Required, validation.NotNil),
			validation.Field(&order.CreatedAt, validation.Required, validation.NotNil),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

type ServiceInsertOrderToLongTermParam struct {
	ID           uuid.UUID       `json:"id"`
	CashierID    uuid.UUID       `json:"cashier_id"`
	StoreID      string          `json:"store_id"`
	PaymentID    string          `json:"payment_id"`
	CustomerID   string          `json:"customer_id"`
	Currency     string          `json:"currency"`
	UsdRate      decimal.Decimal `json:"usd_rate"`
	CreatedAt    time.Time       `json:"created_at"`
	OrderDetails []OrderDetailsLongTerm
}

type OrderDetailsLongTerm struct {
	ID       uuid.UUID       `json:"id"`
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s ServiceInsertOrderToLongTermParam) ToOrderModel() model.Order {
	totalQuantity := null.Int64From(0)
	totalUnit := null.Int64From(0)
	totalPrice := decimal.NewFromInt(0)

	for _, orderDetail := range s.OrderDetails {
		totalQuantity = null.NewInt64(totalQuantity.Int64+orderDetail.Quantity, true)
		totalUnit = null.NewInt64(totalUnit.Int64+orderDetail.Quantity, true)
		totalPrice = totalPrice.Add(orderDetail.Price)
	}

	return model.Order{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         null.StringFrom(s.StoreID),
		PaymentID:       null.StringFrom(s.PaymentID),
		CustomerID:      null.StringFrom(s.CustomerID),
		Currency:        null.StringFrom(s.Currency),
		TotalQuantity:   totalQuantity,
		TotalUnit:       totalUnit,
		TotalPrice:      totalPrice,
		TotalPriceInUSD: totalPrice.Mul(s.UsdRate),
		UsdRate:         s.UsdRate,
		CreatedAt:       s.CreatedAt,
	}
}

func (s ServiceInsertOrderToLongTermParam) ToOrderDetailModel() []model.OrderDetail {
	orderDetails := make([]model.OrderDetail, 0)
	for _, orderDetail := range s.OrderDetails {
		orderDetails = append(orderDetails, model.OrderDetail{
			ID:       orderDetail.ID,
			OrderID:  s.ID,
			ItemID:   null.StringFrom(orderDetail.ItemID),
			Quantity: null.NewInt64(orderDetail.Quantity, true),
			Unit:     null.StringFrom(orderDetail.Unit),
			Price:    orderDetail.Price,
		})
	}

	return orderDetails
}

type ServiceInsertOrdersToLongTermParam []ServiceInsertOrderToLongTermParam

func (s ServiceInsertOrdersToLongTermParam) Validate(ctx context.Context) error {
	for _, order := range s {
		err := validation.ValidateStructWithContext(ctx, &order,
			validation.Field(&order.ID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.StoreID, validation.Required, validation.NotNil),
			validation.Field(&order.PaymentID, validation.Required, validation.NotNil),
			validation.Field(&order.CustomerID, validation.Required, validation.NotNil),
			validation.Field(&order.Currency, validation.Required, validation.NotNil),
			validation.Field(&order.UsdRate, validation.Required, validation.NotNil),
			validation.Field(&order.CreatedAt, validation.Required, validation.NotNil),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
