package domain

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Create a struct that models the structure of a user, both in the request body, and in the DB
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func NewCredentials(username, password string) UserCredentials {
	return UserCredentials{
		Username: username,
		Password: password,
	}
}

func (c *UserCredentials) ClearPassword() {
	c.Password = ""
}

func (c *UserCredentials) EncryptPassword() {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("password decryption error: %s\n", err)
	}

	c.Password = string(encryptedPassword)
}

func (c *UserCredentials) ValidatePassword(inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(inputPassword))
}
