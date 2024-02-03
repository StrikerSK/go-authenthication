package cache

import (
	"context"
	"errors"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"time"
)

type cachedUser struct {
	domain.UserDTO
	createdAt time.Time
}

type InMemoryCache struct {
	users             []cachedUser
	conditionDuration time.Duration
}

func NewInMemoryCache() *InMemoryCache {
	conditionDuration, err := time.ParseDuration("3600s")
	if err != nil {
		log.Fatal("error parsing duration", err)
	}

	return &InMemoryCache{
		users:             make([]cachedUser, 0),
		conditionDuration: conditionDuration,
	}
}

func (r *InMemoryCache) CreateCache(ctx context.Context, user *domain.UserDTO) error {
	tmpUser := cachedUser{
		UserDTO:   *user,
		createdAt: time.Now(),
	}

	r.users = append(r.users, tmpUser)
	return nil
}

func (r *InMemoryCache) RetrieveCache(ctx context.Context, searchedUser *domain.UserDTO) (bool, error) {
	for _, user := range r.users {
		if searchedUser.Username == user.Username {
			if user.createdAt.Sub(time.Now()) > r.conditionDuration {
				return false, nil
			}

			*searchedUser = user.UserDTO
			return true, nil
		}
	}

	return false, errors.New(constants.NotFoundConstant)
}
