package general

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"os/signal"
	"skripsi-be/connection"
	"skripsi-be/lib/encodelib"
	"skripsi-be/type/constant"
	"skripsi-be/type/model"
	"strconv"
	"syscall"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

func LoadGeneralDatabaseData(connections connection.Connection) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
		<-sigterm
		log.Info().Msg("Received terminate or interrupt signal, canceling context...")
		cancel()
	}()

	defer cancel()

	tx, err := connections.GeneralDatabase.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	defer tx.Commit()

	log.Info().Msg("Start load customer data from csv...")
	customers, err := loadCustomer("./dataset/customer_dim.csv")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish load customer data from csv...")

	log.Info().Msg("Start insert customer data to database...")
	err = insertCustomer(ctx, tx, customers)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish insert customer data to database...")
	customers = nil

	log.Info().Msg("Start load store data from csv...")
	stores, err := loadStore("./dataset/store_dim.csv")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish load store data from csv...")

	log.Info().Msg("Start insert store data to database...")
	err = insertStore(ctx, tx, stores)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish insert store data to database...")

	log.Info().Msg("Start create cashier data...")
	cashiers, accounts, err := createCashier(ctx, 2, stores)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish create cashier data...")

	log.Info().Msg("Start insert cashier data to database...")
	err = insertCashier(ctx, tx, cashiers)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish insert cashier data to database...")

	log.Info().Msg("Start insert account data to database...")
	err = insertAccount(ctx, tx, accounts)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish insert account data to database...")

	log.Info().Msg("Start load payment data from csv...")
	payments, err := loadPayment("./dataset/Trans_dim.csv")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish load payment data from csv...")

	log.Info().Msg("Start insert payment data to database...")
	err = insertPayment(ctx, tx, payments)
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Info().Msg("Finish insert payment data to database...")
	payments = nil

	log.Info().Msg("Start load item data from csv...")
	items, err := loadItem("./dataset/item_dim.csv")
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish load item data from csv...")

	log.Info().Msg("Start insert item data to database...")
	err = insertItem(ctx, tx, items)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}
	log.Info().Msg("Finish insert item data to database...")

	return nil
}

func loadCustomer(path string) ([]model.Customer, error) {
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

	customers := make([]model.Customer, 0)
	for idx, record := range records {
		if idx == 0 {
			continue
		}

		name := encodelib.SanitizeUTF8String(record[1])
		if name == "" {
			continue
		}

		customers = append(customers, model.Customer{
			ID:        record[0],
			Name:      name,
			Contact:   record[2],
			CreatedAt: time.Now(),
		})
	}

	return customers, nil
}

func insertCustomer(ctx context.Context, tx *sqlx.Tx, customers []model.Customer) error {
	q := `
		INSERT INTO customers (
			id,
			name,
			contact,
			created_at
		) VALUES (
			:id,
			:name,
			:contact,
			:created_at
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	for _, customer := range customers {
		_, err := stmt.ExecContext(ctx, customer)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, constant.SkipErrorParameter)
		}
	}

	return nil
}

func loadStore(path string) ([]model.Store, error) {
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

	stores := make([]model.Store, 0)
	for idx, record := range records {
		if idx == 0 {
			continue
		}

		stores = append(stores, model.Store{
			ID:          record[0],
			Nation:      "BANGLADESH",
			Region:      record[1],
			District:    record[2],
			SubDistrict: record[3],
			Currency:    "BDT",
			CreatedAt:   time.Now(),
		})
	}

	return stores, nil
}

func insertStore(ctx context.Context, tx *sqlx.Tx, stores []model.Store) error {
	q := `
		INSERT INTO stores (
			id,
			region,
			district,
			sub_district,
			created_at
		) VALUES (
			:id,
			:region,
			:district,
			:sub_district,
			:created_at
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, store := range stores {
		_, err := stmt.ExecContext(ctx, store)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func loadPayment(path string) ([]model.PaymentType, error) {
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

	payments := make([]model.PaymentType, 0)
	for idx, record := range records {
		if idx == 0 {
			continue
		}

		payments = append(payments, model.PaymentType{
			ID:        record[0],
			Type:      record[1],
			Bank:      record[2],
			CreatedAt: time.Now(),
		})
	}

	return payments, nil
}

func insertPayment(ctx context.Context, tx *sqlx.Tx, payments []model.PaymentType) error {
	q := `
		INSERT INTO payment_types (
			id,
			type,
			bank,
			created_at
		) VALUES (
			:id,
			:type,
			:bank,
			:created_at
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, payment := range payments {
		_, err := stmt.ExecContext(ctx, payment)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func createCashier(ctx context.Context, cashierEachStore int, stores []model.Store) ([]model.Cashier, []model.Account, error) {
	cashiers := make([]model.Cashier, 0)
	accounts := make([]model.Account, 0)

	for _, store := range stores {
		for i := 0; i < cashierEachStore; i++ {
			cashierName := fmt.Sprintf("cashier-%d", len(cashiers)+1)
			password, err := bcrypt.GenerateFromPassword([]byte(cashierName), bcrypt.DefaultCost)
			if err != nil {
				return nil, nil, errors.Wrap(err, constant.SkipErrorParameter)
			}

			account := model.Account{
				ID:        uuid.New(),
				Username:  cashierName,
				Password:  string(password),
				CreatedAt: time.Now(),
			}
			accounts = append(accounts, account)

			cashiers = append(cashiers, model.Cashier{
				ID:        uuid.New(),
				AccountID: account.ID,
				StoreID:   store.ID,
				Name:      cashierName,
				CreatedAt: time.Now(),
			})
		}
	}

	return cashiers, accounts, nil
}

func insertCashier(ctx context.Context, tx *sqlx.Tx, cashiers []model.Cashier) error {
	q := `
		INSERT INTO cashiers (
			id,
			account_id,
			store_id,
			name,
			created_at
		) VALUES (
			:id,
			:account_id,
			:store_id,
			:name,
			:created_at
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, cashier := range cashiers {
		_, err := stmt.ExecContext(ctx, cashier)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func insertAccount(ctx context.Context, tx *sqlx.Tx, accounts []model.Account) error {
	q := `
		INSERT INTO accounts (
			id,
			username,
			password,
			created_at
		) VALUES (
			:id,
			:username,
			:password,
			:created_at
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, account := range accounts {
		_, err := stmt.ExecContext(ctx, account)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func loadItem(path string) ([]model.Item, error) {
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

	items := make([]model.Item, 0)
	for idx, record := range records {
		if idx == 0 {
			continue
		}

		price, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, errors.Wrap(err, constant.SkipErrorParameter)
		}

		name := encodelib.SanitizeUTF8String(record[1])
		if name == "" {
			continue
		}

		items = append(items, model.Item{
			ID:            record[0],
			Name:          name,
			Desc:          record[2],
			Price:         decimal.NewFromFloat(price),
			OriginCountry: record[4],
			Supplier:      record[5],
			Unit:          record[6],
			CreatedAt:     time.Now(),
		})
	}

	return items, nil
}

func insertItem(ctx context.Context, tx *sqlx.Tx, items []model.Item) error {
	q := `
		INSERT INTO items (
			id,
			"name",
			"desc",
			price,
			origin_country,
			supplier,
			unit,
			created_at
		) VALUES (
			:id,
			:name,
			:desc,
			:price,
			:origin_country,
			:supplier,
			:unit,
			:created_at
		)
	`

	stmt, err := tx.PrepareNamedContext(ctx, q)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range items {
		_, err := stmt.ExecContext(ctx, item)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}
