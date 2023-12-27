package userServices

import (
	"context"
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

func (r *LocalUserService) ReadUser(ctx context.Context, username string) (domain.UserDTO, error) {
	if cachedUser, isPresent := r.userCache.RetrieveCache(ctx, username); !isPresent {
		user, err := r.userRepository.ReadUser(username)
		if err != nil {
			return domain.UserDTO{}, err
		}

		if err = r.userCache.CreateCache(ctx, user); err != nil {
			return domain.UserDTO{}, err
		}
		//log.Println("persistedUser", user)
		return user, nil
	} else {
		//log.Println("cachedUser", cachedUser)
		return cachedUser, nil
	}
}
