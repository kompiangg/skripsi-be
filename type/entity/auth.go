package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (c CustomJWTClaims) GetRole() (string, error) {
	return c.Role, nil
}
