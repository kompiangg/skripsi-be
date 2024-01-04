package main

import (
	"flag"

	webservice "skripsi-be/cmd/orderservice"
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
	config, err := config.Load()
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

	serviceName := flag.String("service", "", "service name")
	flag.Parse()

	if *serviceName == "" {
		panic("service name must be not empty")
	}

	if *serviceName == "orderservice" {
		err = webservice.Init(
			service,
			config.Microservice.OrderService,
		)
		if err != nil {
			panic(err)
		}
	}
}
