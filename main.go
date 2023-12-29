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
	userEndpointAuthorization := resolveUserEndpointAuthorization(authorizationConfig)
	userPasswordService := userServices.NewBcryptUserPasswordService(authorizationConfig.Encryption.Cost)
	userService := userServices.NewUserService(userRepo, userCache, userPasswordService)

	handlers := []userPorts.IUserHandler{
		userhandlers.NewUserHandler(userService, encodingService, userEndpointAuthorization),
	}

	for _, handler := range handlers {
		handler.RegisterHandler(applicationRouter)
	}

	corsHandler := cors.AllowAll().Handler(applicationRouter)
	address := fmt.Sprintf(":%s", applicationConfig.Port)

	log.Println("Listening on port:", applicationConfig.Port)
	log.Println(http.ListenAndServe(address, corsHandler))
}

func resolveEncodingType(configuration appConfigs.Authorization) userPorts.IEncodingService {
	switch configuration.TokenEncodingType {
	case "jwt":
		log.Println("JWT Token encoding selected")
		return userServices.NewJWTEncodingService(configuration)
	case "base64":
		log.Println("Base64 Token encoding selected")
		return userServices.NewBase64EncodingService()
	default:
		log.Fatal("no token encoding type selected")
		return nil
	}
}

func resolveUserEndpointAuthorization(authorizationConfig appConfigs.Authorization) userPorts.IUserEndpointHandler {
	switch authorizationConfig.AuthorizationType {
	case "header":
		log.Println("Header authorization handling selected")
		return userhandlers.NewHeaderAuthorization(authorizationConfig)
	case "cookies":
		log.Println("Cookies authorization handling selected")
		return userhandlers.NewCookiesAuthorization(authorizationConfig)
	default:
		log.Fatal("No authorization handling selected")
		return nil
	}
}

func resolveCachingInstance(configuration appConfigs.CacheConfiguration) userPorts.IUserCache {
	switch configuration.Name {
	case "memcache":
		log.Println("Memcache instance selected")
		return userRepository.NewMemcacheCache(configuration)
	case "redis":
		log.Println("Redis instance selected")
		return userRepository.NewRedisCache(configuration)
	default:
		log.Fatal("No cache instance selected")
		return nil
	}
}

func resolveDatabaseInstance() userPorts.IUserRepository {
	return userRepository.NewLocalUserRepository()
}
