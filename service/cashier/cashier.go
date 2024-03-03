package cashier

import (
	"context"
	"skripsi-be/repository/cashier"
	"skripsi-be/type/result"

	"github.com/google/uuid"
)

type Service interface {
	FindCashierByID(ctx context.Context, id uuid.UUID) (res result.Cashier, err error)
}

type Config struct{}

type service struct {
	config      Config
	cashierRepo cashier.Repository
}

func New(
	config Config,
	cashierRepo cashier.Repository,
) Service {
	return service{
		config:      config,
		cashierRepo: cashierRepo,
	}
}
