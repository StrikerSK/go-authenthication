package cache

import (
	"context"
	"fmt"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"time"
)

type cachedUser struct {
	domain.UserDTO
	createdAt time.Time
}

type InMemoryCache struct {
	users             map[string]cachedUser
	conditionDuration time.Duration
}

func NewInMemoryCache(configuration config.CacheConfiguration) *InMemoryCache {
	stringDuration := fmt.Sprintf("%ds", configuration.Expiration)
	conditionDuration, err := time.ParseDuration(stringDuration)
	if err != nil {
		log.Fatal("error parsing duration", err)
	}

	return &InMemoryCache{
		users:             make(map[string]cachedUser),
		conditionDuration: conditionDuration,
	}
}

func (r *InMemoryCache) CreateCache(ctx context.Context, user *domain.UserDTO) error {
	tmpUser := cachedUser{
		UserDTO:   *user,
		createdAt: time.Now(),
	}

	r.users[user.Username] = tmpUser
	return nil
}

func (r *InMemoryCache) RetrieveCache(ctx context.Context, searchedUser *domain.UserDTO) (bool, error) {
	requestTime := time.Now()
	foundUser, ok := r.users[searchedUser.Username]
	if !ok {
		return false, nil
	}

	if requestTime.Sub(foundUser.createdAt) > r.conditionDuration {
		return false, nil
	}

	*searchedUser = foundUser.UserDTO
	return true, nil
}
