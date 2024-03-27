package handler

import (
	"fmt"
	"skripsi-be/cmd/middleware"
	"skripsi-be/cmd/shardingloadservice/handler/loadorder"
	"skripsi-be/config"
	"skripsi-be/pkg/errors"
	"skripsi-be/service"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func EventHandlerInit(
	kafkaConfig config.Kafka,
	consumer *kafka.Consumer,
	service service.Service,
) error {
	err := consumer.Subscribe(kafkaConfig.Topic.LoadOrder, nil)
	if err != nil {
		return errors.Wrap(err)
	}

	orderHandler := loadorder.New(
		service.Order,
	)

	log.Info().Msg("Listening to load order event...")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Println("Received message", string(msg.Value))
			err := orderHandler.HandleLoadOrderEvent(msg)
			if err != nil {
				log.Error().Err(err).Msg("Error while handling load order event")
				continue
			}

		} else {
			log.Error().Err(err).Msg("Error while reading message")
			return errors.Wrap(err)
		}
	}
}

func HTTPHandlerInit(
	echo *echo.Echo,
	service service.Service,
	middleware middleware.Middleware,
	config config.Config,
) {
	loadorder.InitHTTPHandler(echo, middleware, config, service.Order)
}
