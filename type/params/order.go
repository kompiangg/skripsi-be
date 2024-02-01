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
	ID              uuid.UUID           `json:"id"`
	CashierID       uuid.UUID           `json:"cashier_id"`
	StoreID         string              `json:"store_id"`
	PaymentID       string              `json:"payment_id"`
	CustomerID      string              `json:"customer_id"`
	TotalQuantity   int64               `json:"total_quantity"`
	TotalUnit       int64               `json:"total_unit"`
	Currency        string              `json:"currency"`
	TotalPrice      decimal.Decimal     `json:"total_price"`
	TotalPriceInUSD decimal.Decimal     `json:"total_price_in_usd"`
	UsdRate         decimal.Decimal     `json:"usd_rate"`
	CreatedAt       time.Time           `json:"created_at"`
	OrderDetails    []OrderDetailsShard `json:"order_details"`
}

type OrderDetailsShard struct {
	ID       uuid.UUID       `json:"id"`
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s ServiceInsertOrderToShard) ToOrderModelInSeeder() model.Order {
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

func (s ServiceInsertOrderToShard) ToOrderModel() model.Order {
	return model.Order{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         null.StringFrom(s.StoreID),
		PaymentID:       null.StringFrom(s.PaymentID),
		TotalQuantity:   null.NewInt64(s.TotalQuantity, true),
		TotalUnit:       null.Int64From(s.TotalUnit),
		TotalPrice:      s.TotalPrice,
		TotalPriceInUSD: s.TotalPriceInUSD,
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
	ID              uuid.UUID              `json:"id"`
	CashierID       uuid.UUID              `json:"cashier_id"`
	StoreID         string                 `json:"store_id"`
	PaymentID       string                 `json:"payment_id"`
	CustomerID      string                 `json:"customer_id"`
	TotalQuantity   int64                  `json:"total_quantity"`
	TotalUnit       int64                  `json:"total_unit"`
	Currency        string                 `json:"currency"`
	TotalPrice      decimal.Decimal        `json:"total_price"`
	TotalPriceInUSD decimal.Decimal        `json:"total_price_in_usd"`
	UsdRate         decimal.Decimal        `json:"usd_rate"`
	CreatedAt       time.Time              `json:"created_at"`
	OrderDetails    []OrderDetailsLongTerm `json:"order_details"`
}

type OrderDetailsLongTerm struct {
	ID       uuid.UUID       `json:"id"`
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s ServiceInsertOrderToLongTermParam) ToOrderModelInSeeder() model.Order {
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
		OrderDetails:    s.ToOrderDetailModel(),
	}
}

func (s ServiceInsertOrderToLongTermParam) ToOrderModel() model.Order {
	return model.Order{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         null.StringFrom(s.StoreID),
		PaymentID:       null.StringFrom(s.PaymentID),
		CustomerID:      null.StringFrom(s.CustomerID),
		Currency:        null.StringFrom(s.Currency),
		TotalQuantity:   null.Int64From(s.TotalQuantity),
		TotalUnit:       null.Int64From(s.TotalUnit),
		TotalPrice:      s.TotalPrice,
		TotalPriceInUSD: s.TotalPriceInUSD,
		UsdRate:         s.UsdRate,
		CreatedAt:       s.CreatedAt,
		OrderDetails:    s.ToOrderDetailModel(),
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

type FindOrderService struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type TransformOrderEventService struct {
	ID           uuid.UUID                        `json:"id"`
	CashierID    uuid.UUID                        `json:"cashier_id"`
	StoreID      string                           `json:"store_id"`
	PaymentID    string                           `json:"payment_id"`
	CustomerID   string                           `json:"customer_id"`
	Currency     string                           `json:"currency"`
	CreatedAt    time.Time                        `json:"created_at"`
	OrderDetails []PublishOrderDetailEventService `json:"order_details"`
}

type PublishOrderDetailEventService struct {
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s TransformOrderEventService) Validate() error {
	err := validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required, validation.NotNil, is.UUIDv4),
		validation.Field(&s.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
		validation.Field(&s.StoreID, validation.Required, validation.NotNil),
		validation.Field(&s.PaymentID, validation.Required, validation.NotNil),
		validation.Field(&s.CustomerID, validation.Required, validation.NotNil),
		validation.Field(&s.Currency, validation.Required, validation.NotNil),
		validation.Field(&s.CreatedAt, validation.Required, validation.NotNil),
		validation.Field(&s.OrderDetails, validation.Required, validation.NotNil, validation.Length(1, 0)),
	)
	if err != nil {
		return err
	}

	for _, orderDetail := range s.OrderDetails {
		err := validation.ValidateStruct(&orderDetail,
			validation.Field(&orderDetail.ItemID, validation.Required, validation.NotNil),
			validation.Field(&orderDetail.Quantity, validation.Required, validation.NotNil, validation.NotIn(0)),
			validation.Field(&orderDetail.Unit, validation.Required, validation.NotNil, validation.NotIn(0)),
			validation.Field(&orderDetail.Price, validation.Required, validation.NotNil, validation.NotIn(0)),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s TransformOrderEventService) ToPublishOrderEventRepository() PublishTransformOrderEventRepository {
	res := PublishTransformOrderEventRepository{
		ID:           s.ID,
		CashierID:    s.CashierID,
		StoreID:      s.StoreID,
		PaymentID:    s.PaymentID,
		CustomerID:   s.CustomerID,
		Currency:     s.Currency,
		CreatedAt:    s.CreatedAt,
		OrderDetails: []PublishOrderDetailEventService{},
	}

	for _, orderDetail := range s.OrderDetails {
		res.OrderDetails = append(res.OrderDetails, PublishOrderDetailEventService{
			ItemID:   orderDetail.ItemID,
			Quantity: orderDetail.Quantity,
			Unit:     orderDetail.Unit,
			Price:    orderDetail.Price,
		})
	}

	return res
}

type PublishTransformOrderEventRepository struct {
	ID           uuid.UUID                        `json:"id"`
	CashierID    uuid.UUID                        `json:"cashier_id"`
	StoreID      string                           `json:"store_id"`
	PaymentID    string                           `json:"payment_id"`
	CustomerID   string                           `json:"customer_id"`
	Currency     string                           `json:"currency"`
	CreatedAt    time.Time                        `json:"created_at"`
	OrderDetails []PublishOrderDetailEventService `json:"order_details"`
}

type PublishTransformOrderDetailEventRepository struct {
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}
