package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (r repository) PublishTransformOrderEvent(ctx context.Context, param []params.RepositoryPublishTransformOrderEvent) error {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err)
	}

	err = r.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &r.config.TransformOrderTopic,
			Partition: kafka.PartitionAny,
		},
		Value: jsonParam,
	}, nil)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (r repository) PublishTransformOrderHTTPRequest(ctx context.Context, param []params.RepositoryPublishTransformOrderEvent) error {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err)
	}

	url := fmt.Sprintf("%s/v1/transform", r.config.TransformBaseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParam))
	if err != nil {
		return errors.Wrap(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Failed to publish to transform service, status code: %d", resp.StatusCode))
	}

	return nil
}
