package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/strikersk/user-auth/src/domain"
)

type UserClaims struct {
	User domain.UserDTO
	jwt.StandardClaims
}
