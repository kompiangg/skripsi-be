package httpx

import (
	"skripsi-be/pkg/errors"
	"skripsi-be/type/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetJWTClaimsFromContext(c echo.Context, ctxKey string) (entity.CustomJWTClaims, error) {
	token, ok := c.Get(ctxKey).(*jwt.Token)
	if !ok {
		return entity.CustomJWTClaims{}, errors.New("cannot get jwt token from context")
	}

	claims, ok := token.Claims.(entity.CustomJWTClaims)
	if !ok {
		return entity.CustomJWTClaims{}, errors.New("cannot get jwt claims from context")
	}

	return claims, nil
}
