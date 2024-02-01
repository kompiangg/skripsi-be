package handler

import (
	"skripsi-be/cmd/shardingloadservice/handler/loadorder"
	"skripsi-be/config"
	"skripsi-be/service"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Init(
	kafkaConfig config.Kafka,
	consumer *kafka.Consumer,
	service service.Service,
) error {
	err := consumer.Subscribe(kafkaConfig.Topic.LoadOrder, nil)
	if err != nil {
		return err
	}

	orderHandler := loadorder.New()

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			err := orderHandler.HandleLoadOrderEvent(msg)
			if err != nil {
				continue
			}

		} else {
			return err
		}
	}
}
