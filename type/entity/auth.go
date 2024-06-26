package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	Role string `json:"role"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}
