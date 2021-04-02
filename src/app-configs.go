package src

import (
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var cache redis.Conn
var JwtSecret string

func InitCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}

	// Assign the connection to the package level `cache` variable
	cache = conn
}

func InitEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JwtSecret = os.Getenv("TOKEN_SECRET")
}
