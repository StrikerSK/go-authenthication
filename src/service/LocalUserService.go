package userServices

import (
	"context"
	"errors"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
)

type LocalUserService struct {
	userRepository ports.IUserRepository
	userCache      ports.IUserCache
}

func NewUserService(userRepository ports.IUserRepository, userCache ports.IUserCache) LocalUserService {
	return LocalUserService{
		userRepository: userRepository,
		userCache:      userCache,
	}
}

func (r *LocalUserService) CreateUser(ctx context.Context, user domain.UserDTO) error {
	return r.userRepository.CreateUser(user)
}

func (r *LocalUserService) ReadUser(ctx context.Context, credentials domain.UserCredentials) (domain.UserDTO, error) {
	var user domain.UserDTO
	username := credentials.Username

	user, isPresent, err := r.userCache.RetrieveCache(ctx, username)
	if err != nil {
		return domain.UserDTO{}, err
	}

	if !isPresent {
		user, err = r.userRepository.ReadUser(username)
		if err != nil {
			return domain.UserDTO{}, err
		}

		if err = r.userCache.CreateCache(ctx, user); err != nil {
			return domain.UserDTO{}, err
		}
	}

	if user.Password != credentials.Password {
		return domain.UserDTO{}, errors.New("passwords are not matching")
	}

	return user, nil
}
