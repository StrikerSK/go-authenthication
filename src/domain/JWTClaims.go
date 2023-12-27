package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	User UserDTO
	jwt.RegisteredClaims
}
