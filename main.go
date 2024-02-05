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
	"skripsi-be/cmd/transformservice"
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
		log.Info().Msg("Starting auth service...")
		err = authservice.Init(
			service,
			config.Microservice.AuthService,
			mw,
		)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down auth service...")
	} else if *serviceName == "ingestionservice" {
		log.Info().Msg("Starting ingestion service...")
		err = ingestionservice.Init(
			service,
			config.Microservice.IngestionService,
			mw,
		)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down ingestion service...")
	} else if *serviceName == "longtermloadservice" {
		log.Info().Msg("Starting long term load service...")
		err = longtermloadservice.Init(
			service,
			config.Kafka,
			mw,
		)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down long term load service...")
	} else if *serviceName == "orderservice" {
		log.Info().Msg("Starting order service...")
		err = orderservice.Init(
			service,
			config.Microservice.OrderService,
			mw,
		)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down order service...")
	} else if *serviceName == "scheduler" {
		log.Info().Msg("Starting scheduler...")
		scheduler.Init(
			service,
			config.Scheduler,
		)
		log.Info().Msg("Shutting down scheduler...")
	} else if *serviceName == "servingservice" {
		log.Info().Msg("Starting serving service...")
		err = servingservice.Init(
			service,
			config.Microservice.ServingService,
			config,
			mw,
		)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down serving service...")
	} else if *serviceName == "shardingloadservice" {
		log.Info().Msg("Starting sharding load service...")
		err = shardingloadservice.Init(
			service,
			config.Kafka,
			mw,
		)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down sharding load service...")
	} else if *serviceName == "transformservice" {
		log.Info().Msg("Starting transform service...")
		err = transformservice.Init(service, config.Kafka)
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Shutting down transform service...")
	} else {
		panic("service name is not valid")
	}
}
