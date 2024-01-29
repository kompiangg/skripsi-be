package ingestionservice

import (
	"fmt"
	"skripsi-be/config"
	"skripsi-be/service"

	inmiddleware "skripsi-be/cmd/middleware"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Init(
	service service.Service,
	kafkaConfig config.Kafka,
	mw inmiddleware.Middleware,
) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Server,
		"client.id":         "ingestionservice",
		"acks":              "all",
	})
	if err != nil {
		return err
	}

	defer producer.Close()

	for {
		var word string
		fmt.Scanf("%s", &word)

		producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &kafkaConfig.Topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}

	return nil
}
