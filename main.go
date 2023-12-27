package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/handlers"
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

	authorizationService := userServices.NewJWTService(authorizationConfig)

	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewCacheConfig(applicationConfiguration.Cache)
	userService := userServices.NewUserService(&userRepo, userCache)
	userHandling := handlers.NewUserHandler(&userService)
	userHandling.EnrichRouter(appRoute)

	switch authorizationConfig.AuthorizationType {
	case "jwt":
		jwtHandling := handlers.NewJwtHandler(&userService, authorizationService)
		jwtHandling.EnrichRouter(appRoute)
		break
	case "cookies":
		cookiesHandling := handlers.NewCookiesHandler(&userService, authorizationService, authorizationConfig)
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
