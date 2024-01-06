package main

import (
	"context"
	"os"
	"skripsi-be/cmd/scheduler/task"
	"skripsi-be/config"
	"skripsi-be/connection"
	"skripsi-be/repository"
	"skripsi-be/service"

	"github.com/go-errors/errors"
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

	task := task.New(service)

	err = task.Sharding(context.Background())
	if err != nil {
		errx := err.(*errors.Error)
		log.Error().Err(err).Msg(errx.ErrorStack())
		return
	}
}
