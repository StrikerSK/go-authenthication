package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/strikersk/user-auth/src/domain"
	"time"
)

type JWTService struct {
	JwtSecret string
}

func NewConfigStruct() JWTService {
	return JWTService{
		JwtSecret: retrieveFromEnvironment(),
	}
}

func (receiver JWTService) ParseToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(receiver.JwtSecret), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		err = errors.New("could not parse claims")
		return "", err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return "", err
	}

	return claims.User.Username, nil
}

func (receiver JWTService) GenerateToken(user domain.UserDTO) (signedToken string, err error) {
	claims := &UserClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * 600).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(receiver.JwtSecret))
	if err != nil {
		return
	}

	return
}
