package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"net/http"
	"time"
)

type CookiesHandler struct {
	tokenName    string
	expiration   time.Duration
	userService  ports.IUserService
	tokenService ports.IAuthorizationService
}

func NewCookiesHandler(userService ports.IUserService, tokenService ports.IAuthorizationService, configuration config.Authorization) CookiesHandler {
	return CookiesHandler{
		tokenName:    configuration.AuthorizationHeader,
		expiration:   time.Duration(configuration.TokenExpiration),
		userService:  userService,
		tokenService: tokenService,
	}
}

func (h CookiesHandler) RegisterHandler(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	userRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h CookiesHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	sessionToken, err := h.tokenService.GenerateToken(domain.UserDTO{UserCredentials: userCredentials})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    h.tokenName,
		Value:   sessionToken,
		Expires: time.Now().Add(h.expiration * time.Second),
	})

	w.Header().Set("Content-Type", "application/json")
	user, err := json.Marshal(persistedUser)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write(user)
}

func (h CookiesHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie(h.tokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

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
