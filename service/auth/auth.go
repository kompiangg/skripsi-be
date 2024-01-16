package auth

import (
	"context"
	"skripsi-be/config"
	"skripsi-be/repository/account"
	"skripsi-be/repository/admin"
	"skripsi-be/type/params"
	"skripsi-be/type/result"
)

type Service interface {
	RegisterAdmin(ctx context.Context, param params.ServiceRegisterAdmin) error
	AdminLogin(ctx context.Context, param params.ServiceLoginAdmin) (result.ServiceAdminLoginResult, error)
}

type Config struct {
	Admin config.AdminJWTConfig
}

type service struct {
	config      Config
	accountRepo account.Repository
	adminRepo   admin.Repository
}

func New(
	config Config,
	accountRepo account.Repository,
	adminRepo admin.Repository,
) Service {
	return service{
		config:      config,
		accountRepo: accountRepo,
		adminRepo:   adminRepo,
	}
}
