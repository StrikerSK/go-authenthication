package userServices

import (
	"context"
	"github.com/strikersk/user-auth/src/domain"
	"github.com/strikersk/user-auth/src/ports"
)

type LocalUserService struct {
	repository ports.IUserRepository
	cache      ports.IUserCache
}

func NewLocalUserRepository(repository ports.IUserRepository, cache ports.IUserCache) LocalUserService {
	return LocalUserService{
		repository: repository,
		cache:      cache,
	}
}

func (r *LocalUserService) CreateUser(ctx context.Context, user domain.User) error {
	return r.repository.CreateUser(user)
}

func (r *LocalUserService) ReadUser(ctx context.Context, username string) (domain.User, error) {
	if cachedUser, isPresent := r.cache.RetrieveCache(ctx, username); !isPresent {
		user, err := r.repository.ReadUser(username)
		if err != nil {
			return domain.User{}, err
		}

		if err = r.cache.CreateCache(ctx, user); err != nil {
			return domain.User{}, err
		}
		//log.Println("persistedUser", user)
		return user, nil
	} else {
		//log.Println("cachedUser", cachedUser)
		return cachedUser, nil
	}
}
