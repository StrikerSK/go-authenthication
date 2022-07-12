package domain

// Create a struct that models the structure of a user, both in the request body, and in the DB
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
