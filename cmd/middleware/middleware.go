package middleware

import (
	"skripsi-be/config"
	"skripsi-be/type/constant"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type middleware struct {
	config Config
	jwt    echojwt.Config
}

type Middleware interface {
	JWTRestricted() echo.MiddlewareFunc
}

type Config struct {
	JWT config.JWT
}

func New(config Config) middleware {
	return middleware{
		config: config,
		jwt: echojwt.Config{
			SigningKey: []byte(config.JWT.Secret),
			ContextKey: constant.AuthContextKey,
		},
	}
}
