package main

import (
	"flag"
	"os"
	"skripsi-be/config"
	"skripsi-be/connection"
	"skripsi-be/repository"
	"skripsi-be/service"
	"skripsi-be/tools/seeder/entity"
	"skripsi-be/tools/seeder/general"
	"skripsi-be/tools/seeder/order"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	seederFlag := entity.SeederFlag{}

	seederFlag.ConnectionType = flag.String("type", "", "connection type (else, and general)")
	flag.Parse()

	err := seederFlag.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid migration flag")
	}

	cfg, err := config.Load("./etc/config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	connections, err := connection.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create connection")
	}
	defer func() {
		log.Info().Msg("Closing all connections...")

		err = connections.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to close connection")
		}

		log.Info().Msg("All connections closed, RIP 🙏")
	}()

	repo, err := repository.New(cfg, connections.LongTermDatabase, connections.GeneralDatabase, connections.ShardingDatabase, connections.KafkaProducer)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create repository")
	}

	svc, err := service.New(repo, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create service")
	}

	if *seederFlag.ConnectionType == "general" {
		err := general.LoadGeneralDatabaseData(connections)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to load data to general database, %s", err.(*errors.Error).ErrorStack())
		}
	} else if *seederFlag.ConnectionType == "else" {
		err := order.LoadOrderData(cfg, connections, svc, repo, 10)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to load data to else database, %s", err.(*errors.Error).ErrorStack())
		}
	} else {
		log.Fatal().Msg("invalid connection type")
	}

	log.Info().Msg("Seeder finished")
}
