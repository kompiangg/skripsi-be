package cashier

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/result"

	"github.com/google/uuid"
)

func (s service) FindCashierByID(ctx context.Context, id uuid.UUID) (res result.Cashier, err error) {
	cashier, err := s.cashierRepo.FindByID(ctx, id)
	if err != nil {
		return res, errors.Wrap(err)
	}

	res = result.Cashier{
		ID:        cashier.ID,
		Name:      cashier.Name,
		StoreID:   cashier.StoreID,
		CreatedAt: cashier.CreatedAt,
	}

	return res, nil
}
