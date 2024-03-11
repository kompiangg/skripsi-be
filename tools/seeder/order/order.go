package order

import (
	"context"
	"encoding/csv"
	"math/rand"
	"os"
	"os/signal"
	"skripsi-be/config"
	"skripsi-be/connection"
	"skripsi-be/lib/encodelib"
	"skripsi-be/repository"
	"skripsi-be/service"
	"skripsi-be/type/constant"
	"skripsi-be/type/model"
	"skripsi-be/type/params"
	"strconv"
	"syscall"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"github.com/volatiletech/null/v9"
)

func LoadOrderData(config config.Config, connections connection.Connection, service service.Service, repository repository.Repository, iteration int) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
		<-sigterm
		log.Info().Msg("Received terminate or interrupt signal, canceling context...")
		cancel()
	}()

	defer cancel()

	payments, err := findAllPayment(ctx, connections.GeneralDatabase)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	customers, err := findAllCustomers(ctx, connections.GeneralDatabase)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	stores, err := findAllStores(ctx, connections.GeneralDatabase)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	items, err := findAllItems(ctx, connections.GeneralDatabase)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	cashiers, err := findAllCashier(ctx, connections.GeneralDatabase)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	storeCashier := make(map[string][]model.Cashier)

	for _, store := range stores {
		storeCashier[store.ID] = make([]model.Cashier, 0)

		for _, cashier := range cashiers {
			if cashier.StoreID == store.ID {
				storeCashier[store.ID] = append(storeCashier[store.ID], cashier)
			}
		}
	}

	usdRate, err := repository.Currency.FindByBaseAndQuote(ctx, "USD", "BDT")
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	for i := 0; i < iteration; i++ {
		log.Info().Msgf("Iteration %d", i+1)

		orders, detailOrders, err := loadOrder("./dataset/fact_table.csv", stores, storeCashier, customers, payments, items, config, usdRate.Rate)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		longTermParams := make(params.ServiceInsertOrdersToLongTermParam, 0)
		shardParams := make(params.ServiceInsertOrdersToShardParam, 0)

		orderIDMap := map[string]int{}

		for idx, order := range orders {
			if _, ok := orderIDMap[order.ID]; !ok {
				orderIDMap[order.ID] = idx
			}

			orderID, err := ulid.Parse(order.ID)
			if err != nil {
				return errors.Wrap(err, constant.SkipErrorParameter)
			}

			longTermParams = append(longTermParams, params.ServiceInsertOrderToLongTermParam{
				ID:         orderID.String(),
				CashierID:  order.CashierID,
				StoreID:    order.StoreID.String,
				PaymentID:  order.PaymentID.String,
				CustomerID: order.CustomerID,
				Currency:   order.Currency.String,
				UsdRate:    order.UsdRate,
				CreatedAt:  order.CreatedAt,
			})

			shardParams = append(shardParams, params.ServiceInsertOrderToShard{
				ID:         orderID.String(),
				CashierID:  order.CashierID,
				StoreID:    order.StoreID.String,
				PaymentID:  order.PaymentID.String,
				CustomerID: order.CustomerID,
				Currency:   order.Currency.String,
				UsdRate:    order.UsdRate,
				CreatedAt:  order.CreatedAt,
			})
		}

		for _, orderDetail := range detailOrders {
			idx := orderIDMap[orderDetail.OrderID]

			longTermParams[idx].OrderDetails = append(longTermParams[idx].OrderDetails, params.OrderDetailsLongTerm{
				ID:       orderDetail.ID,
				ItemID:   orderDetail.ItemID.String,
				Quantity: orderDetail.Quantity.Int64,
				Unit:     orderDetail.Unit.String,
				Price:    orderDetail.Price,
			})

			shardParams[idx].OrderDetails = append(shardParams[idx].OrderDetails, params.OrderDetailsShard{
				ID:       orderDetail.ID,
				ItemID:   orderDetail.ItemID.String,
				Quantity: orderDetail.Quantity.Int64,
				Unit:     orderDetail.Unit.String,
				Price:    orderDetail.Price,
			})
		}

		err = service.Order.InsertToLongTermSeeder(ctx, longTermParams)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		err = service.Order.InsertToShardSeeder(ctx, shardParams)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}
	}

	return nil
}

