package auth

import (
	"context"
	"skripsi-be/lib/tokenlib"
	"skripsi-be/pkg/errors"
	"skripsi-be/type/constant"
	"skripsi-be/type/entity"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s service) Login(ctx context.Context, param params.ServiceLogin) (res result.ServiceLoginResult, err error) {
	err = param.Validate()
	if err != nil {
		return res, errors.New(err)
	}

	account, err := s.accountRepo.FindAccountByUsername(ctx, param.Username)
	if errors.Is(err, errors.ErrRecordNotFound) {
		return res, errors.Wrap(errors.ErrAccountNotFound)
	} else if err != nil {
		return res, err
	}

	isMatch := account.ComparePassword(param.Password)
	if !isMatch {
		return res, errors.Wrap(errors.ErrIncorrectPassword)
	}

	var name string
	var id string

	cashier, err := s.cashierRepo.FindCashierAccountByID(ctx, account.ID)
	if errors.Is(err, errors.ErrRecordNotFound) {
		admin, err := s.adminRepo.FindAdminAccountByID(ctx, account.ID)
		if errors.Is(err, errors.ErrRecordNotFound) {
			return res, errors.Wrap(errors.ErrAccountNotFound)
		} else if err != nil {
			return res, errors.Wrap(err)
		}

		id = admin.ID.String()
		res.Role = constant.RoleEnum.Admin
	} else if err != nil {
		return res, errors.Wrap(err)
	} else {
		id = cashier.ID.String()
		res.Role = constant.RoleEnum.Cashier
		name = cashier.Name
	}

	now := time.Now()
	token, err := tokenlib.GenerateJWT(entity.CustomJWTClaims{
		Role: res.Role,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * time.Duration(s.config.JWT.ExpInDay))),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}, s.config.JWT.Secret)
	if err != nil {
		return res, err
	}

	res.Token = token

	return res, nil
}
