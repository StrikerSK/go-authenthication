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
	applicationConfig := applicationConfiguration.Application
	authorizationConfig := applicationConfiguration.Authorization
	appRoute := mux.NewRouter().PathPrefix(applicationConfig.ContextPath).Subrouter()

	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewCacheConfig()
	userService := userServices.NewLocalUserRepository(&userRepo, userCache)
	userHandling := handlers.NewUserHandler(&userService)
	userHandling.EnrichRouter(appRoute)

	switch authorizationConfig.AuthorizationType {
	case "jwt":
		jwtConfig := jwt.NewConfigStruct()
		jwtHandling := handlers.NewJwtHandler(&userService, jwtConfig)
		jwtHandling.EnrichRouter(appRoute)
		break
	case "cookies":
		cookiesHandling := handlers.NewCookiesHandler("session_token", &userService)
		cookiesHandling.EnrichRouter(appRoute)
		break
	default:
		log.Fatal("unrecognized authorization type")
	}

	corsHandler := cors.AllowAll().Handler(appRoute)
	address := fmt.Sprintf(":%s", applicationConfig.Port)

	log.Println("Listening on port:", applicationConfig.Port)
	log.Println(http.ListenAndServe(address, corsHandler))
}
