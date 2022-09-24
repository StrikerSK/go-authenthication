package userRepository

import (
	"github.com/strikersk/user-auth/src/domain"
	appErrors "github.com/strikersk/user-auth/src/errors"
)

type LocalUserRepository struct {
	Users []domain.User
}

func NewLocalUserRepository() LocalUserRepository {
	var localRepository LocalUserRepository
	localRepository.Users = []domain.User{
		{
			domain.UserCredentials{
				Username: "test",
				Password: "test",
			},
		},
	}
	return localRepository
}

func (r *LocalUserRepository) CreateUser(user domain.User) error {
	for _, tmpUser := range r.Users {
		if user.Username == tmpUser.Username {
			return appErrors.NewRepositoryError().FromMessage("user already created")
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

	return domain.User{}, appErrors.NewRepositoryError().FromMessage("user cannot be found")
}
