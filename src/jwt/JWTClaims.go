package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/strikersk/user-auth/src/domain"
)

type CustomClaims struct {
	User domain.User
	jwt.StandardClaims
}
