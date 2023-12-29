package ports

import "github.com/strikersk/user-auth/src/domain"

type IUserRepository interface {
	CreateEntry(domain.UserDTO) error
	ReadEntry(string) (domain.UserDTO, bool, error)
}
