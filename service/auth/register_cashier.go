package auth

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/google/uuid"
)

func (s service) RegisterCashier(ctx context.Context, param params.ServiceRegisterCashier) error {
	err := param.Validate()
	if err != nil {
		return errors.New(err)
	}

	account, err := s.accountRepo.FindAccountByID(ctx, param.RequestBy)
	if errors.Is(err, errors.ErrRecordNotFound) {
		return errors.New(errors.ErrUnauthorized)
	} else if err != nil {
		return err
	}

	_, err = s.adminRepo.FindAdminAccountByID(ctx, account.ID)
	if errors.Is(err, errors.ErrRecordNotFound) {
		return errors.New(errors.ErrUnauthorized)
	} else if err != nil {
		return err
	}

	adminUsername, err := s.accountRepo.FindAccountByUsername(ctx, param.Username)
	if errors.Is(err, errors.ErrRecordNotFound) {
		err = nil
	} else if err != nil {
		return err
	}

	if adminUsername.ID != uuid.Nil {
		return errors.New(errors.ErrUsernameDuplicated)
	}

	cashierAccount := param.ToAccountModel()
	cashierAccount.AssignUUID()

	err = cashierAccount.HashPassword()
	if err != nil {
		return err
	}

	cashier := param.ToCashierModel(cashierAccount.ID)
	cashier.AssignUUID()

	err = s.accountRepo.InsertNewAccount(ctx, cashierAccount)
	if err != nil {
		return err
	}

	err = s.adminRepo.InsertNewCashier(ctx, cashier)
	if err != nil {
		return err
	}

	return nil
}
