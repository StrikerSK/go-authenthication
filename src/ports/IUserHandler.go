package ports

import (
	"github.com/gorilla/mux"
	"net/http"
)

type IUserHandler interface {
	RegisterHandler(router *mux.Router)
}

type IUserEndpointHandler interface {
	WriteAuthorizationHeader(token string, w http.ResponseWriter)
	ReadAuthorizationHeader(r *http.Request) (string, error)
}
