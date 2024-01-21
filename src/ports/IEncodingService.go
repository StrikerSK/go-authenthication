package ports

import (
	"github.com/strikersk/user-auth/src/domain"
)

type IEncodingService interface {
	ParseToken(string) (string, error)
	GenerateToken(domain.UserDTO) (string, error)
}

type IPasswordEncryptionService interface {
	SetPassword(*domain.UserCredentials) error
	ValidatePassword(domain.UserCredentials, domain.UserCredentials) error
}
