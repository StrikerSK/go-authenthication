package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	appConfigs "github.com/strikersk/user-auth/config"
	appConstants "github.com/strikersk/user-auth/constants"
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

	databaseRepository := resolveDatabaseInstance(applicationConfiguration.Database)
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
	case appConstants.JWT:
		log.Println("JWT Token encoding selected")
		return userServices.NewJWTEncodingService(configuration)
	case appConstants.Base64:
		log.Println("Base64 Token encoding selected")
		return userServices.NewBase64EncodingService()
	default:
		log.Fatalf("no token encoding type selected, available options: %s, %s", appConstants.JWT, appConstants.Base64)
		return nil
	}
}

func resolveUserEndpointAuthorization(authorizationConfig appConfigs.Authorization) userPorts.IUserEndpointHandler {
	switch authorizationConfig.AuthorizationType {
	case appConstants.Header:
		log.Println("Header authorization handling selected")
		return userhandlers.NewHeaderAuthorization(authorizationConfig)
	case appConstants.Cookies:
		log.Println("Cookies authorization handling selected")
		return userhandlers.NewCookiesAuthorization(authorizationConfig)
	default:
		log.Fatalf("No authorization handling selected, available options: %s, %s", appConstants.Header, appConstants.Cookies)
		return nil
	}
}

func resolveCachingInstance(configuration appConfigs.CacheConfiguration) userPorts.IUserCache {
	switch configuration.Name {
	case appConstants.InMemory:
		log.Println("InMemory cache instance selected")
		return caching.NewInMemoryCache()
	case appConstants.MemCache:
		log.Println("MemCache instance selected")
		return caching.NewMemcacheCache(configuration)
	case appConstants.Redis:
		log.Println("Redis instance selected")
		return caching.NewRedisCache(configuration)
	default:
		log.Fatalf("No cache instance selected, available caches: %s, %s, %s", appConstants.MemCache, appConstants.Redis, appConstants.InMemory)
		return nil
	}
}

func resolvePasswordService(configuration *appConfigs.EncryptionConfiguration) userPorts.IUserPasswordService {
	return userServices.NewBcryptUserPasswordService(configuration)
}

func resolveDatabaseInstance(configuration appConfigs.DatabaseConfiguration) userPorts.IUserRepository {
	switch configuration.Name {
	case appConstants.SQLite, appConstants.Postgres:
		log.Println("SQLite database instance selected")
		return database.NewGormUserRepository(configuration)
	case appConstants.InMemory:
		log.Println("InMemory database instance selected")
		return database.NewInMemoryUserDatabase()
	default:
		log.Fatalf("No database instance selected, available options are: %s, %s, %s", appConstants.InMemory, appConstants.SQLite, appConstants.Postgres)
		return nil
	}
}
