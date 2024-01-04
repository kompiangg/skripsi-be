package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"skripsi-be/config"
	"skripsi-be/pkg/errors"
	"skripsi-be/tools/migration/entity"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	migrationFlag := entity.MigrationFlag{}

	migrationFlag.Operation = flag.String("operation", "", "operation to perform (new, up, down, or drop)")
	migrationFlag.ConnectionType = flag.String("type", "", "connection type (longterm, shard, and general)")
	migrationFlag.TableName = flag.String("tableName", "", "table name")
	flag.Parse()

	err := migrationFlag.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg("invalid migration flag")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	shardMigration := make([]*migrate.Migrate, len(cfg.ShardingDatabase.Shards))
	for i, shard := range cfg.ShardingDatabase.Shards {
		shardMigration[i], err = migrate.New(
			"file://./migration/shard",
			shard.URIConnection,
		)
		if err != nil {
			log.Fatal().Err(errors.New(err)).Stack().Msg("failed to create shard migration")
		}
	}

	longtermMigration, err := migrate.New(
		"file://./migration/longterm",
		cfg.LongTermDatabase.URIConnection,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create longterm migration")
	}

	generalMigration, err := migrate.New(
		"file://./migration/general",
		cfg.GeneralDatabase.URIConnection,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create general migration")
	}

	if *migrationFlag.Operation == "up" {
		if *migrationFlag.ConnectionType == "shard" {
			for _, shard := range shardMigration {
				err = shard.Up()
				if errors.Is(err, migrate.ErrNoChange) {
					log.Info().Msg("no change in longterm migration")
				} else if err != nil {
					log.Fatal().Err(err).Msg("failed to run shard migration")
				}
			}
		} else if *migrationFlag.ConnectionType == "longterm" {
			err = longtermMigration.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info().Msg("no change in longterm migration")
			} else if err != nil {
				log.Fatal().Err(err).Msg("failed to run longterm migration")
			}
		} else if *migrationFlag.ConnectionType == "general" {
			err = generalMigration.Up()
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info().Msg("no change in general migration")
			} else if err != nil {
				log.Fatal().Err(err).Msg("failed to run general migration")
			}
		}
	} else if *migrationFlag.Operation == "down" {
		if *migrationFlag.ConnectionType == "shard" {
			for _, shard := range shardMigration {
				err = shard.Down()
				if errors.Is(err, migrate.ErrNoChange) {
					log.Info().Msg("no change in longterm migration")
				} else if err != nil {
					log.Fatal().Err(err).Msg("failed to run shard migration")
				}
			}
		} else if *migrationFlag.ConnectionType == "longterm" {
			err = longtermMigration.Down()
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info().Msg("no change in longterm migration")
			} else if err != nil {
				log.Fatal().Err(err).Msg("failed to run longterm migration")
			}
		} else if *migrationFlag.ConnectionType == "general" {
			err = generalMigration.Down()
			if errors.Is(err, migrate.ErrNoChange) {
				log.Info().Msg("no change in general migration")
			} else if err != nil {
				log.Fatal().Err(err).Msg("failed to run general migration")
			}
		}
	} else if *migrationFlag.Operation == "drop" {
		if *migrationFlag.ConnectionType == "shard" {
			for _, shard := range shardMigration {
				err = shard.Drop()
				if err != nil {
					log.Fatal().Err(err).Msg("failed to run shard migration")
				}
			}
		} else if *migrationFlag.ConnectionType == "longterm" {
			err = longtermMigration.Drop()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to run longterm migration")
			}
		} else if *migrationFlag.ConnectionType == "general" {
			err = generalMigration.Drop()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to run general migration")
			}
		}
	} else if *migrationFlag.Operation == "new" {
		if *migrationFlag.ConnectionType == "shard" {
			cmd := exec.CommandContext(
				context.Background(),
				"sh", "-c",
				fmt.Sprintf("migrate create -ext sql -dir ./migration/shard %s", *migrationFlag.TableName))

			err := cmd.Run()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create new shard migration")
			}
		} else if *migrationFlag.ConnectionType == "longterm" {
			cmd := exec.CommandContext(
				context.Background(),
				"sh", "-c",
				fmt.Sprintf("migrate create -ext sql -dir ./migration/longterm %s", *migrationFlag.TableName))

			err := cmd.Run()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create new longterm migration")
			}
		} else if *migrationFlag.ConnectionType == "general" {
			cmd := exec.CommandContext(
				context.Background(),
				"sh", "-c",
				fmt.Sprintf("migrate create -ext sql -dir ./migration/general %s", *migrationFlag.TableName))

			err := cmd.Run()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to create new general migration")
			}
		}
	}
}
