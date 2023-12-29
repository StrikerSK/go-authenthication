package handlers

import (
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/ports"
	"net/http"
	"time"
)

type CookiesHandler struct {
	AbstractHandler
}

func NewCookiesHandler(userService ports.IUserService, tokenService ports.IAuthorizationService, configuration config.Authorization) CookiesHandler {
	return CookiesHandler{
		AbstractHandler: NewAbstractHandler(userService, tokenService, configuration),
	}
}

func (h *CookiesHandler) writeHeader(token string, w http.ResponseWriter) {
	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    h.tokenName,
		Value:   token,
		Expires: time.Now().Add(h.expiration * time.Second),
	})

	w.WriteHeader(http.StatusOK)
}

func (h *CookiesHandler) readAuthorizationHeader(r *http.Request) (string, error) {
	c, err := r.Cookie(h.tokenName)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}
