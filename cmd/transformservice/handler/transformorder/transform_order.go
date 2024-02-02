package transformorder

import (
	"context"
	"encoding/json"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (h handler) HandleTransformOrderEvent(msg *kafka.Message) error {
	var req []params.ServiceTransformOrder
	err := json.Unmarshal(msg.Value, &req)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.orderService.TransformOrder(context.Background(), req)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
