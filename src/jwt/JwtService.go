package jwt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		log.Printf("Login() error: %s\n", err)
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if user.Password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userToken, err := Configuration.GenerateToken(customUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authentication", userToken)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	test, err := Configuration.ParseToken(r.Header.Get("Authentication"))
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Welcome() error: %s\n", err)
		return
	}

	// Finally, return the welcome message to the user
	w.Write([]byte(fmt.Sprintf("Welcome %s!", test.User.Username)))
}
