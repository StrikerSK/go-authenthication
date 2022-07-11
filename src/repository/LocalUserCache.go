package userRepository

import (
	"github.com/gomodule/redigo/redis"
)

type UserCache struct {
	Cache     redis.Conn
	TokenName string
}

func NewCacheConfig() (connection UserCache) {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}

	// Assign the connection to the package level `cache` variable
	connection.Cache = conn
	connection.TokenName = "session_token"
	return
}

func (receiver UserCache) CreateCache(sessionToken string, input interface{}) error {
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	_, err := receiver.Cache.Do("SETEX", sessionToken, "120", input)
	if err != nil {
		return err
	}

	return nil
}

func (receiver UserCache) RetrieveCache(sessionToken string) (interface{}, error) {
	// We then get the name of the user from our cache, where we set the session token
	return receiver.Cache.Do("GET", sessionToken)
}
