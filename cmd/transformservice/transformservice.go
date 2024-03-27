package transformservice

import (
	"context"
	inmiddleware "skripsi-be/cmd/middleware"
	"skripsi-be/cmd/transformservice/handler"
	"skripsi-be/config"
	"skripsi-be/pkg/http"
	"skripsi-be/service"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func Init(
	service service.Service,
	kafkaConfig config.Kafka,
	config config.TransformService,
	globalConfig config.Config,
	mw inmiddleware.Middleware,
) error {
	if globalConfig.KappaArchitecture.IsUsingKappaArchitecture {
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": kafkaConfig.Server,
			"group.id":          kafkaConfig.Group.Transform,
			"auto.offset.reset": "earliest",
		})
		if err != nil {
			return err
		}

		defer consumer.Close()

		err = handler.EventHandlerInit(kafkaConfig, consumer, service)
		if err != nil {
			return err
		}
	} else {
		e := echo.New()
		e.Use(middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: config.WhiteListAllowOrigin,
			},
		))

		handler.HTTPHandlerInit(e, service, mw, globalConfig)

		log.Info().Msgf("Starting Auth Service HTTP server on %s:%d", config.Host, config.Port)
		err := http.Start(http.HTTPServerConfig{
			Echo: e,
			Port: config.Port,
			Host: config.Host,
		})
		if err != nil {
			return err
		}

		log.Info().Msg("Starting graceful shutdown HTTP Server...")

		err = e.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Error while shutting down HTTP server")
			return err
		}

		log.Info().Msg("HTTP Server shutdown gracefully, RIP 🙏")
	}

	return nil
}
