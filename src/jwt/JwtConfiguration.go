package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var Configuration = InitJwtConfig()

type ConfigStruct struct {
	JwtSecret string
}

func InitJwtConfig() ConfigStruct {
	var jwtConfig ConfigStruct

	jwtConfig.JwtSecret = initEnvFile()

	return jwtConfig
}

func initEnvFile() (secret string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret = os.Getenv("TOKEN_SECRET")
	return
}

type CustomClaims struct {
	User User
	jwt.StandardClaims
}

func (receiver ConfigStruct) ParseToken(signedToken string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(receiver.JwtSecret), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*CustomClaims)
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

func (receiver ConfigStruct) GenerateToken(user User) (signedToken string, err error) {
	claims := &CustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(receiver.JwtSecret))
	if err != nil {
		return
	}

	return
}
