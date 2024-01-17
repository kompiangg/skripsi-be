package auth

import (
	"context"
	"skripsi-be/lib/tokenlib"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
)

func (s service) CashierLogin(ctx context.Context, param params.ServiceLoginCashier) (result.ServiceCashierLoginResult, error) {
	err := param.Validate()
	if err != nil {
		return result.ServiceCashierLoginResult{}, errors.New(err)
	}

	cashier, err := s.accountRepo.FindAccountByUsername(ctx, param.Username)
	if errors.Is(err, errors.ErrRecordNotFound) {
		return result.ServiceCashierLoginResult{}, errors.New(errors.ErrAccountNotFound)
	} else if err != nil {
		return result.ServiceCashierLoginResult{}, err
	}

	isMatch := cashier.ComparePassword(param.Password)
	if !isMatch {
		return result.ServiceCashierLoginResult{}, errors.New(errors.ErrIncorrectPassword)
	}

	token, err := tokenlib.GenerateJWT(cashier.ID.String(), s.config.Cashier.Secret, s.config.Cashier.ExpInDay)
	if err != nil {
		return result.ServiceCashierLoginResult{}, err
	}

	return result.ServiceCashierLoginResult{
		Token: token,
	}, nil
}
