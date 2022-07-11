package userRepository

import (
	"errors"
	"github.com/strikersk/user-auth/src/domain"
)

type LocalUserRepository struct {
	Users []domain.User
}

func NewLocalUserRepository() LocalUserRepository {
	return LocalUserRepository{}
}

func (r *LocalUserRepository) CreateUser(user domain.User) error {
	for _, tmpUser := range r.Users {
		if user.Username == tmpUser.Username {
			return errors.New("user already created")
		}
	}

	r.Users = append(r.Users, user)
	return nil
}

func (r *LocalUserRepository) ReadUser(username string) (domain.User, error) {
	for _, user := range r.Users {
		if username == user.Username {
			return user, nil
		}
	}

	return domain.User{}, errors.New("user cannot be found")
}
