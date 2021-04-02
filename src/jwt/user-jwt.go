package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/strikersk/user-auth/src"
	"net/http"
	"time"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//A sample use
var customUser = User{
	ID:       1,
	Username: "Tester User",
	Password: "password",
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if user.Password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userToken, _ := GenerateToken(customUser)
	w.Header().Set("Authentication", userToken)
}

type CustomClaims struct {
	User User
	jwt.StandardClaims
}

func GenerateToken(user User) (signedToken string, err error) {
	claims := &CustomClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * 15).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(src.JwtSecret))
	if err != nil {
		return
	}

	return
}

func ParseToken(signedToken string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(src.JwtSecret), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authentication")
	test, _ := ParseToken(tokenHeader)
	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", test.User.Username)))
}
