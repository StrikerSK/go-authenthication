package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/src/handlers"
	"github.com/strikersk/user-auth/src/ports"
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

	authorizationService := resolveEncodingType(authorizationConfig)

	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewCacheConfig(applicationConfiguration.Cache)
	userService := userServices.NewUserService(&userRepo, userCache)
	userHandling := handlers.NewUserHandler(&userService)

	userHandling.RegisterHandler(appRoute)

	switch authorizationConfig.AuthorizationType {
	case "jwt":
		log.Println("JWT endpoint handling  selected")
		jwtHandling := handlers.NewJwtHandler(&userService, authorizationService)
		jwtHandling.RegisterHandler(appRoute)
		break
	case "cookies":
		log.Println("Cookies endpoint handling  selected")
		cookiesHandling := handlers.NewCookiesHandler(&userService, authorizationService, authorizationConfig)
		cookiesHandling.RegisterHandler(appRoute)
		break
	default:
		log.Fatal("unrecognized authorization type")
	}

	corsHandler := cors.AllowAll().Handler(appRoute)
	address := fmt.Sprintf(":%s", applicationConfig.Port)

	log.Println("Listening on port:", applicationConfig.Port)
	log.Println(http.ListenAndServe(address, corsHandler))
}

func resolveEncodingType(configuration config.Authorization) ports.IAuthorizationService {
	switch configuration.TokenEncodingType {
	case "jwt":
		log.Println("JWT Token encoding selected")
		return userServices.NewJWTService(configuration)
	case "base64":
		log.Println("Base64 Token encoding selected")
		return userServices.NewBase64EncodingService()
	default:
		log.Fatal("no token encoding type selected")
		return nil
	}
}
