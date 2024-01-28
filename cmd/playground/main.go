package main

import (
	"context"
	"fmt"
	"os"
	"skripsi-be/cmd/scheduler/task"
	"skripsi-be/config"
	"skripsi-be/connection"
	"skripsi-be/repository"
	"skripsi-be/service"
	"skripsi-be/type/params"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
		connections.GeneralDatabase,
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

	_ = task.New(service)

	paramStart, err := time.Parse("2006-01-02", "2023-12-01")
	if err != nil {
		panic(err)
	}

	paramEnd, err := time.Parse("2006-01-02", "2023-12-31")
	if err != nil {
		panic(err)
	}

	startProcess := time.Now()

	shardOrders, err := service.Order.FindOrder(context.Background(), params.FindOrderService{
		StartDate: paramStart,
		EndDate:   paramEnd,
	})
	if err != nil {
		panic(err)
	}

	var shardingDatabaseRuntime float64 = float64(time.Since(startProcess).Milliseconds())

	detailOrdersCount := 0
	for _, order := range shardOrders {
		detailOrdersCount += len(order.OrderDetails)
	}

	fmt.Println("Sharding Database")
	fmt.Printf("Querying %d row in %vms\n", detailOrdersCount, shardingDatabaseRuntime)

	shardOrders = nil
	startProcess = time.Now()

	ctx := context.WithValue(context.Background(), "isUsingSharding", false)
	shardOrders, err = service.Order.FindOrder(ctx, params.FindOrderService{
		StartDate: paramStart,
		EndDate:   paramEnd,
	})
	if err != nil {
		panic(err)
	}

	var longtermDatabaseRuntime float64 = float64(time.Since(startProcess).Milliseconds())

	detailOrdersCount = 0
	for _, order := range shardOrders {
		detailOrdersCount += len(order.OrderDetails)
	}

	fmt.Println("Longterm Database")
	fmt.Printf("Querying %d row in %vms\n", detailOrdersCount, longtermDatabaseRuntime)

	speedIncrease := ((longtermDatabaseRuntime - shardingDatabaseRuntime) / longtermDatabaseRuntime) * 100
	fmt.Printf("Sharding database is %f%% faster than Longterm database\n", speedIncrease)
}
