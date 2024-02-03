package ports

import "github.com/strikersk/user-auth/src/domain"

type IUserRepository interface {
	CreateEntry(*domain.UserDTO) error
	ReadEntry(*domain.UserDTO) (bool, error)
}
