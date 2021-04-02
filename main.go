package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/strikersk/user-auth/src"
	"net/http"
)

func main() {
	src.InitCache()
	src.InitEnvFile()

	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/signin", src.Signin).Methods("POST")
	myRouter.HandleFunc("/welcome", src.Welcome).Methods("GET")
	myRouter.HandleFunc("/jwt/login", src.JwtLogin).Methods("POST")
	myRouter.HandleFunc("/jwt/welcome", src.JwtWelcome).Methods("Get")

	fmt.Println(http.ListenAndServe(":5000", myRouter))
}
