package userServices

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"time"
)

type JWTEncodingService struct {
	secret     string
	expiration time.Duration
	sessionID  string
}

func NewJWTEncodingService(authorization config.Authorization) JWTEncodingService {
	return JWTEncodingService{
		secret:     authorization.JWT.TokenEncoding,
		expiration: time.Duration(authorization.TokenExpiration),
		sessionID:  uuid.NewString(),
	}
}

func (receiver JWTEncodingService) ParseToken(signedToken string) (string, error) {
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

	if !token.Valid {
		return "", errors.New(constants.InvalidAuthorizationToken)
	}

	claims, ok := token.Claims.(*domain.UserClaims)
	if !ok {
		err = errors.New("could not parse claims")
		return "", err
	}

	if claims.ID != receiver.sessionID {
		log.Print("Token from different session used")
		return "", errors.New(constants.InvalidAuthorizationToken)
	}

	if claims.ExpiresAt.Before(time.Now().Local()) {
		log.Print("Expired authorization token")
		err = errors.New(constants.ExpiredAuthorizationToken)
		return "", err
	}

	return claims.User.Username, nil
}

func (receiver JWTEncodingService) GenerateToken(user domain.UserDTO) (signedToken string, err error) {
	currentTime := time.Now().Local()
	claims := &domain.UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Second * receiver.expiration)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ID:        receiver.sessionID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(receiver.secret))
	if err != nil {
		return
	}

	return
}
