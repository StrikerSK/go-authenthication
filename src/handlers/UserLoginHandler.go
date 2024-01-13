package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"log"
	"net/http"
)

type UserLoginHandler struct {
	userService  ports.IUserService
	tokenService ports.IAuthorizationService
	userEndpoint ports.IUserEndpointHandler
}

func NewUserLoginHandler(userService ports.IUserService, tokenService ports.IAuthorizationService, userEndpoint ports.IUserEndpointHandler) UserLoginHandler {
	return UserLoginHandler{
		userService:  userService,
		tokenService: tokenService,
		userEndpoint: userEndpoint,
	}
}

func (h UserLoginHandler) RegisterHandler(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	userRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h UserLoginHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	user, err := json.Marshal(persistedUser)
	if err != nil {
		log.Println("Error marshalling data", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.userEndpoint.AddAuthorization(token, w)
	w.Write(user)
}

func (h UserLoginHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	token, err := h.userEndpoint.GetAuthorization(r)
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
