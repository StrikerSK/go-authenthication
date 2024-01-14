package constants

import "net/http"

func ResolveResponse(w http.ResponseWriter, err error) {
	switch err.Error() {
	case ExpiredAuthorizationToken, MissingAuthorizationToken, InvalidAuthorizationToken:
		w.WriteHeader(http.StatusUnauthorized)
	case NotFoundConstant:
		w.WriteHeader(http.StatusNotFound)
	case ConflictConstant:
		w.WriteHeader(http.StatusConflict)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}
