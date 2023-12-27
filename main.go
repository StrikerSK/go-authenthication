package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/handlers"
	"github.com/strikersk/user-auth/src/jwt"
	userRepository "github.com/strikersk/user-auth/src/repository"
	userServices "github.com/strikersk/user-auth/src/service"
	"log"
	"net/http"
)

func main() {
	applicationConfiguration := config.ReadConfiguration()

	coreRoute := mux.NewRouter()
	appRoute := coreRoute.PathPrefix(applicationConfiguration.Application.ContextPath).Subrouter()

	jwtConfig := jwt.NewConfigStruct()
	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewCacheConfig()
	userService := userServices.NewLocalUserRepository(&userRepo, userCache)

	userHandling := handlers.NewUserHandler(&userService)
	jwtHandling := handlers.NewJwtHandler(&userService, jwtConfig)
	cookiesHandling := handlers.NewCookiesHandler("session_token", &userService)

	userHandling.EnrichRouter(appRoute)
	jwtHandling.EnrichRouter(appRoute)
	cookiesHandling.EnrichRouter(appRoute)

	corsHandler := cors.AllowAll().Handler(appRoute)
	address := fmt.Sprintf(":%s", applicationConfiguration.Application.Port)

	log.Println("Listening on port:", applicationConfiguration.Application.Port)
	log.Println(http.ListenAndServe(address, corsHandler))
}
