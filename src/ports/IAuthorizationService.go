package ports

import (
	"github.com/strikersk/user-auth/src/domain"
)

type IAuthorizationService interface {
	ParseToken(signedToken string) (string, error)
	GenerateToken(user domain.UserDTO) (signedToken string, err error)
}
