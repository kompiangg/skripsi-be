package httpx

import (
	"encoding/json"
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

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return entity.CustomJWTClaims{}, errors.New("cannot get jwt claims from context")
	}

	marshalled, err := json.Marshal(jwtClaims)
	if err != nil {
		return entity.CustomJWTClaims{}, errors.Wrap(err)
	}

	var claims entity.CustomJWTClaims
	err = json.Unmarshal(marshalled, &claims)
	if err != nil {
		return entity.CustomJWTClaims{}, errors.Wrap(err)
	}

	return claims, nil
}
