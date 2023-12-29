package userServices

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"time"
)

type JWTService struct {
	secret     string
	expiration time.Duration
}

func NewJWTService(authorization config.Authorization) JWTService {
	return JWTService{
		secret:     authorization.JWT.TokenEncoding,
		expiration: time.Duration(authorization.TokenExpiration),
	}
}

func (receiver JWTService) ParseToken(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&domain.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(receiver.secret), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*domain.UserClaims)
	if !ok {
		err = errors.New("could not parse claims")
		return "", err
	}

	if claims.ExpiresAt.Before(time.Now().Local()) {
		err = errors.New(constants.ExpiredTokenConstant)
		return "", err
	}

	return claims.User.Username, nil
}

func (receiver JWTService) GenerateToken(user domain.UserDTO) (signedToken string, err error) {
	currentTime := time.Now().Local()
	claims := &domain.UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Second * receiver.expiration)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(receiver.secret))
	if err != nil {
		return
	}

	return
}
