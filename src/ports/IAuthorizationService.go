package ports

import (
	"github.com/strikersk/user-auth/src/domain"
)

type IAuthorizationService interface {
	ParseToken(string) (string, error)
	GenerateToken(domain.UserDTO) (string, error)
}
