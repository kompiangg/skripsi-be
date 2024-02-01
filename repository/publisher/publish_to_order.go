package publisher

import (
	"context"
	"encoding/json"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (r repository) PublishLoadOrderEvent(ctx context.Context, param []params.PublishTransformOrderEventRepository) error {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err)
	}

	err = r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &r.config.LoadOrderTopic,
			Partition: kafka.PartitionAny,
		},
		Value: jsonParam,
	}, nil)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
