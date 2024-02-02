package shardingloadservice

import (
	"skripsi-be/config"
	"skripsi-be/pkg/errors"
	"skripsi-be/service"

	inmiddleware "skripsi-be/cmd/middleware"
	"skripsi-be/cmd/shardingloadservice/handler"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Init(
	service service.Service,
	kafkaConfig config.Kafka,
	mw inmiddleware.Middleware,
) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Server,
		"group.id":          kafkaConfig.Group.Shard,
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
