package userServices

import (
	"github.com/strikersk/user-auth/src/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type DefaultUserPasswordService struct{}

func (ps *DefaultUserPasswordService) SetPassword(user *domain.UserCredentials) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password encryption error: ", err)
		return err
	}

	user.Password = string(encryptedPassword)
	return nil
}

func (ps *DefaultUserPasswordService) ValidatePassword(persistedCredentials, requestCredentials domain.UserCredentials) error {
	return bcrypt.CompareHashAndPassword([]byte(persistedCredentials.Password), []byte(requestCredentials.Password))
}
