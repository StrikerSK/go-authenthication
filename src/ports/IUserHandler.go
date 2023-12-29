package ports

import "github.com/gorilla/mux"

type IUserHandler interface {
	RegisterHandler(router *mux.Router)
}
