package userRepository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/strikersk/user-auth/src/domain"
	"log"
	"time"
)

type UserCache struct {
	Cache     *redis.Client
	TokenName string
}

func NewCacheConfig() (connection UserCache) {
	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Assign the connection to the package level `cache` variable
	connection.Cache = conn
	connection.TokenName = "session_token"
	return
}

func (receiver UserCache) CreateCache(ctx context.Context, inputUser domain.UserDTO) error {
	err := receiver.Cache.Set(ctx, inputUser.Username, inputUser, time.Second*300).Err()
	if err != nil {
		panic(err)
	}

	return nil
}

func (receiver UserCache) RetrieveCache(ctx context.Context, userName string) (domain.UserDTO, bool) {
	var user domain.UserDTO

	val, err := receiver.Cache.Get(ctx, userName).Result()
	if err == redis.Nil {
		log.Println("user not found in cache")
		return domain.UserDTO{}, false
	}

	_ = json.Unmarshal([]byte(val), &user)
	//log.Println("Cache user", user)
	return user, true
}
