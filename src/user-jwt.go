package src

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
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

func JwtLogin(w http.ResponseWriter, r *http.Request) {
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

	userToken, _ := CreateToken(customUser)
	w.Header().Set("Authentication", userToken)
}

func CreateToken(user User) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID
	atClaims["username"] = user.Username
	atClaims["exp"] = time.Now().Add(time.Second * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func JwtWelcome(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authentication")
	claims, _ := extractClaims(tokenHeader)
	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims["username"])))
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return []byte(JwtSecret), nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
