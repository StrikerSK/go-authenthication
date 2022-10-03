package ports

import (
	"context"
	userDomain "github.com/strikersk/user-auth/src/domain"
)

type IUserClient interface {
	RegisterUser(context.Context, userDomain.User) error
	LoginUser(context.Context, userDomain.User) (string, error)
}
