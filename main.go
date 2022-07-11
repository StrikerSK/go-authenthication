package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/strikersk/user-auth/src/handlers"
	"github.com/strikersk/user-auth/src/jwt"
	"github.com/strikersk/user-auth/src/redis"
	userRepository "github.com/strikersk/user-auth/src/repository"
	userServices "github.com/strikersk/user-auth/src/service"
	"net/http"
)

func main() {
	myRouter := mux.NewRouter()

	jwtConfig := jwt.NewConfigStruct()
	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewCacheConfig()
	userService := userServices.NewLocalUserRepository(&userRepo, userCache)

	userHandling := handlers.NewUserHandler(&userService)
	jwtHandling := handlers.NewJwtHandler(&userService, jwtConfig)

	redisRouter := myRouter.PathPrefix("/redis").Subrouter()
	redisRouter.HandleFunc("/login", redis.Signin).Methods("POST")
	redisRouter.HandleFunc("/welcome", redis.Welcome).Methods("GET")

	userHandling.EnrichRouter(myRouter)
	jwtHandling.EnrichRouter(myRouter)

	fmt.Println(http.ListenAndServe(":4000", cors.AllowAll().Handler(myRouter)))
}
