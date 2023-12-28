package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"net/http"
)

type UserHandler struct {
	service ports.IUserService
}

func NewUserHandler(service ports.IUserService) UserHandler {
	return UserHandler{service: service}
}

func (h UserHandler) RegisterHandler(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", h.createUser).Methods("POST")
}

func (h UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user domain.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.Header()
		return
	}

	_ = h.service.CreateUser(r.Context(), user)

	w.WriteHeader(http.StatusCreated)
	return
}
