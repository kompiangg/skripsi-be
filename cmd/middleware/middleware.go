package middleware

import (
	"skripsi-be/config"
	"skripsi-be/type/constant"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type middleware struct {
	config           Config
	cashierJWTConfig echojwt.Config
	adminJWTConfig   echojwt.Config
}

type Middleware interface {
	Admin() echo.MiddlewareFunc
	Cashier() echo.MiddlewareFunc
}

type Config struct {
	JWTConfig config.JWT
}

func New(config Config) middleware {
	return middleware{
		config: config,
		cashierJWTConfig: echojwt.Config{
			SigningKey: []byte(config.JWTConfig.Cashier.Secret),
			ContextKey: constant.CashierContextKey,
		},
		adminJWTConfig: echojwt.Config{
			SigningKey: []byte(config.JWTConfig.Admin.Secret),
			ContextKey: constant.AdminContextKey,
		},
	}
}
