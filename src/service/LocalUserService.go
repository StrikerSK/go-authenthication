package userServices

import (
	"context"
	"errors"
	userDomain "github.com/strikersk/user-auth/src/domain"
	userJwt "github.com/strikersk/user-auth/src/jwt"
	userPorts "github.com/strikersk/user-auth/src/ports"
)

type LocalUserService struct {
	repository userPorts.IUserRepository
	cache      userPorts.IUserCache
	jwt        userJwt.JWTConfiguration
}

func NewLocalUserRepository(repository userPorts.IUserRepository, cache userPorts.IUserCache, jwt userJwt.JWTConfiguration) LocalUserService {
	return LocalUserService{
		repository: repository,
		cache:      cache,
		jwt:        jwt,
	}
}

func (r *LocalUserService) CreateUser(ctx context.Context, user userDomain.User) error {
	return r.repository.CreateUser(user)
}

func (r *LocalUserService) ReadUser(ctx context.Context, username string) (userDomain.User, error) {
	if cachedUser, isPresent := r.cache.RetrieveCache(ctx, username); isPresent {
		//log.Println("cachedUser", cachedUser)
		return cachedUser, nil
	} else {
		user, err := r.repository.ReadUser(username)
		if err != nil {
			return userDomain.User{}, err
		}

		if err = r.cache.CreateCache(ctx, user); err != nil {
			return userDomain.User{}, err
		}
		//log.Println("persistedUser", user)
		return user, nil
	}
}

func (r *LocalUserService) LoginUser(ctx context.Context, credentials userDomain.UserCredentials) (string, error) {
	persistedUser, err := r.ReadUser(ctx, credentials.Username)
	if err != nil {
		return "", err
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if credentials.Password != persistedUser.Password {
		return "", errors.New("user unauthorized")
	}

	userToken, err := r.jwt.GenerateToken(persistedUser)
	if err != nil {
		return "", err
	}

	return userToken, nil
}
