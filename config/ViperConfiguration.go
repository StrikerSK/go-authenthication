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
	TokenEncoding string `mapstructure:"TokenEncoding"`
}

type Authorization struct {
	AuthorizationType   string           `mapstructure:"AuthorizationType"`
	AuthorizationHeader string           `mapstructure:"AuthorizationHeader"`
	TokenExpiration     int              `mapstructure:"TokenExpiration"`
	ExcludedPaths       []string         `mapstructure:"ExcludedPaths"`
	JWT                 JWTConfiguration `mapstructure:"JWT"`
}

type CacheConfiguration struct {
	URL        string `mapstructure:"URL"`
	Host       string `mapstructure:"Host"`
	Port       string `mapstructure:"Port"`
	Expiration int    `mapstructure:"Expiration"`
}

type ApplicationConfiguration struct {
	Application   Application        `mapstructure:"Application"`
	Authorization Authorization      `mapstructure:"Authorization"`
	Cache         CacheConfiguration `mapstructure:"Cache"`
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
			AuthorizationType:   "cookies",
			AuthorizationHeader: "Authorization",
			TokenExpiration:     3600,
			JWT: JWTConfiguration{
				TokenEncoding: "Wow, much safe",
			},
		},
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		os.Exit(-1)
	}

	return configuration
}
