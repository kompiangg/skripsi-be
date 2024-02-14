package auth

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/account"
	"skripsi-be/repository/admin"
	"skripsi-be/repository/cashier"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
)

type Service interface {
	Login(ctx context.Context, param params.ServiceLogin) (result.ServiceLoginResult, error)
	RegisterAdmin(ctx context.Context, param params.ServiceRegisterAdmin) error
	RegisterCashier(ctx context.Context, param params.ServiceRegisterCashier) error
}

type Config struct {
	JWT config.JWT
}

type service struct {
	config      Config
	accountRepo account.Repository
	adminRepo   admin.Repository
	cashierRepo cashier.Repository
}

func New(
	config Config,
	accountRepo account.Repository,
	adminRepo admin.Repository,
	cashierRepo cashier.Repository,
) Service {
	return service{
		config:      config,
		accountRepo: accountRepo,
		adminRepo:   adminRepo,
		cashierRepo: cashierRepo,
	}
}
