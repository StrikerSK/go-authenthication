package handlers

import (
	"errors"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"net/http"
	"time"
)

type CookiesAuthorization struct {
	tokenName  string
	expiration time.Duration
}

func NewCookiesAuthorization(configuration config.Authorization) *CookiesAuthorization {
	return &CookiesAuthorization{
		tokenName:  configuration.AuthorizationHeader,
		expiration: time.Duration(configuration.TokenExpiration),
	}
}

func (h *CookiesAuthorization) AddAuthorization(token string, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    h.tokenName,
		Value:   token,
		Expires: time.Now().Add(h.expiration * time.Second),
	})
}

func (h *CookiesAuthorization) GetAuthorization(r *http.Request) (string, error) {
	c, err := r.Cookie(h.tokenName)

	if c.Value == "" {
		return "", errors.New(constants.MissingAuthorizationToken)
	}

	if err != nil {
		return "", err
	}

	return c.Value, nil
}
