package ports

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
)

type IUserService interface {
	LoginUser(context.Context, domain.UserCredentials) error
	CreateUser(context.Context, *domain.UserDTO) error
	ReadUser(context.Context, *domain.UserDTO) error
}
