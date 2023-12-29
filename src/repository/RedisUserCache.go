package userRepository

import (
	"context"
	"encoding/json"
	"fmt"
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

const cachePrefix = "user:"

func NewRedisCache(configuration config.CacheConfiguration) (connection RedisCache) {
	var address string

	if configuration.URL != "" {
		address = configuration.URL
	} else if configuration.Host != "" && configuration.Port != "" {
		address = fmt.Sprintf("%s:%s", configuration.Host, configuration.Port)
	} else {
		log.Fatal("cache address not provide")
	}

	redisConnection := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

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
		if err == redis.Nil {
			log.Println("user not found in cache")
			return domain.UserDTO{}, false, nil
		} else {
			return domain.UserDTO{}, false, err
		}
	}

	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return domain.UserDTO{}, false, err
	}

	//log.Println("Cache user", user)
	return user, true, nil
}
