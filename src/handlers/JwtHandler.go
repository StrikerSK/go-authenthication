package handlers

import (
	"errors"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"log"
	"net/http"
)

type JwtHandler struct {
	tokenName string
}

func NewJwtHandler(configuration config.Authorization) *JwtHandler {
	return &JwtHandler{
		tokenName: configuration.AuthorizationHeader,
	}
}

func (h *JwtHandler) AddAuthorization(token string, w http.ResponseWriter) {
	w.Header().Set(h.tokenName, token)
}

func (h *JwtHandler) GetAuthorization(r *http.Request) (string, error) {
	token := r.Header.Get(h.tokenName)

	if token == "" {
		log.Println("Cannot get token from header")
		return "", errors.New(constants.MissingJwtToken)
	}

	return token, nil
}
