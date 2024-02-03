package cache

import (
	"context"
	"errors"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"time"
)

type cachedUser struct {
	domain.UserDTO
	createdAt time.Time
}

type InMemoryCache struct {
	Users []cachedUser
}

func (r *InMemoryCache) CreateCache(ctx context.Context, user *domain.UserDTO) error {
	tmpUser := cachedUser{
		UserDTO:   *user,
		createdAt: time.Now(),
	}

	r.Users = append(r.Users, tmpUser)
	return nil
}

func (r *InMemoryCache) RetrieveCache(ctx context.Context, searchedUser *domain.UserDTO) (bool, error) {
	for _, user := range r.Users {
		if searchedUser.Username == user.Username {
			*searchedUser = user.UserDTO
			return true, nil
		}
	}

	return false, errors.New(constants.NotFoundConstant)
}
