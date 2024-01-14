package params

import (
	"context"
	"skripsi-be/type/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v9"
)

type ServiceInsertOrderToShardParam struct {
	ID         uuid.UUID     `json:"id"`
	ItemID     null.String   `json:"item_id"`
	StoreID    null.String   `json:"store_id"`
	CustomerID null.String   `json:"customer_id"`
	CashierID  uuid.NullUUID `json:"cashier_id"`
	Unit       null.String   `json:"unit"`
	CreatedAt  time.Time     `json:"created_at"`
	Price      float64       `json:"price"`
	TotalPrice float64       `json:"total_price"`
	PaymentID  null.String   `json:"payment_id"`
	Quantity   int           `json:"quantity"`
}

func (s ServiceInsertOrderToShardParam) ToOrderModel() model.Order {
	return model.Order{
		ID:         s.ID,
		CashierID:  s.CashierID,
		CustomerID: s.CustomerID,
		ItemID:     s.ItemID,
		StoreID:    s.StoreID,
		PaymentID:  s.PaymentID,
		Quantity:   null.NewInt(s.Quantity, true),
		Unit:       null.NewString(s.Unit.String, true),
		Price:      null.NewFloat64(s.Price, true),
		TotalPrice: null.NewFloat64(s.TotalPrice, true),
		CreatedAt:  s.CreatedAt,
	}
}

type ServiceInsertOrdersToShardParam []ServiceInsertOrderToShardParam

func (s ServiceInsertOrdersToShardParam) Validate(ctx context.Context) error {
	for _, order := range s {
		err := validation.ValidateStructWithContext(ctx, &order,
			validation.Field(&order.ID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.ItemID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.StoreID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.CustomerID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.PaymentID, validation.Required, validation.NotNil),
			validation.Field(&order.Quantity, validation.Required, validation.NotNil),
			validation.Field(&order.Unit, validation.Required, validation.NotNil),
			validation.Field(&order.Price, validation.Required, validation.NotNil),
			validation.Field(&order.TotalPrice, validation.Required, validation.NotNil),
			validation.Field(&order.CreatedAt, validation.Required, validation.NotNil),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

type ServiceInsertOrderToLongTermParam struct {
	ID         uuid.UUID     `json:"id"`
	CashierID  uuid.NullUUID `json:"cashier_id"`
	ItemID     null.String   `json:"item_id"`
	StoreID    null.String   `json:"store_id"`
	CustomerID null.String   `json:"customer_id"`
	Unit       null.String   `json:"unit"`
	CreatedAt  time.Time     `json:"created_at"`
	Price      float64       `json:"price"`
	TotalPrice float64       `json:"total_price"`
	PaymentID  null.String   `json:"payment_id"`
	Quantity   int           `json:"quantity"`
}

func (s ServiceInsertOrderToLongTermParam) ToOrderModel() model.Order {
	return model.Order{
		ID:         s.ID,
		CashierID:  s.CashierID,
		CustomerID: s.CustomerID,
		ItemID:     s.ItemID,
		StoreID:    s.StoreID,
		PaymentID:  s.PaymentID,
		Quantity:   null.NewInt(s.Quantity, true),
		Unit:       null.NewString(s.Unit.String, true),
		Price:      null.NewFloat64(s.Price, true),
		TotalPrice: null.NewFloat64(s.TotalPrice, true),
		CreatedAt:  s.CreatedAt,
	}
}

type ServiceInsertOrdersToLongTermParam []ServiceInsertOrderToLongTermParam

func (s ServiceInsertOrdersToLongTermParam) Validate(ctx context.Context) error {
	for _, order := range s {
		err := validation.ValidateStructWithContext(ctx, &order,
			validation.Field(&order.ID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.ItemID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.StoreID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.CustomerID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.PaymentID, validation.Required, validation.NotNil),
			validation.Field(&order.Quantity, validation.Required, validation.NotNil),
			validation.Field(&order.Unit, validation.Required, validation.NotNil),
			validation.Field(&order.Price, validation.Required, validation.NotNil),
			validation.Field(&order.TotalPrice, validation.Required, validation.NotNil),
			validation.Field(&order.CreatedAt, validation.Required, validation.NotNil),
		)

		if err != nil {
			return err
		}
	}

	return nil
}
