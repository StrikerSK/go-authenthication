package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	appConfigs "github.com/strikersk/user-auth/config"
	userhandlers "github.com/strikersk/user-auth/src/handlers"
	userPorts "github.com/strikersk/user-auth/src/ports"
	caching "github.com/strikersk/user-auth/src/repository/cache"
	"github.com/strikersk/user-auth/src/repository/database"
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

	databaseRepository := resolveDatabaseInstance()
	cacheRepository := resolveCachingInstance(cacheConfiguration)

	userEndpointAuthorization := resolveUserEndpointAuthorization(authorizationConfig)
	userPasswordService := resolvePasswordService(&authorizationConfig.Encryption)
	userService := userServices.NewUserService(databaseRepository, cacheRepository, userPasswordService)

	handlers := []userPorts.IUserHandler{
		userhandlers.NewUserHandler(userService, encodingService, userEndpointAuthorization),
	}

	for _, handler := range handlers {
		handler.RegisterHandler(applicationRouter)
	}

	corsHandler := cors.AllowAll().Handler(applicationRouter)
	address := fmt.Sprintf(":%s", applicationConfig.Port)

	log.Println("Listening on port:", applicationConfig.Port)

	err := http.ListenAndServe(address, corsHandler)
	if err != nil {
		log.Fatal(err)
	}
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
		return caching.NewMemcacheCache(configuration)
	case "redis":
		log.Println("Redis instance selected")
		return caching.NewRedisCache(configuration)
	default:
		log.Fatal("No cache instance selected")
		return nil
	}
}

func resolvePasswordService(configuration *appConfigs.EncryptionConfiguration) userPorts.IUserPasswordService {
	return userServices.NewBcryptUserPasswordService(configuration)
}

func resolveDatabaseInstance() userPorts.IUserRepository {
	return database.NewLocalUserRepository()
}
