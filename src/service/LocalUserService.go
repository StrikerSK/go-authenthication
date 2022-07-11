package userServices

import (
	uuid "github.com/satori/go.uuid"
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

func (r *LocalUserService) CreateUser(user domain.User) error {
	return r.repository.CreateUser(user)
}

func (r *LocalUserService) ReadUser(username string) (domain.User, error) {
	if cachedUser, _ := r.cache.RetrieveCache(username); cachedUser == nil {
		user, err := r.repository.ReadUser(username)
		if err != nil {
			return domain.User{}, err
		}

		// Create a new random session token
		sessionToken := uuid.NewV4().String()
		if err = r.cache.CreateCache(sessionToken, user); err != nil {
			return domain.User{}, err
		}

		return user, nil
	} else {
		return cachedUser.(domain.User), nil
	}
}
