package publisher

import (
	"context"
	"skripsi-be/type/params"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Repository interface {
	PublishLoadOrderEvent(ctx context.Context, param []params.RepositoryPublishLoadOrderEvent) error
	PublishTransformOrderEvent(ctx context.Context, param []params.RepositoryPublishTransformOrderEvent) error
}

type Config struct {
	LoadOrderTopic      string
	TransformOrderTopic string
}

type repository struct {
	config   Config
	producer *kafka.Producer
}

func New(
	config Config,
	producer *kafka.Producer,
) Repository {
	return repository{
		config:   config,
		producer: producer,
	}
}
