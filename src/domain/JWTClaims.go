package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Username string
	jwt.RegisteredClaims
}
