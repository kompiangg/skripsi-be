package order

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
	"skripsi-be/type/params"

	"github.com/shopspring/decimal"
)

func (s service) TransformOrder(ctx context.Context, param []params.ServiceTransformOrder) error {
	for _, v := range param {
		err := v.Validate()
		if err != nil {
			return errors.Wrap(err)
		}
	}

	mapCurrency := make(map[string]decimal.Decimal)
	for _, v := range param {
		if _, ok := mapCurrency[v.Currency]; !ok {
			mapCurrency[v.Currency] = decimal.Zero
		}
	}

	for currency := range mapCurrency {
		usdRate, err := s.currencyRepo.FindUSDCurrencyRate(ctx, currency)
		if errors.Is(err, constant.ErrRedisNil) {
			usdRate = decimal.NewFromFloat(0)
		} else if err != nil {
			return errors.Wrap(err)
		}

		mapCurrency[currency] = usdRate
	}

	repoParam := make([]params.RepositoryPublishLoadOrderEvent, len(param))
	for i, v := range param {
		repoParam[i] = v.TransformOrder(mapCurrency[v.Currency])
	}

	err := s.publisherRepo.PublishLoadOrderEvent(ctx, repoParam)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}
