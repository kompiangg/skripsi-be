package auth

import (
	"context"
	"skripsi-be/lib/tokenlib"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
)

func (s service) AdminLogin(ctx context.Context, param params.ServiceLoginAdmin) (result.ServiceAdminLoginResult, error) {
	err := param.Validate()
	if err != nil {
		return result.ServiceAdminLoginResult{}, errors.New(err)
	}

	admin, err := s.accountRepo.FindAccountByUsername(ctx, param.Username)
	if errors.Is(err, errors.ErrRecordNotFound) {
		return result.ServiceAdminLoginResult{}, errors.New(errors.ErrAccountNotFound)
	} else if err != nil {
		return result.ServiceAdminLoginResult{}, err
	}

	isMatch := admin.ComparePassword(param.Password)
	if !isMatch {
		return result.ServiceAdminLoginResult{}, errors.New(errors.ErrIncorrectPassword)
	}

	token, err := tokenlib.GenerateJWT(admin.ID.String(), s.config.Admin.Secret, s.config.Admin.ExpInDay)
	if err != nil {
		return result.ServiceAdminLoginResult{}, err
	}

	return result.ServiceAdminLoginResult{
		Token: token,
	}, nil
}
