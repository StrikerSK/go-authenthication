package ports

import "github.com/strikersk/user-auth/src/domain"

type IUserService interface {
	CreateUser(domain.User) error
	ReadUser(string) (domain.User, error)
}
