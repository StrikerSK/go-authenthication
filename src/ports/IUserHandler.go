package ports

import (
	"github.com/gorilla/mux"
	"net/http"
)

type IUserHandler interface {
	RegisterHandler(router *mux.Router)
}

type IUserEndpointHandler interface {
	AddAuthorization(token string, w http.ResponseWriter)
	GetAuthorization(r *http.Request) (string, error)
}
