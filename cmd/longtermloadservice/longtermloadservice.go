package longtermloadservice

import (
	"skripsi-be/cmd/longtermloadservice/handler"
	"skripsi-be/config"
	"skripsi-be/pkg/errors"
	"skripsi-be/service"

	inmiddleware "skripsi-be/cmd/middleware"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Init(
	service service.Service,
	kafkaConfig config.Kafka,
	mw inmiddleware.Middleware,
) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Server,
		"group.id":          kafkaConfig.Group.LongTerm,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return errors.Wrap(err)
	}

	defer consumer.Close()

	err = handler.Init(kafkaConfig, consumer, service)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
