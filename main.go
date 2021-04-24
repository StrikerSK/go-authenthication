package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/strikersk/user-auth/src/jwt"
	"github.com/strikersk/user-auth/src/redis"
	"net/http"
)

func main() {
	myRouter := mux.NewRouter()

	redisRouter := myRouter.PathPrefix("/redis").Subrouter()
	redisRouter.HandleFunc("/login", redis.Signin).Methods("POST")
	redisRouter.HandleFunc("/welcome", redis.Welcome).Methods("GET")

	jwtRouter := myRouter.PathPrefix("/jwt").Subrouter()
	jwtRouter.HandleFunc("/login", jwt.Login).Methods("POST")
	jwtRouter.HandleFunc("/welcome", jwt.Welcome).Methods("Get")

	fmt.Println(http.ListenAndServe(":5000", cors.AllowAll().Handler(myRouter)))
}
