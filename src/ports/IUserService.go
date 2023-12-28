package ports

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
)

type IUserService interface {
	CreateUser(context.Context, domain.UserDTO) error
	ReadUser(context.Context, domain.UserCredentials) (domain.UserDTO, error)
}