func loadOrder(path string, stores []model.Store, storeCashier map[string][]model.Cashier, customers []model.Customer, payments []model.PaymentType, items []model.Item, config config.Config, usdRate decimal.Decimal) ([]model.Order, []model.OrderDetail, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	lenCustomers := len(customers)
	lenPayments := len(payments)
	lenItems := len(items)
	lenStores := len(stores)

	orders := make([]model.Order, 0)
	orderDetails := make([]model.OrderDetail, 0)
	idMap := map[string]string{}

	for idx, record := range records {
		if idx == 0 {
			continue
		}

		name := encodelib.SanitizeUTF8String(record[1])
		if name == "" {
			continue
		}

		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)

		randStore := r.Intn(lenStores)
		randCustomer := r.Intn(lenCustomers)
		randPayment := r.Intn(lenPayments)
		randItem := r.Intn(lenItems)

		cashierStore := storeCashier[stores[randStore].ID]
		randCashier := r.Intn(len(cashierStore))

		orderQty, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
		}

		orderPrice, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
		}

		start := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		end := config.Date.Now()
		diff := end.Unix() - start.Unix()
		randSeconds := r.Int63n(diff)
		randTime := start.Add(time.Duration(randSeconds) * time.Second)

		idExists := map[string]bool{}
		uuidOrderDetailsExists := map[string]bool{}

		if orderID, ok := idMap[record[2]]; ok {
			uuidOrderID := ulid.MustParse(orderID)
			uuidOrderDetailsID := uuid.New()

			for uuidOrderDetailsExists[uuidOrderDetailsID.String()] {
				uuidOrderDetailsID = uuid.New()
			}

			orderDetails = append(orderDetails, model.OrderDetail{
				ID:       uuidOrderDetailsID,
				OrderID:  uuidOrderID.String(),
				ItemID:   null.NewString(items[randItem].ID, true),
				Quantity: null.NewInt64(orderQty, true),
				Unit:     null.NewString(record[6], true),
				Price:    decimal.NewFromFloat(orderPrice),
			})
		} else {
			entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
			ms := ulid.Timestamp(randTime)
			orderID, err := ulid.New(ms, entropy)
			if err != nil {
				return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
			}

			for idExists[orderID.String()] {
				orderID, err = ulid.New(ms, entropy)
				if err != nil {
					return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
				}
			}

			uuidOrderDetailsID := uuid.New()
			for uuidOrderDetailsExists[uuidOrderDetailsID.String()] {
				uuidOrderDetailsID = uuid.New()
			}

			idMap[record[2]] = orderID.String()

			orders = append(orders, model.Order{
				ID:            orderID.String(),
				CashierID:     cashierStore[randCashier].ID,
				StoreID:       null.StringFrom(stores[randStore].ID),
				PaymentID:     null.StringFrom(payments[randPayment].ID),
				CustomerID:    null.StringFrom(customers[randCustomer].ID),
				TotalQuantity: null.NewInt64(0, true),
				TotalPrice:    decimal.NewFromInt(0),
				Currency:      null.StringFrom("AFN"),
				UsdRate:       usdRate,
				CreatedAt:     randTime,
			})

			orderDetails = append(orderDetails, model.OrderDetail{
				ID:       uuidOrderDetailsID,
				OrderID:  orderID.String(),
				ItemID:   null.NewString(items[randItem].ID, true),
				Quantity: null.NewInt64(orderQty, true),
				Unit:     null.NewString(record[6], true),
				Price:    decimal.NewFromFloat(orderPrice),
			})
		}

	}

	return orders, orderDetails, nil
}

func findAllCashier(ctx context.Context, generalDB *sqlx.DB) ([]model.Cashier, error) {
	q := `
		SELECT
			id,
			account_id,
			store_id,
			name,
			created_at
		FROM
			cashiers
	`

	cashiers := make([]model.Cashier, 0)
	err := generalDB.SelectContext(ctx, &cashiers, q)
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	return cashiers, nil
}

func findAllStores(ctx context.Context, generalDB *sqlx.DB) ([]model.Store, error) {
	q := `
		SELECT
			id,
			region,
			district,
			sub_district,
			created_at
		FROM
			stores
	`

	stores := make([]model.Store, 0)
	err := generalDB.SelectContext(ctx, &stores, q)
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	return stores, nil
}

func findAllCustomers(ctx context.Context, generalDB *sqlx.DB) ([]model.Customer, error) {
	q := `
		SELECT
			id,
			name,
			contact,
			created_at
		FROM
			customers
	`

	customers := make([]model.Customer, 0)
	err := generalDB.SelectContext(ctx, &customers, q)
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	return customers, nil
}

func findAllPayment(ctx context.Context, generalDB *sqlx.DB) ([]model.PaymentType, error) {
	q := `
		SELECT
			id,
			"type",
			bank,
			created_at
		FROM
			payment_types
	`

	paymentTypes := make([]model.PaymentType, 0)
	err := generalDB.SelectContext(ctx, &paymentTypes, q)
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	return paymentTypes, nil
}

func findAllItems(ctx context.Context, generalDB *sqlx.DB) ([]model.Item, error) {
	q := `
		SELECT
			id,
			name,
			created_at
		FROM
			items
	`

	items := make([]model.Item, 0)
	err := generalDB.SelectContext(ctx, &items, q)
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	return items, nil
}
