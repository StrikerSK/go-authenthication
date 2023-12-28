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

func (r *LocalUserRepository) CreateUser(user domain.UserDTO) error {
	for _, tmpUser := range r.Users {
		if user.Username == tmpUser.Username {
			return errors.New("user already created")
		}
	}

	r.Users = append(r.Users, user)
	return nil
}

func (r *LocalUserRepository) ReadUser(username string) (domain.UserDTO, error) {
	for _, user := range r.Users {
		if username == user.Username {
			return user, nil
		}
	}

	return domain.UserDTO{}, errors.New("user cannot be found")
}
