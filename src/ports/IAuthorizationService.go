package ports

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
)

type IAuthorizationService interface {
	ParseToken(context.Context, string) (string, error)
	GenerateToken(context.Context, domain.UserDTO) (string, error)
}
