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
	for idx, v := range param {
		cashier, err := s.cashierRepo.FindByID(ctx, v.CashierID)
		if err != nil {
			return errors.Wrap(err)
		}

		store, err := s.storeRepo.FindByID(ctx, cashier.StoreID)
		if err != nil {
			return errors.Wrap(err)
		}

		param[idx].StoreID = store.ID
		param[idx].Currency = store.Currency

		if _, ok := mapCurrency[store.Currency]; !ok {
			mapCurrency[store.Currency] = model.Currency{}
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

	if s.config.KappaArchitecture.IsUsingKappaArchitecture {
		err := s.publisherRepo.PublishLoadOrderEvent(ctx, repoParam)
		if err != nil {
			return errors.Wrap(err)
		}
	} else {
		err := s.publisherRepo.PublishLongtermLoadOrderHTTPRequest(ctx, repoParam)
		if err != nil {
			return errors.Wrap(err)
		}

		err = s.publisherRepo.PublishShardingLoadOrderHTTPRequest(ctx, repoParam)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	return nil
}
