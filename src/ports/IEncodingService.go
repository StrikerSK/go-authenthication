package ports

import (
	"github.com/strikersk/user-auth/src/domain"
)

type IEncodingService interface {
	ParseToken(*domain.UserCredentials) error
	GenerateToken(domain.UserDTO) (string, error)
}

type IPasswordEncryptionService interface {
	SetPassword(user *domain.UserCredentials) error
	ValidatePassword(*domain.UserCredentials, string) error
}
