package middleware

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func (m middleware) JWTRestricted() echo.MiddlewareFunc {
	return echojwt.WithConfig(m.jwt)
}
