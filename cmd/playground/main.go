package main

import (
	"context"
	"os"
	"skripsi-be/config"
	"skripsi-be/connection"
	"skripsi-be/repository"
	"skripsi-be/service"
	"skripsi-be/type/params"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v9"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config, err := config.Load("./etc/config.yaml")
	if err != nil {
		panic(err)
	}

	connections, err := connection.New(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Info().Msg("Closing all connections...")

		err = connections.Close()
		if err != nil {
			panic(err)
		}

		log.Info().Msg("All connections closed, RIP 🙏")
	}()

	repository, err := repository.New(
		config,
		connections.LongTermDatabase,
		connections.ShardingDatabase,
		connections.Redis,
	)
	if err != nil {
		panic(err)
	}

	service, err := service.New(
		repository,
		config,
	)
	if err != nil {
		panic(err)
	}

	// task := task.New(service)

	orders := make(params.ServiceInsertOrdersToShardParam, 0)

	// 0
	dateTime, err := time.Parse(time.RFC3339, "2023-12-31T12:00:00+08:00")
	if err != nil {
		panic(err)
	}
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	// 1
	dateTime, _ = time.Parse(time.RFC3339, "2023-12-30T12:00:00+08:00")
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	// 1
	dateTime, _ = time.Parse(time.RFC3339, "2023-12-24T12:00:00+08:00")
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	// 2
	dateTime, _ = time.Parse(time.RFC3339, "2023-12-25T12:00:00+08:00")
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	// 2
	dateTime, _ = time.Parse(time.RFC3339, "2023-12-01T12:00:00+08:00")
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	// 3
	dateTime, _ = time.Parse(time.RFC3339, "2023-11-01T12:00:00+08:00")
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	// not inserted
	dateTime, _ = time.Parse(time.RFC3339, "2023-01-01T12:00:00+08:00")
	orders = append(orders, params.ServiceInsertOrderToShardParam{
		ID:         uuid.New(),
		ItemID:     uuid.New(),
		StoreID:    uuid.New(),
		CustomerID: uuid.New(),
		Unit:       null.String{String: "kg", Valid: true},
		CreatedAt:  dateTime,
		Price:      10000,
		TotalPrice: 10000,
		PaymentID:  1,
		Quantity:   1,
	})

	err = service.Order.InsertToShard(context.Background(), orders)
	if err != nil {
		panic(err)
	}
}
