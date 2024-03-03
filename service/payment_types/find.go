package payment_types

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/result"
)

func (s service) FindLikeOneOfAllColumn(ctx context.Context, req string) ([]result.PaymentTypes, error) {
	paymentTypes, err := s.paymentTypesRepo.FindByLikeOneOfAllColumn(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	res := make([]result.PaymentTypes, len(paymentTypes))
	for i, paymentType := range paymentTypes {
		res[i].FromModel(paymentType)
	}

	return res, nil
}
