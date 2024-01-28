package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

// TODO Consider using only username

type UserClaims struct {
	User UserDTO
	jwt.RegisteredClaims
}
