package userServices

import (
	"context"
	"errors"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
)

type LocalUserService struct {
	userRepository      ports.IUserRepository
	userCache           ports.IUserCache
	userPasswordService ports.IUserPasswordService
}

func NewUserService(userRepository ports.IUserRepository, userCache ports.IUserCache, userPasswordService ports.IUserPasswordService) *LocalUserService {
	return &LocalUserService{
		userRepository:      userRepository,
		userCache:           userCache,
		userPasswordService: userPasswordService,
	}
}

func (r *LocalUserService) CreateUser(ctx context.Context, user *domain.UserDTO) error {
	exists, err := r.fetchUser(ctx, user)
	if err != nil && err.Error() != constants.NotFoundConstant {
		return err
	}

	if exists {
		return errors.New(constants.ConflictConstant)
	}

	err = r.userPasswordService.SetPassword(&user.UserCredentials)
	if err != nil {
		return err
	}

	return r.userRepository.CreateEntry(user)
}

func (r *LocalUserService) ReadUser(ctx context.Context, searchedUser *domain.UserDTO) error {
	exists, err := r.fetchUser(ctx, searchedUser)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New(constants.NotFoundConstant)
	}

	return nil
}

func (r *LocalUserService) LoginUser(ctx context.Context, credentials domain.UserCredentials) error {
	persistedUser := domain.UserDTO{
		UserCredentials: credentials,
	}

	err := r.ReadUser(ctx, &persistedUser)
	if err != nil {
		return err
	}

	return r.userPasswordService.ValidatePassword(persistedUser.UserCredentials, credentials)
}

func (r *LocalUserService) fetchUser(ctx context.Context, searchedUser *domain.UserDTO) (bool, error) {
	user, isPresent, err := r.userCache.RetrieveCache(ctx, searchedUser.Username)
	if err != nil {
		return false, err
	}

	if !isPresent {
		isPresent, err = r.userRepository.ReadEntry(searchedUser)
		if err != nil {
			return false, err
		}

		if !isPresent {
			return false, nil
		}

		if err = r.userCache.CreateCache(ctx, user); err != nil {
			return false, err
		}
	}

	return true, nil
}
