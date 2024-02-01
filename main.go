package main

import (
	"flag"

	"skripsi-be/cmd/authservice"
	"skripsi-be/cmd/ingestionservice"
	"skripsi-be/cmd/longtermloadservice"
	"skripsi-be/cmd/middleware"
	"skripsi-be/cmd/orderservice"
	"skripsi-be/cmd/scheduler"
	"skripsi-be/cmd/servingservice"
	"skripsi-be/cmd/shardingloadservice"
	"skripsi-be/config"
	"skripsi-be/connection"
	_ "skripsi-be/pkg/errors"

	"skripsi-be/repository"
	"skripsi-be/service"

	"github.com/rs/zerolog/log"
)

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer " before the token value
func main() {
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
		connections.KafkaProducer,
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

	serviceName := flag.String("service", "", "service name")
	flag.Parse()

	mw := middleware.New(
		middleware.Config{
			JWTConfig: config.JWT,
		},
	)

	if *serviceName == "" {
		panic("service name must be not empty")
	}

	if *serviceName == "authservice" {
		err = authservice.Init(
			service,
			config.Microservice.AuthService,
			mw,
		)
		if err != nil {
			panic(err)
		}
	} else if *serviceName == "ingestionservice" {
		err = ingestionservice.Init(
			service,
			config.Microservice.IngestionService,
			mw,
		)
		if err != nil {
			panic(err)
		}
	} else if *serviceName == "longtermloadservice" {
		err = longtermloadservice.Init(
			service,
			config.Kafka,
			mw,
		)
		if err != nil {
			panic(err)
		}
	} else if *serviceName == "orderservice" {
		err = orderservice.Init(
			service,
			config.Microservice.OrderService,
			mw,
		)
		if err != nil {
			panic(err)
		}
	} else if *serviceName == "scheduler" {
		scheduler.Init(
			service,
			config.Scheduler,
		)
	} else if *serviceName == "servingservice" {
		err = servingservice.Init(
			service,
			config.Microservice.ServingService,
			mw,
		)
		if err != nil {
			panic(err)
		}
	} else if *serviceName == "shardingloadservice" {
		err = shardingloadservice.Init(
			service,
			config.Kafka,
			mw,
		)
		if err != nil {
			panic(err)
		}
	} else if *serviceName == "transformservice" {
		// err = tra.Init(
		// 	service,
		// 	config.Microservice.OrderService,
		// )
		// if err != nil {
		// 	panic(err)
		// }
	} else {
		panic("service name is not valid")
	}
}
