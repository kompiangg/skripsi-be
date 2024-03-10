package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/model"
	"skripsi-be/type/params"
)

func (s service) TransformOrder(ctx context.Context, param []params.ServiceTransformOrder) error {
	for _, v := range param {
		err := v.Validate()
		if err != nil {
			return errors.Wrap(err)
		}
	}

	mapCurrency := make(map[string]model.Currency)
	for _, v := range param {
		if _, ok := mapCurrency[v.Currency]; !ok {
			mapCurrency[v.Currency] = model.Currency{}
		}
	}

	for currency := range mapCurrency {
		usdRate, err := s.currencyRepo.FindByBaseAndQuote(ctx, "USD", currency)
		if err != nil {
			return errors.Wrap(err)
		}

		mapCurrency[currency] = usdRate
	}

	repoParam := make([]params.RepositoryPublishLoadOrderEvent, len(param))
	for i, v := range param {
		repoParam[i] = v.TransformOrder(mapCurrency[v.Currency].Rate)
	}

	err := s.publisherRepo.PublishLoadOrderEvent(ctx, repoParam)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
