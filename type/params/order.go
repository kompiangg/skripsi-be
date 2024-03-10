package params

import (
	"context"
	"math/rand"
	"skripsi-be/type/model"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

type ServiceInsertOrderToShard struct {
	ID              string              `json:"id"`
	CashierID       uuid.UUID           `json:"cashier_id"`
	StoreID         string              `json:"store_id"`
	PaymentID       string              `json:"payment_id"`
	CustomerID      null.String         `json:"customer_id"`
	TotalQuantity   int64               `json:"total_quantity"`
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
	totalPrice := decimal.NewFromInt(0)

	for _, orderDetail := range s.OrderDetails {
		totalQuantity = null.NewInt64(totalQuantity.Int64+orderDetail.Quantity, true)
		totalPrice = totalPrice.Add(orderDetail.Price)
	}

	return model.Order{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         null.StringFrom(s.StoreID),
		PaymentID:       null.StringFrom(s.PaymentID),
		TotalQuantity:   totalQuantity,
		TotalPrice:      totalPrice,
		TotalPriceInUSD: totalPrice.Div(s.UsdRate),
		CustomerID:      s.CustomerID,
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
		TotalPrice:      s.TotalPrice,
		TotalPriceInUSD: s.TotalPriceInUSD,
		CustomerID:      s.CustomerID,
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
			validation.Field(&order.ID, validation.Required, validation.NotNil),
			validation.Field(&order.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.StoreID, validation.Required, validation.NotNil),
			validation.Field(&order.PaymentID, validation.Required, validation.NotNil),
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
	ID              string                 `json:"id"`
	CashierID       uuid.UUID              `json:"cashier_id"`
	StoreID         string                 `json:"store_id"`
	PaymentID       string                 `json:"payment_id"`
	CustomerID      null.String            `json:"customer_id"`
	TotalQuantity   int64                  `json:"total_quantity"`
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
	totalPrice := decimal.NewFromInt(0)

	for _, orderDetail := range s.OrderDetails {
		totalQuantity = null.NewInt64(totalQuantity.Int64+orderDetail.Quantity, true)
		totalPrice = totalPrice.Add(orderDetail.Price)
	}

	return model.Order{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         null.StringFrom(s.StoreID),
		PaymentID:       null.StringFrom(s.PaymentID),
		CustomerID:      s.CustomerID,
		Currency:        null.StringFrom(s.Currency),
		TotalQuantity:   totalQuantity,
		TotalPrice:      totalPrice,
		TotalPriceInUSD: totalPrice.Div(s.UsdRate),
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
		CustomerID:      s.CustomerID,
		Currency:        null.StringFrom(s.Currency),
		TotalQuantity:   null.Int64From(s.TotalQuantity),
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
			validation.Field(&order.ID, validation.Required, validation.NotNil),
			validation.Field(&order.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
			validation.Field(&order.StoreID, validation.Required, validation.NotNil),
			validation.Field(&order.PaymentID, validation.Required, validation.NotNil),
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
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`

	CashierID uuid.NullUUID `json:"-"`
}

func (s FindOrderService) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.StartDate, validation.Required, validation.NotNil),
		validation.Field(&s.EndDate, validation.Required, validation.NotNil),
	)
}

type ServiceIngestionOrder struct {
	StoreID      string                        `json:"store_id"`
	PaymentID    string                        `json:"payment_id"`
	CustomerID   null.String                   `json:"customer_id"`
	Currency     string                        `json:"currency"`
	PaymentDate  time.Time                     `json:"payment_date"`
	OrderDetails []ServiceIngestionOrderDetail `json:"order_details"`

	CashierID uuid.UUID `json:"-"`
}

type ServiceIngestionOrderDetail struct {
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
	Unit     string          `json:"unit"`
}

func (s ServiceIngestionOrder) Validate() error {
	err := validation.ValidateStruct(&s,
		validation.Field(&s.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
		validation.Field(&s.StoreID, validation.Required, validation.NotNil),
		validation.Field(&s.PaymentID, validation.Required, validation.NotNil),
		validation.Field(&s.Currency, validation.Required, validation.NotNil),
		validation.Field(&s.PaymentDate, validation.Required, validation.NotNil),
		validation.Field(&s.OrderDetails, validation.Required, validation.NotNil, validation.Length(1, 0)),
	)
	if err != nil {
		return err
	}

	for _, orderDetail := range s.OrderDetails {
		err := validation.ValidateStruct(&orderDetail,
			validation.Field(&orderDetail.ItemID, validation.Required, validation.NotNil),
			validation.Field(&orderDetail.Quantity, validation.Required, validation.NotNil, validation.NotIn(0)),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ServiceIngestionOrder) ToRepositoryPublishTransformOrderEvent() RepositoryPublishTransformOrderEvent {
	id, err := ulid.New(uint64(s.PaymentDate.Unix()), rand.New(rand.NewSource(time.Now().Unix())))
	if err != nil {
		panic(err)
	}

	res := RepositoryPublishTransformOrderEvent{
		ID:           id.String(),
		CashierID:    s.CashierID,
		StoreID:      s.StoreID,
		PaymentID:    s.PaymentID,
		CustomerID:   s.CustomerID,
		Currency:     s.Currency,
		CreatedAt:    s.PaymentDate,
		OrderDetails: make([]RepositoryPublishTransformOrderDetailEvent, 0),
	}

	return res
}

type RepositoryPublishTransformOrderEvent struct {
	ID           string                                       `json:"id"`
	CashierID    uuid.UUID                                    `json:"cashier_id"`
	StoreID      string                                       `json:"store_id"`
	PaymentID    string                                       `json:"payment_id"`
	CustomerID   null.String                                  `json:"customer_id"`
	Currency     string                                       `json:"currency"`
	CreatedAt    time.Time                                    `json:"created_at"`
	OrderDetails []RepositoryPublishTransformOrderDetailEvent `json:"order_details"`
}

type RepositoryPublishTransformOrderDetailEvent struct {
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s ServiceIngestionOrderDetail) ToRepositoryPublishTransformOrderDetailEvent(item model.Item) RepositoryPublishTransformOrderDetailEvent {
	return RepositoryPublishTransformOrderDetailEvent{
		ItemID:   s.ItemID,
		Quantity: s.Quantity,
		Unit:     s.Unit,
		Price:    s.Price,
	}
}

type ServiceTransformOrder struct {
	ID           string                        `json:"id"`
	CashierID    uuid.UUID                     `json:"cashier_id"`
	StoreID      string                        `json:"store_id"`
	PaymentID    string                        `json:"payment_id"`
	CustomerID   null.String                   `json:"customer_id"`
	Currency     string                        `json:"currency"`
	CreatedAt    time.Time                     `json:"created_at"`
	OrderDetails []ServiceTransformOrderDetail `json:"order_details"`
}

type ServiceTransformOrderDetail struct {
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s ServiceTransformOrder) Validate() error {
	err := validation.ValidateStruct(&s,
		validation.Field(&s.CashierID, validation.Required, validation.NotNil, is.UUIDv4),
		validation.Field(&s.StoreID, validation.Required, validation.NotNil),
		validation.Field(&s.PaymentID, validation.Required, validation.NotNil),
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
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s ServiceTransformOrder) TransformOrder(usdRate decimal.Decimal) RepositoryPublishLoadOrderEvent {
	var totalQuantity int64
	totalPrice := decimal.NewFromInt(0)

	for _, orderDetail := range s.OrderDetails {
		totalQuantity += orderDetail.Quantity
		totalPrice = totalPrice.Add(orderDetail.Price)
	}

	res := RepositoryPublishLoadOrderEvent{
		ID:              s.ID,
		CashierID:       s.CashierID,
		StoreID:         s.StoreID,
		PaymentID:       s.PaymentID,
		CustomerID:      s.CustomerID,
		Currency:        s.Currency,
		TotalQuantity:   totalQuantity,
		TotalPrice:      totalPrice,
		TotalPriceInUSD: totalPrice.Div(usdRate),
		UsdRate:         usdRate,
		CreatedAt:       s.CreatedAt,
		OrderDetails:    []RepositoryPublishLoadOrderDetailEvent{},
	}

	for _, orderDetail := range s.OrderDetails {
		res.OrderDetails = append(res.OrderDetails, orderDetail.toRepositoryPublishLoadOrderDetailEvent())
	}

	return res
}

func (s ServiceTransformOrderDetail) toRepositoryPublishLoadOrderDetailEvent() RepositoryPublishLoadOrderDetailEvent {
	return RepositoryPublishLoadOrderDetailEvent{
		ID:       uuid.New(),
		ItemID:   s.ItemID,
		Quantity: s.Quantity,
		Unit:     s.Unit,
		Price:    s.Price,
	}
}

type RepositoryPublishLoadOrderEvent struct {
	ID              string                                  `json:"id"`
	CashierID       uuid.UUID                               `json:"cashier_id"`
	StoreID         string                                  `json:"store_id"`
	PaymentID       string                                  `json:"payment_id"`
	CustomerID      null.String                             `json:"customer_id"`
	TotalQuantity   int64                                   `json:"total_quantity"`
	Currency        string                                  `json:"currency"`
	TotalPrice      decimal.Decimal                         `json:"total_price"`
	TotalPriceInUSD decimal.Decimal                         `json:"total_price_in_usd"`
	UsdRate         decimal.Decimal                         `json:"usd_rate"`
	CreatedAt       time.Time                               `json:"created_at"`
	OrderDetails    []RepositoryPublishLoadOrderDetailEvent `json:"order_details"`
}

type RepositoryPublishLoadOrderDetailEvent struct {
	ID       uuid.UUID       `json:"id"`
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

type FindOrderDetailsService struct {
	OrderID string `param:"id"`
}

func (s FindOrderDetailsService) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.OrderID, validation.Required, validation.NotNil),
	)
}
