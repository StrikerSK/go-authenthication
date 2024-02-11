package constants

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ResolveResponse(w http.ResponseWriter, err interface{}) {
	switch resolveString(err) {
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

func resolveString(value interface{}) string {
	switch value.(type) {
	case error:
		return value.(error).Error()
	case string:
		return value.(string)
	default:
		return "value could not be resolved"
	}
}
