package database

import (
	"errors"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
)

type LocalUserRepository struct {
	Users []domain.UserDTO
}

func NewLocalUserRepository() *LocalUserRepository {
	return &LocalUserRepository{}
}

func (r *LocalUserRepository) CreateEntry(user *domain.UserDTO) error {
	for _, tmpUser := range r.Users {
		if user.Username == tmpUser.Username {
			return errors.New(constants.ConflictConstant)
		}
	}

	r.Users = append(r.Users, *user)
	return nil
}

func (r *LocalUserRepository) ReadEntry(searchedUser *domain.UserDTO) (bool, error) {
	for _, user := range r.Users {
		if searchedUser.Username == user.Username {
			*searchedUser = user
			return true, nil
		}
	}

	return false, errors.New(constants.NotFoundConstant)
}
