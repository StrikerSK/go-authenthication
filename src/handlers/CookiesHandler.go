package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"net/http"
	"time"
)

type CookiesHandler struct {
	TokenName string
	Service   ports.IUserService
}

func NewCookiesHandler(tokenName string, service ports.IUserService) CookiesHandler {
	return CookiesHandler{
		TokenName: tokenName,
		Service:   service,
	}
}

func (h CookiesHandler) EnrichRouter(router *mux.Router) {
	jwtRouter := router.PathPrefix("/cookies").Subrouter()
	jwtRouter.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	jwtRouter.HandleFunc("/welcome", h.Welcome).Methods(http.MethodGet)
}

func (h CookiesHandler) Login(w http.ResponseWriter, r *http.Request) {
	var reqUser domain.Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	persistedUser, err := h.Service.ReadUser(r.Context(), reqUser.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if persistedUser.Password != reqUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := base64.URLEncoding.EncodeToString([]byte(reqUser.Username))

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    h.TokenName,
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})
}

func (h CookiesHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie(h.TokenName)
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
	response, err := base64.URLEncoding.DecodeString(sessionToken)

	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user
	_, _ = w.Write([]byte(fmt.Sprintf("Welcome %s!", response)))
}
