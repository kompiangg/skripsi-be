package task

import (
	"context"
	"skripsi-be/type/constant"

	"github.com/go-errors/errors"
)

func (t task) Sharding(ctx context.Context) error {
	err := t.service.Order.MoveDataThroughShard(ctx)
	if err != nil {
		return errors.Wrap(err, constant.SkipErrorParameter)
	}

	return nil
}
