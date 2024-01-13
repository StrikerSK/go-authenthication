package handlers

import (
	"github.com/strikersk/user-auth/config"
	"net/http"
	"time"
)

type CookiesHandler struct {
	tokenName  string
	expiration time.Duration
}

func NewCookiesHandler(configuration config.Authorization) CookiesHandler {
	return CookiesHandler{
		tokenName:  configuration.AuthorizationHeader,
		expiration: time.Duration(configuration.TokenExpiration),
	}
}

func (h *CookiesHandler) WriteAuthorizationHeader(token string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.tokenName,
		Value:   token,
		Expires: time.Now().Add(h.expiration * time.Second),
	})

	w.WriteHeader(http.StatusOK)
}

func (h *CookiesHandler) ReadAuthorizationHeader(r *http.Request) (string, error) {
	c, err := r.Cookie(h.tokenName)
	if err != nil {
		return "", err
	}

	return c.Value, nil
}
