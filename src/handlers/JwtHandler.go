package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"log"
	"net/http"
)

type JwtHandler struct {
	userService ports.IUserService
	authService ports.IAuthorizationService
}

func NewJwtHandler(userService ports.IUserService, authService ports.IAuthorizationService) JwtHandler {
	return JwtHandler{
		userService: userService,
		authService: authService,
	}
}

func (h JwtHandler) EnrichRouter(router *mux.Router) {
	jwtRouter := router.PathPrefix("/user").Subrouter()
	jwtRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	jwtRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h JwtHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqUser domain.UserCredentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Login() error: %s\n", err)
		return
	}

	persistedUser, err := h.userService.ReadUser(r.Context(), reqUser.Username)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Login() error: %s\n", err)
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if reqUser.Password != persistedUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userToken, err := h.authService.GenerateToken(persistedUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", userToken)
}

func (h JwtHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username, err := h.authService.ParseToken(token)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Welcome() error: %s\n", err)
		return
	}

	// Finally, return the welcome message to the user
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s!", username)))
}
