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
	address := cacheUrlResolver(configuration)

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

func (receiver MemcacheCache) CreateCache(ctx context.Context, inputUser *domain.UserDTO) error {
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

func (receiver MemcacheCache) RetrieveCache(ctx context.Context, user *domain.UserDTO) (bool, error) {
	item, err := receiver.cacheClient.Get(cachePrefix + user.Username)
	if err != nil {
		if err.Error() == memcache.ErrCacheMiss.Error() {
			return false, nil
		} else {
			return false, err
		}
	}

	err = json.Unmarshal(item.Value, user)
	if err != nil {
		return false, err
	}

	return true, nil
}
