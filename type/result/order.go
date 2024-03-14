package result

import (
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
	"skripsi-be/type/params"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

type Order struct {
	ID              ulid.ULID       `json:"id"`
	CashierID       uuid.UUID       `json:"cashier_id"`
	StoreID         null.String     `json:"store_id"`
	PaymentID       null.String     `json:"payment_id"`
	CustomerID      null.String     `json:"customer_id"`
	TotalQuantity   null.Int64      `json:"total_quantity"`
	TotalPrice      decimal.Decimal `json:"total_price"`
	TotalPriceInUSD decimal.Decimal `json:"total_price_in_usd"`
	Currency        null.String     `json:"currency"`
	UsdRate         decimal.Decimal `json:"usd_rate"`
	CreatedAt       time.Time       `json:"created_at"`
	OrderDetails    []OrderDetail   `json:"order_details"`
}

type OrderDetail struct {
	ID       uuid.UUID       `json:"id"`
	OrderID  ulid.ULID       `json:"order_id"`
	ItemID   null.String     `json:"item_id"`
	Quantity null.Int64      `json:"quantity"`
	Unit     null.String     `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (o *Order) FromModel(orderModel model.Order) error {
	var err error
	o.ID, err = ulid.Parse(orderModel.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	o.CashierID = orderModel.CashierID
	o.StoreID = orderModel.StoreID
	o.PaymentID = orderModel.PaymentID
	o.CustomerID = orderModel.CustomerID
	o.TotalQuantity = orderModel.TotalQuantity
	o.TotalPrice = orderModel.TotalPrice
	o.TotalPriceInUSD = orderModel.TotalPriceInUSD
	o.Currency = orderModel.Currency
	o.UsdRate = orderModel.UsdRate
	o.CreatedAt = orderModel.CreatedAt
	o.OrderDetails = make([]OrderDetail, len(orderModel.OrderDetails))

	for idx := range o.OrderDetails {
		o.OrderDetails[idx] = OrderDetail{
			ID:       orderModel.OrderDetails[idx].ID,
			OrderID:  o.ID,
			ItemID:   orderModel.OrderDetails[idx].ItemID,
			Quantity: orderModel.OrderDetails[idx].Quantity,
			Unit:     orderModel.OrderDetails[idx].Unit,
			Price:    orderModel.OrderDetails[idx].Price,
		}
	}

	return nil
}

type ServiceIngestOrder struct {
	ID           string                     `json:"id"`
	CashierID    uuid.UUID                  `json:"cashier_id"`
	PaymentID    string                     `json:"payment_id"`
	CustomerID   null.String                `json:"customer_id"`
	CreatedAt    time.Time                  `json:"created_at"`
	OrderDetails []ServiceIngestOrderDetail `json:"order_details"`
}

type ServiceIngestOrderDetail struct {
	ItemID   string          `json:"item_id"`
	Quantity int64           `json:"quantity"`
	Unit     string          `json:"unit"`
	Price    decimal.Decimal `json:"price"`
}

func (s *ServiceIngestOrder) FromParamServiceIngestionOrder(param params.ServiceIngestionOrder, id string) {
	s.ID = id
	s.CashierID = param.CashierID
	s.PaymentID = param.PaymentID
	s.CustomerID = param.CustomerID
	s.CreatedAt = param.PaymentDate
	s.OrderDetails = make([]ServiceIngestOrderDetail, len(param.OrderDetails))

	for idx := range s.OrderDetails {
		s.OrderDetails[idx] = ServiceIngestOrderDetail{
			ItemID:   param.OrderDetails[idx].ItemID,
			Quantity: param.OrderDetails[idx].Quantity,
			Unit:     param.OrderDetails[idx].Unit,
			Price:    param.OrderDetails[idx].Price,
		}
	}
}

type OrderBriefInformation struct {
	ID              ulid.ULID       `json:"id"`
	CashierID       uuid.UUID       `json:"cashier_id"`
	StoreID         null.String     `json:"store_id"`
	PaymentID       null.String     `json:"payment_id"`
	CustomerID      null.String     `json:"customer_id"`
	TotalQuantity   null.Int64      `json:"total_quantity"`
	TotalPrice      decimal.Decimal `json:"total_price"`
	TotalPriceInUSD decimal.Decimal `json:"total_price_in_usd"`
	Currency        null.String     `json:"currency"`
	UsdRate         decimal.Decimal `json:"usd_rate"`
	CreatedAt       time.Time       `json:"created_at"`
}

func (o *OrderBriefInformation) FromModel(orderModel model.Order) error {
	var err error
	o.ID, err = ulid.Parse(orderModel.ID)
	if err != nil {
		return errors.Wrap(err)
	}

	o.CashierID = orderModel.CashierID
	o.StoreID = orderModel.StoreID
	o.PaymentID = orderModel.PaymentID
	o.CustomerID = orderModel.CustomerID
	o.TotalQuantity = orderModel.TotalQuantity
	o.TotalPrice = orderModel.TotalPrice
	o.TotalPriceInUSD = orderModel.TotalPriceInUSD
	o.Currency = orderModel.Currency
	o.UsdRate = orderModel.UsdRate
	o.CreatedAt = orderModel.CreatedAt

	return nil
}

type GetAggregateOrderService struct {
	TopProducts []GetAggregateOrderTopProductService `json:"top_products"`
	Chart       []GetAggregateOrderChartService      `json:"chart"`

	ItemSoldTotalQuantity int64           `json:"item_sold_total_quantity"`
	ItemSoldTotalPrice    decimal.Decimal `json:"item_sold_total_price"`

	TotalCustomerOrderQuantity    int64 `json:"total_customer_order_quantity"`
	TotalNotCustomerOrderQuantity int64 `json:"total_not_customer_order_quantity"`
}

type GetAggregateOrderTopProductService struct {
	ItemID                string `json:"item_id"`
	ItemName              string `json:"item_name"`
	ItemSoldTotalQuantity int64  `json:"item_sold_total_quantity"`
}

type GetAggregateOrderChartService struct {
	TotalOrderQuantity int64           `json:"total_order_quantity"`
	TotalOrderPrice    decimal.Decimal `json:"total_order_price"`
	Date               time.Time       `json:"date"`
}
