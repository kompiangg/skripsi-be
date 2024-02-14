package tokenlib

import (
	"skripsi-be/type/entity"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(claims entity.CustomJWTClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
