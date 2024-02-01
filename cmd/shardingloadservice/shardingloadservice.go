package shardingloadservice

import (
	"skripsi-be/config"
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
		"auto.offset.reset": "latest",
	})
	if err != nil {
		return err
	}

	defer consumer.Close()

	err = handler.Init(kafkaConfig, consumer, service)
	if err != nil {
		return err
	}

	return nil
}
