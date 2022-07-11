package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"os"
	"time"
)

type JWTConfiguration struct {
	JwtSecret string
}

func NewConfigStruct() JWTConfiguration {
	return JWTConfiguration{
		JwtSecret: retrieveFromEnvironment(),
	}
}

func retrieveFromEnvironment() (secret string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret = os.Getenv("TOKEN_SECRET")
	return
}

func (receiver JWTConfiguration) ParseToken(signedToken string) (claims *UserClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(receiver.JwtSecret), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		err = errors.New("could not parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return
}

func (receiver JWTConfiguration) GenerateToken(user domain.User) (signedToken string, err error) {
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
