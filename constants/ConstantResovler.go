package constants

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ResolveResponse(w http.ResponseWriter, err error) {
	switch err.Error() {
	case ExpiredAuthorizationToken, MissingAuthorizationToken, InvalidAuthorizationToken, bcrypt.ErrMismatchedHashAndPassword.Error():
		w.WriteHeader(http.StatusUnauthorized)
	case NotFoundConstant:
		w.WriteHeader(http.StatusNotFound)
	case ConflictConstant:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
