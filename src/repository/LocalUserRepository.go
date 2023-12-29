package userRepository

import (
	"errors"
	"github.com/strikersk/user-auth/src/domain"
)

type LocalUserRepository struct {
	Users []domain.UserDTO
}

func NewLocalUserRepository() LocalUserRepository {
	return LocalUserRepository{}
}

func (r *LocalUserRepository) CreateEntry(user domain.UserDTO) error {
	for _, tmpUser := range r.Users {
		if user.Username == tmpUser.Username {
			return errors.New("user already created")
		}
	}

	r.Users = append(r.Users, user)
	return nil
}

func (r *LocalUserRepository) ReadEntry(username string) (domain.UserDTO, bool, error) {
	for _, user := range r.Users {
		if username == user.Username {
			return user, true, nil
		}
	}

	return domain.UserDTO{}, false, nil
}
