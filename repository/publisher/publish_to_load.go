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

func (r repository) PublishLoadOrderEvent(ctx context.Context, param []params.RepositoryPublishLoadOrderEvent) error {
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

func (r repository) PublishLongtermLoadOrderHTTPRequest(ctx context.Context, param []params.RepositoryPublishLoadOrderEvent) error {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err)
	}

	url := fmt.Sprintf("%s/v1/load", r.config.LongtermLoadBaseURl)

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
		return errors.New(fmt.Sprintf("Failed to publish to longterm load service, status code: %d", resp.StatusCode))
	}

	return nil
}

func (r repository) PublishShardingLoadOrderHTTPRequest(ctx context.Context, param []params.RepositoryPublishLoadOrderEvent) error {
	jsonParam, err := json.Marshal(param)
	if err != nil {
		return errors.Wrap(err)
	}

	url := fmt.Sprintf("%s/v1/load", r.config.ShardingLoadBaseURL)

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
		return errors.New(fmt.Sprintf("Failed to publish to sharding load service, status code: %d", resp.StatusCode))
	}

	return nil
}
