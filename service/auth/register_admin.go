package auth

import (
	"context"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"

	"github.com/google/uuid"
)

func (s service) RegisterAdmin(ctx context.Context, param params.ServiceRegisterAdmin) error {
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

	adminAccount := param.ToAccountModel()
	adminAccount.AssignUUID()

	err = adminAccount.HashPassword()
	if err != nil {
		return err
	}

	admin := param.ToAdminModel(adminAccount.ID)
	admin.AssignUUID()

	err = s.accountRepo.InsertNewAccount(ctx, adminAccount)
	if err != nil {
		return err
	}

	err = s.adminRepo.InsertNewAdmin(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}
