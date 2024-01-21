package userServices

import (
	"github.com/strikersk/user-auth/src/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserPasswordService struct{}

func (ps *UserPasswordService) SetPassword(user *domain.UserCredentials) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password encryption error: ", err)
		return err
	}

	user.Password = string(encryptedPassword)
	return nil
}

func (ps *UserPasswordService) ValidatePassword(user domain.UserCredentials, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
