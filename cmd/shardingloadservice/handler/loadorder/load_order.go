package loadorder

import (
	"context"
	"encoding/json"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (h handler) HandleLoadOrderEvent(msg *kafka.Message) error {
	var req []params.ServiceInsertOrderToShard
	err := json.Unmarshal(msg.Value, &req)
	if err != nil {
		return errors.Wrap(err)
	}

	err = h.orderService.InsertToShard(context.Background(), req)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
