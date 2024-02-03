package userServices

import (
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/domain"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type BcryptUserPasswordService struct {
	cost int
}

func NewBcryptUserPasswordService(configuration *config.EncryptionConfiguration) *BcryptUserPasswordService {
	cost := configuration.Cost
	if cost > bcrypt.MaxCost {
		log.Fatalf("Encryption value cannot be higher than %d", bcrypt.MaxCost)
	} else if cost < bcrypt.MinCost {
		log.Fatalf("Encryption value cannot be lower than %d", bcrypt.MinCost)
	}

	return &BcryptUserPasswordService{
		cost: configuration.Cost,
	}
}

func (ps *BcryptUserPasswordService) SetPassword(user *domain.UserCredentials) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), ps.cost)
	if err != nil {
		log.Println("password encryption error: ", err)
		return err
	}

	user.Password = string(encryptedPassword)
	return nil
}

func (ps *BcryptUserPasswordService) ValidatePassword(persistedCredentials *domain.UserCredentials, requestCredentials *domain.UserCredentials) error {
	return bcrypt.CompareHashAndPassword([]byte(persistedCredentials.Password), []byte(requestCredentials.Password))
}
