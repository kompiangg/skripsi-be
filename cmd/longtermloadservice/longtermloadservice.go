package longtermloadservice

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
	fmt.Println(kafkaConfig.Group.LongTerm)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Server,
		"group.id":          kafkaConfig.Group.LongTerm,
		"auto.offset.reset": "latest",
	})
	if err != nil {
		return err
	}

	fmt.Println("Starting Sharding Load Service...")
	defer consumer.Close()
	fmt.Println(kafkaConfig.Topic)

	err = consumer.Subscribe(kafkaConfig.Topic, nil)
	if err != nil {
		return err
	}

	fmt.Println("Starting Sharding Load Service...")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	return nil
}
