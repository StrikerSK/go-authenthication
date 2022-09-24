package ports

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
)

type IUserService interface {
	CreateUser(context.Context, domain.User) error
	ReadUser(context.Context, string) (domain.User, error)
	LoginUser(context.Context, domain.UserCredentials) (string, error)
}
