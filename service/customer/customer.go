package customer

import (
	"context"
	"skripsi-be/repository/customer"
	"skripsi-be/type/result"
)

type Service interface {
	FindLikeOneOfAllColumn(ctx context.Context, req string) ([]result.Customer, error)
}

type Config struct {
}

type service struct {
	config       Config
	customerRepo customer.Repository
}

func New(
	config Config,
	customerRepo customer.Repository,
) Service {
	return service{
		config:       config,
		customerRepo: customerRepo,
	}
}
