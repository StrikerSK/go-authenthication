package ports

import "github.com/strikersk/user-auth/src/domain"

type IUserRepository interface {
	CreateUser(domain.UserDTO) error
	ReadUser(string) (domain.UserDTO, error)
}
