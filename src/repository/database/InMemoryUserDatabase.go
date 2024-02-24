package database

import (
	"errors"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
)

type InMemoryUserDatabase struct {
	Users []domain.UserDTO
}

func NewInMemoryUserDatabase() *InMemoryUserDatabase {
	return &InMemoryUserDatabase{}
}

func (r *InMemoryUserDatabase) CreateEntry(user *domain.UserDTO) error {
	for _, tmpUser := range r.Users {
		if user.Username == tmpUser.Username {
			return errors.New(constants.ConflictConstant)
		}
	}

	r.Users = append(r.Users, *user)
	return nil
}

func (r *InMemoryUserDatabase) ReadEntry(searchedUser *domain.UserDTO) (bool, error) {
	for _, user := range r.Users {
		if searchedUser.Username == user.Username {
			*searchedUser = user
			return true, nil
		}
	}

	return false, errors.New(constants.NotFoundConstant)
}
