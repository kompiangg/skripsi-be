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
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v9"
)

func LoadOrderData(config config.Config, connections connection.Connection, service service.Service) error {
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

	iteration := 1
	for i := 0; i < iteration; i++ {
		orders, err := loadOrder("./dataset/fact_table.csv", stores, storeCashier, customers, payments, items, config)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		longTermParams := make(params.ServiceInsertOrdersToLongTermParam, 0)
		shardParams := make(params.ServiceInsertOrdersToShardParam, 0)

		for _, order := range orders {
			longTermParams = append(longTermParams, params.ServiceInsertOrderToLongTermParam{
				ID:         order.ID,
				ItemID:     order.ItemID,
				StoreID:    order.StoreID,
				CashierID:  order.CashierID,
				CustomerID: order.CustomerID,
				Unit:       order.Unit,
				CreatedAt:  order.CreatedAt,
				Price:      order.Price.Float64,
				TotalPrice: order.TotalPrice.Float64,
				PaymentID:  order.PaymentID,
				Quantity:   order.Quantity.Int,
			})

			shardParams = append(shardParams, params.ServiceInsertOrderToShardParam{
				ID:         order.ID,
				ItemID:     order.ItemID,
				StoreID:    order.StoreID,
				CashierID:  order.CashierID,
				CustomerID: order.CustomerID,
				Unit:       order.Unit,
				CreatedAt:  order.CreatedAt,
				Price:      order.Price.Float64,
				TotalPrice: order.TotalPrice.Float64,
				PaymentID:  order.PaymentID,
				Quantity:   order.Quantity.Int,
			})
		}

		err = service.Order.InsertToLongTerm(ctx, longTermParams)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}

		err = service.Order.InsertToShard(ctx, shardParams)
		if err != nil {
			return errors.Wrap(err, constant.SkipErrorParameter)
		}
	}

	return nil
}

func loadOrder(path string, stores []model.Store, storeCashier map[string][]model.Cashier, customers []model.Customer, payments []model.PaymentType, items []model.Item, config config.Config) ([]model.Order, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, constant.SkipErrorParameter)
	}

	lenCustomers := len(customers)
	lenPayments := len(payments)
	lenItems := len(items)
	lenStores := len(stores)

	orders := make([]model.Order, 0)
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

		orderQty, err := strconv.Atoi(record[5])
		if err != nil {
			return nil, errors.Wrap(err, constant.SkipErrorParameter)
		}

		orderPrice, err := strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, errors.Wrap(err, constant.SkipErrorParameter)
		}

		totalPrice, err := strconv.ParseFloat(record[8], 64)
		if err != nil {
			return nil, errors.Wrap(err, constant.SkipErrorParameter)
		}

		start := time.Date(2023, 6, 1, 0, 0, 0, 0, time.FixedZone("Asia/Makassar", 8*60*60))
		end := config.Date.Now()
		diff := end.Unix() - start.Unix()
		randSeconds := r.Int63n(diff)
		randTime := start.Add(time.Duration(randSeconds) * time.Second)

		uuidExists := map[string]bool{}
		uuidOrderID := uuid.New()

		for uuidExists[uuidOrderID.String()] {
			uuidOrderID = uuid.New()
		}

		orders = append(orders, model.Order{
			ID:         uuidOrderID,
			ItemID:     null.NewString(items[randItem].ID, true),
			StoreID:    null.NewString(stores[randStore].ID, true),
			CashierID:  uuid.NullUUID{UUID: cashierStore[randCashier].ID, Valid: true},
			CustomerID: null.NewString(customers[randCustomer].ID, true),
			PaymentID:  null.NewString(payments[randPayment].ID, true),
			Quantity:   null.NewInt(orderQty, true),
			Unit:       null.NewString(record[6], true),
			Price:      null.NewFloat64(orderPrice, true),
			TotalPrice: null.NewFloat64(totalPrice, true),
			CreatedAt:  randTime,
		})
	}

	return orders, nil
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
