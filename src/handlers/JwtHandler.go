package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"log"
	"net/http"
)

type JwtHandler struct {
	tokenName   string
	userService ports.IUserService
	authService ports.IAuthorizationService
}

func NewJwtHandler(userService ports.IUserService, authService ports.IAuthorizationService, configuration config.Authorization) JwtHandler {
	return JwtHandler{
		tokenName:   configuration.AuthorizationHeader,
		userService: userService,
		authService: authService,
	}
}

func (h JwtHandler) RegisterHandler(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	userRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h JwtHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userCredentials domain.UserCredentials

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		log.Println("Body decoding error:", err)
		constants.ResolveResponse(w, err)
		return
	}

	persistedUser, err := h.userService.ReadUser(r.Context(), userCredentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		log.Println("User read error:", err)
		constants.ResolveResponse(w, err)
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if userCredentials.Password != persistedUser.Password {
		log.Println("User authorization did not pass")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userToken, err := h.authService.GenerateToken(persistedUser)
	if err != nil {
		log.Println("Token generation error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(h.tokenName, userToken)
	w.Header().Set("Content-Type", "application/json")
	user, err := json.Marshal(persistedUser)
	if err != nil {
		log.Println("User marshalling error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(user)
}

func (h JwtHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get(h.tokenName)
	if token == "" {
		log.Println("Cannot get token from header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	username, err := h.authService.ParseToken(token)
	if err != nil {
		log.Println("Token parsing error:", err)
		constants.ResolveResponse(w, err)
		return
	}

	// Finally, return the welcome message to the user
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s!", username)))
}
