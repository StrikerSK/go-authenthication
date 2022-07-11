package ports

import "github.com/strikersk/user-auth/src/domain"

type IUserRepository interface {
	CreateUser(domain.User) error
	ReadUser(string) (domain.User, error)
}
