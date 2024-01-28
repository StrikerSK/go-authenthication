package cache

import (
	"context"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"time"
)

type MemcacheCache struct {
	cacheClient *memcache.Client
	expiration  time.Duration
}

func NewMemcacheCache(configuration config.CacheConfiguration) (connection MemcacheCache) {
	address := cacheUrsResolver(configuration)

	cacheConnection := memcache.New(address)

	if err := cacheConnection.Ping(); err != nil {
		log.Fatalf("Error during connecting Redis: %v\n", err)
	}

	if err := cacheConnection.FlushAll(); err != nil {
		log.Fatalf("Error during cleaning caches: %v\n", err)
	}

	// Assign the connection to the package level `cache` variable
	return MemcacheCache{
		cacheClient: cacheConnection,
		expiration:  time.Duration(configuration.Expiration),
	}
}

func (receiver MemcacheCache) CreateCache(ctx context.Context, inputUser domain.UserDTO) error {
	userData, err := json.Marshal(inputUser)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        cachePrefix + inputUser.Username,
		Value:      userData,
		Expiration: int32(receiver.expiration),
	}

	err = receiver.cacheClient.Set(item)
	if err != nil {
		return err
	}

	return nil
}

func (receiver MemcacheCache) RetrieveCache(ctx context.Context, username string) (domain.UserDTO, bool, error) {
	var user domain.UserDTO

	item, err := receiver.cacheClient.Get(username)
	if err != nil {
		if err.Error() == memcache.ErrCacheMiss.Error() {
			return domain.UserDTO{}, false, nil
		} else {
			return domain.UserDTO{}, false, err
		}
	}

	err = user.UnmarshalBinary(item.Value)
	if err != nil {
		return domain.UserDTO{}, false, err
	}

	return user, true, nil
}
