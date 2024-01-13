package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
	"log"
	"net/http"
)

type UserRegisterHandler struct {
	service ports.IUserService
}

func NewUserRegisterHandler(service ports.IUserService) UserRegisterHandler {
	return UserRegisterHandler{service: service}
}

func (h UserRegisterHandler) RegisterHandler(router *mux.Router) {
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", h.createUser).Methods("POST")
}

func (h UserRegisterHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var user domain.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println("User decoding error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.service.CreateUser(r.Context(), user); err != nil {
		log.Println("User register error:", err)
		constants.ResolveResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
