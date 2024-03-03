package payment_types

import (
	"context"
	"skripsi-be/repository/payment_types"
	"skripsi-be/type/result"
)

type Service interface {
	FindLikeOneOfAllColumn(ctx context.Context, req string) ([]result.PaymentTypes, error)
}

type Config struct{}

type service struct {
	config           Config
	paymentTypesRepo payment_types.Repository
}

func New(
	config Config,
	paymentTypesRepo payment_types.Repository,
) Service {
	return service{
		config:           config,
		paymentTypesRepo: paymentTypesRepo,
	}
}
