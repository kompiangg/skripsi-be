package handler

import (
	"skripsi-be/cmd/transformservice/handler/transformorder"
	"skripsi-be/config"
	"skripsi-be/pkg/errors"
	"skripsi-be/service"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

func Init(
	kafkaConfig config.Kafka,
	consumer *kafka.Consumer,
	service service.Service,
) error {
	err := consumer.Subscribe(kafkaConfig.Topic.TransformOrder, nil)
	if err != nil {
		return errors.Wrap(err)
	}

	transformHandler := transformorder.New(
		service.Order,
	)

	log.Info().Msg("Listening to transform order event...")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			err := transformHandler.HandleTransformOrderEvent(msg)
			if err != nil {
				log.Error().Err(err).Msg("Error while handling transform order event")
				continue
			}

		} else {
			log.Error().Err(err).Msg("Error while reading message")
			return err
		}
	}
}
