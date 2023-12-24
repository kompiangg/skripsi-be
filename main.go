package main

import (
	"flag"

	"skripsi-be/cmd/webservice"
	"skripsi-be/config"
	"skripsi-be/connection"
	_ "skripsi-be/pkg/errors"

	"skripsi-be/repository"
	"skripsi-be/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer " before the token value
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	connections, err := connection.New(config)
	if err != nil {
		panic(err)
	}

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

	if *serviceName == "webservice" {
		err = webservice.Init(
			service,
			config.Microservice.WebService,
		)
		if err != nil {
			panic(err)
		}
	}

	log.Info().Msg("Closing all connections...")

	err = connections.Close()
	if err != nil {
		panic(err)
	}

	log.Info().Msg("All connections closed, RIP 🙏")
}
