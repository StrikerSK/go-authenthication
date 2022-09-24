package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/jwt"
	"github.com/strikersk/user-auth/src/ports"
	"log"
	"net/http"
)

type JwtHandler struct {
	Service ports.IUserService
	Config  jwt.JWTConfiguration
}

func NewJwtHandler(service ports.IUserService, config jwt.JWTConfiguration) JwtHandler {
	return JwtHandler{
		Service: service,
		Config:  config,
	}
}

func (h JwtHandler) EnrichRouter(router *mux.Router) {
	jwtRouter := router.PathPrefix("/jwt").Subrouter()
	jwtRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	jwtRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h JwtHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials domain.UserCredentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Login() error: %s\n", err)
		return
	}

	token, err := h.Service.LoginUser(context.Background(), credentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Login error: %s\n", err)
		return
	}

	w.Header().Set("Authorization", token)
}

func (h JwtHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	user, err := h.Config.ParseToken(token)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Welcome() error: %s\n", err)
		return
	}

	// Finally, return the welcome message to the user
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s!", user.User.Username)))
}
