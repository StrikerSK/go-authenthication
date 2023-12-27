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

	var authorizationService ports.IAuthorizationService
	switch authorizationConfig.TokenEncodingType {
	case "jwt":
		log.Println("JWT Token encoding selected")
		authorizationService = userServices.NewJWTService(authorizationConfig)
		break
	case "base64":
		log.Println("Base64 Token encoding selected")
		authorizationService = userServices.NewBase64EncodingService()
		break
	default:
		log.Println("no token encoding type selected")
	}

	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewCacheConfig(applicationConfiguration.Cache)
	userService := userServices.NewUserService(&userRepo, userCache)
	userHandling := handlers.NewUserHandler(&userService)
	userHandling.EnrichRouter(appRoute)

	switch authorizationConfig.AuthorizationType {
	case "jwt":
		log.Println("JWT endpoint handling  selected")
		jwtHandling := handlers.NewJwtHandler(&userService, authorizationService)
		jwtHandling.EnrichRouter(appRoute)
		break
	case "cookies":
		log.Println("Cookies endpoint handling  selected")
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
