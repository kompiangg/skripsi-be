package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func (m middleware) Admin() echo.MiddlewareFunc {
	return echojwt.WithConfig(m.adminJWTConfig)
}

func (m middleware) Cashier() echo.MiddlewareFunc {
	return echojwt.WithConfig(m.cashierJWTConfig)
}
