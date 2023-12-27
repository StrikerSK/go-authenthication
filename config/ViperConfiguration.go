package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Application struct {
	Port        string `mapstructure:"Port"`
	ContextPath string `mapstructure:"ContextPath"`
}

type JWTConfiguration struct {
	TokenExpiration     int    `mapstructure:"TokenExpiration"`
	TokenEncoding       string `mapstructure:"TokenEncoding"`
	AuthorizationHeader string `mapstructure:"AuthorizationHeader"`
}

type CookiesConfiguration struct {
	TokenExpiration     int    `mapstructure:"TokenExpiration"`
	AuthorizationHeader string `mapstructure:"AuthorizationHeader"`
}

type Authorization struct {
	AuthorizationType string           `mapstructure:"AuthorizationType"`
	ExcludedPaths     []string         `mapstructure:"ExcludedPaths"`
	JWT               JWTConfiguration `mapstructure:"JWT"`
}

type ApplicationConfiguration struct {
	Application   Application   `mapstructure:"Application"`
	Authorization Authorization `mapstructure:"Authorization"`
}

// ReadConfiguration - read file from the current directory and marshal into the conf config struct.
func ReadConfiguration() *ApplicationConfiguration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(-1)
	}

	configuration := &ApplicationConfiguration{
		Application: Application{
			Port:        "8080",
			ContextPath: "/api",
		},
		Authorization: Authorization{
			AuthorizationHeader: "Authorization",
		},
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		os.Exit(-1)
	}

	return configuration
}