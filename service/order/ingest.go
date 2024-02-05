package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
)

func (s service) IngestOrder(ctx context.Context, param []params.ServiceIngestionOrder) ([]result.ServiceIngestOrder, error) {
	for _, v := range param {
		err := v.Validate()
		if err != nil {
			return nil, errors.Wrap(err)
		}

		if v.CreatedAt.After(s.config.Date.Now()) {
			return nil, errors.Wrap(errors.ErrDataParamMustNotBeforeCurrentTime)
		}
	}

	res := make([]result.ServiceIngestOrder, len(param))
	repoParam := make([]params.RepositoryPublishTransformOrderEvent, len(param))
	for i, v := range param {
		repoParam[i] = v.ToRepositoryPublishTransformOrderEvent()
		res[i].FromParamServiceIngestionOrder(v, repoParam[i].ID)
	}

	err := s.publisherRepo.PublishTransformOrderEvent(ctx, repoParam)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return res, nil
}
