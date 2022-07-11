package ports

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
)

type IUserCache interface {
	CreateCache(context.Context, domain.User) error
	RetrieveCache(context.Context, string) (domain.User, bool)
}
