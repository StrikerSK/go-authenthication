package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	appConfigs "github.com/strikersk/user-auth/config"
	userhandlers "github.com/strikersk/user-auth/src/handlers"
	userPorts "github.com/strikersk/user-auth/src/ports"
	userRepository "github.com/strikersk/user-auth/src/repository"
	userServices "github.com/strikersk/user-auth/src/service"
	"log"
	"net/http"
)

func main() {
	applicationConfiguration := appConfigs.ReadConfiguration()
	applicationConfig := applicationConfiguration.Application
	authorizationConfig := applicationConfiguration.Authorization
	cacheConfiguration := applicationConfiguration.Cache

	appRoute := mux.NewRouter().PathPrefix(applicationConfig.ContextPath).Subrouter()

	authorizationService := resolveEncodingType(authorizationConfig)

	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewRedisCache(cacheConfiguration)
	userService := userServices.NewUserService(&userRepo, userCache)
	userHandling := userhandlers.NewUserHandler(&userService)

	handlers := []userPorts.IUserHandler{
		userHandling,
	}

	switch authorizationConfig.AuthorizationType {
	case "jwt":
		log.Println("JWT endpoint handling  selected")
		jwtHandling := userhandlers.NewJwtHandler(&userService, authorizationService)
		handlers = append(handlers, jwtHandling)
		break
	case "cookies":
		log.Println("Cookies endpoint handling selected")
		cookiesHandling := userhandlers.NewCookiesHandler(&userService, authorizationService, authorizationConfig)
		handlers = append(handlers, cookiesHandling)
		break
	default:
		log.Fatal("unrecognized authorization type")
	}

	for _, handler := range handlers {
		handler.RegisterHandler(appRoute)
	}

	corsHandler := cors.AllowAll().Handler(appRoute)
	address := fmt.Sprintf(":%s", applicationConfig.Port)

	log.Println("Listening on port:", applicationConfig.Port)
	log.Println(http.ListenAndServe(address, corsHandler))
}

func resolveEncodingType(configuration appConfigs.Authorization) userPorts.IAuthorizationService {
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
