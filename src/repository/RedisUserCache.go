package userRepository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"time"
)

type RedisCache struct {
	redisClient *redis.Client
	expiration  time.Duration
}

func NewRedisCache(configuration config.CacheConfiguration) (connection RedisCache) {
	address := cacheUrsResolver(configuration)

	redisConnection := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := redisConnection.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Error during connecting Redis: %v\n", err)
	}

	if err := redisConnection.FlushDB(context.Background()).Err(); err != nil {
		log.Fatalf("Error during cleaning caches: %v\n", err)
	}

	// Assign the connection to the package level `cache` variable
	return RedisCache{
		redisClient: redisConnection,
		expiration:  time.Duration(configuration.Expiration),
	}
}

func (receiver RedisCache) CreateCache(ctx context.Context, inputUser domain.UserDTO) error {
	err := receiver.redisClient.Set(ctx, cachePrefix+inputUser.Username, inputUser, time.Second*receiver.expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (receiver RedisCache) RetrieveCache(ctx context.Context, username string) (domain.UserDTO, bool, error) {
	var user domain.UserDTO

	val, err := receiver.redisClient.Get(ctx, cachePrefix+username).Result()

	if err != nil {
		if err.Error() == redis.Nil.Error() {
			return domain.UserDTO{}, false, nil
		} else {
			return domain.UserDTO{}, false, err
		}
	}

	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return domain.UserDTO{}, false, err
	}

	return user, true, nil
}
