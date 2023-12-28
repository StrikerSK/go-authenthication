package ports

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
)

type IUserCache interface {
	CreateCache(context.Context, domain.UserDTO) error
	RetrieveCache(context.Context, string) (domain.UserDTO, bool, error)
}
