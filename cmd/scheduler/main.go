package scheduler

import (
	"context"
	"skripsi-be/cmd/scheduler/task"
	"skripsi-be/config"
	"skripsi-be/service"

	"github.com/go-errors/errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func Init(
	service service.Service,
	config config.Scheduler,
) {
	c := cron.New()
	task := task.New(service)

	if config.MoveShardingData.Enable {
		_, err := c.AddFunc(config.MoveShardingData.Duration, func() {
			log.Info().Msg("Running sharding task...")

			err := task.Sharding(context.Background())
			if err != nil {
				log.Error().Err(err.(*errors.Error)).Msg("Error while running sharding task")
				return
			}

			log.Info().Msg("Sharding task done")
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Error while adding sharding task")
		}
	}
}
