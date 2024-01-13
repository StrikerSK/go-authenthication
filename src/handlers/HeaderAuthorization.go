package handlers

import (
	"errors"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"log"
	"net/http"
)

type HeaderAuthorization struct {
	tokenName string
}

func NewHeaderAuthorization(configuration config.Authorization) *HeaderAuthorization {
	return &HeaderAuthorization{
		tokenName: configuration.AuthorizationHeader,
	}
}

func (h *HeaderAuthorization) AddAuthorization(token string, w http.ResponseWriter) {
	w.Header().Set(h.tokenName, token)
}

func (h *HeaderAuthorization) GetAuthorization(r *http.Request) (string, error) {
	token := r.Header.Get(h.tokenName)

	if token == "" {
		log.Println("Cannot get token from header")
		return "", errors.New(constants.MissingJwtToken)
	}

	return token, nil
}
