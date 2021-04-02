package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/strikersk/user-auth/src"
	"github.com/strikersk/user-auth/src/jwt"
	"net/http"
)

func main() {
	src.InitCache()
	src.InitEnvFile()

	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/signin", src.Signin).Methods("POST")
	myRouter.HandleFunc("/welcome", src.Welcome).Methods("GET")
	myRouter.HandleFunc("/jwt/login", jwt.Login).Methods("POST")
	myRouter.HandleFunc("/jwt/welcome", jwt.Welcome).Methods("Get")

	fmt.Println(http.ListenAndServe(":5000", cors.AllowAll().Handler(myRouter)))
}
