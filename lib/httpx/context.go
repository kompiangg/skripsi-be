package httpx

import (
	"skripsi-be/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetJWTClaimsFromContext(c echo.Context, adminCashierContextKey string) (jwt.MapClaims, error) {
	token, ok := c.Get(adminCashierContextKey).(*jwt.Token)
	if !ok {
		return nil, errors.New("cannot get jwt token from context")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot get jwt claims from context")
	}

	return claims, nil
}
