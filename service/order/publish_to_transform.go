package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"
)

func (s service) PublishToTransform(ctx context.Context, param []params.TransformOrderEventService) error {
	for _, v := range param {
		err := v.Validate()
		if err != nil {
			return errors.Wrap(err)
		}
	}

	repoParam := make([]params.PublishTransformOrderEventRepository, len(param))
	for i, v := range param {
		repoParam[i] = v.ToPublishOrderEventRepository()
	}

	err := s.publisherRepo.PublishLoadOrderEvent(ctx, repoParam)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
