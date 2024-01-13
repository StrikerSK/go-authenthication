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

	applicationRouter := mux.NewRouter().PathPrefix(applicationConfig.ContextPath).Subrouter()
	encodingService := resolveEncodingType(authorizationConfig)

	userRepo := userRepository.NewLocalUserRepository()
	userCache := userRepository.NewRedisCache(cacheConfiguration)
	userService := userServices.NewUserService(&userRepo, userCache)
	userRegisterHandler := userhandlers.NewUserHandler(userService)

	userAuthorization := resolveUserAuthorization(authorizationConfig)
	userAccessHandler := userhandlers.NewAbstractHandler(userService, encodingService, userAuthorization)

	handlers := []userPorts.IUserHandler{
		userAccessHandler,
		userRegisterHandler,
	}

	for _, handler := range handlers {
		handler.RegisterHandler(applicationRouter)
	}

	corsHandler := cors.AllowAll().Handler(applicationRouter)
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

func resolveUserAuthorization(authorizationConfig appConfigs.Authorization) userPorts.IUserEndpointHandler {
	switch authorizationConfig.TokenEncodingType {
	case "jwt":
		log.Println("JWT Token handling selected")
		return userhandlers.NewJwtHandler(authorizationConfig)
	case "cookies":
		log.Println("Cookies handling selected")
		return userhandlers.NewCookiesHandler(authorizationConfig)
	default:
		log.Fatal("Authorization handling selected")
		return nil
	}
}
