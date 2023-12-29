package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"log"
	"net/http"
	"time"
)

type AbstractHandler struct {
	tokenName    string
	expiration   time.Duration
	userService  ports.IUserService
	tokenService ports.IAuthorizationService
}

func NewAbstractHandler(userService ports.IUserService, tokenService ports.IAuthorizationService, configuration config.Authorization) AbstractHandler {
	return AbstractHandler{
		tokenName:    configuration.AuthorizationHeader,
		expiration:   time.Duration(configuration.TokenExpiration),
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h AbstractHandler) RegisterHandler(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	userRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h AbstractHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userCredentials domain.UserCredentials

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	persistedUser, err := h.userService.ReadUser(r.Context(), userCredentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If it is the same as the password we received, then we can move ahead if NOT, then we return an "Unauthorized" status
	if persistedUser.Password != userCredentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	token, err := h.tokenService.GenerateToken(domain.UserDTO{UserCredentials: userCredentials})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.writeHeader(token, w)

	w.Header().Set("Content-Type", "application/json")
	user, err := json.Marshal(persistedUser)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(user)
}

func (h *AbstractHandler) writeHeader(token string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h AbstractHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	token, err := h.readAuthorizationHeader(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := token

	// We then get the name of the user from our cache, where we set the session token
	// Create a new random session token
	username, err := h.tokenService.ParseToken(sessionToken)
	if err != nil {
		switch err.Error() {
		case constants.ExpiredTokenConstant:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// Finally, return the welcome message to the user
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s!", username)))
}

func (h *AbstractHandler) readAuthorizationHeader(r *http.Request) (string, error) {
	return "", errors.New("not implemented")
}
