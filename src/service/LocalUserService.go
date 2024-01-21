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
	userPasswordService ports.IPasswordEncryptionService
}

func NewUserService(userRepository ports.IUserRepository, userCache ports.IUserCache) *LocalUserService {
	return &LocalUserService{
		userRepository:      userRepository,
		userCache:           userCache,
		userPasswordService: &UserPasswordService{},
	}
}

func (r *LocalUserService) CreateUser(ctx context.Context, user domain.UserDTO) error {
	persistedUser, exists, err := r.fetchUser(ctx, user.UserCredentials)
	if err != nil && err.Error() != constants.NotFoundConstant {
		return err
	}

	if exists && persistedUser.Username == user.Username {
		return errors.New(constants.ConflictConstant)
	}

	err = r.userPasswordService.SetPassword(&user.UserCredentials)
	if err != nil {
		return err
	}

	return r.userRepository.CreateEntry(user)
}

func (r *LocalUserService) ReadUser(ctx context.Context, credentials domain.UserCredentials) (domain.UserDTO, error) {
	user, exists, err := r.fetchUser(ctx, credentials)
	if err != nil {
		return domain.UserDTO{}, err
	}

	if !exists {
		return domain.UserDTO{}, errors.New(constants.NotFoundConstant)
	}

	return user, nil
}

func (r *LocalUserService) LoginUser(ctx context.Context, credentials domain.UserCredentials) error {
	user, err := r.ReadUser(ctx, credentials)
	if err != nil {
		return err
	}

	return r.userPasswordService.ValidatePassword(user.UserCredentials, credentials.Password)
}

func (r *LocalUserService) fetchUser(ctx context.Context, credentials domain.UserCredentials) (domain.UserDTO, bool, error) {
	var user domain.UserDTO
	username := credentials.Username

	user, isPresent, err := r.userCache.RetrieveCache(ctx, username)
	if err != nil {
		return domain.UserDTO{}, false, err
	}

	if !isPresent {
		user, isPresent, err = r.userRepository.ReadEntry(username)
		if err != nil {
			return domain.UserDTO{}, false, err
		}

		if !isPresent {
			return domain.UserDTO{}, false, nil
		}

		if err = r.userCache.CreateCache(ctx, user); err != nil {
			return user, false, err
		}
	}

	return user, true, nil
}
