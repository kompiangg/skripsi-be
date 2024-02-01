package handler

import (
	"skripsi-be/cmd/transformservice/handler/transformorder"
	"skripsi-be/config"
	"skripsi-be/service"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Init(
	kafkaConfig config.Kafka,
	consumer *kafka.Consumer,
	service service.Service,
) error {
	err := consumer.Subscribe(kafkaConfig.Topic.TransformOrder, nil)
	if err != nil {
		return err
	}

	transformHandler := transformorder.New()

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			err := transformHandler.HandleTransformOrderEvent(msg)
			if err != nil {
				continue
			}

		} else {
			return err
		}
	}
}
