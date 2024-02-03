package config

import (
	"github.com/spf13/viper"
	"github.com/strikersk/user-auth/constants"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Application struct {
	Port        string `mapstructure:"Port"`
	ContextPath string `mapstructure:"ContextPath"`
}

type JWTConfiguration struct {
	TokenEncoding string `mapstructure:"TokenEncoding"`
}

type Authorization struct {
	AuthorizationType   string                  `mapstructure:"AuthorizationType"`
	AuthorizationHeader string                  `mapstructure:"AuthorizationHeader"`
	TokenExpiration     int                     `mapstructure:"TokenExpiration"`
	TokenEncodingType   string                  `mapstructure:"TokenEncodingType"`
	ExcludedPaths       []string                `mapstructure:"ExcludedPaths"`
	JWT                 JWTConfiguration        `mapstructure:"JWT"`
	Encryption          EncryptionConfiguration `mapstructure:"Encryption"`
}

type EncryptionConfiguration struct {
	Cost int `mapstructure:"Cost"`
}

type CacheConfiguration struct {
	Name       string `mapstructure:"Name"`
	URL        string `mapstructure:"URL"`
	Host       string `mapstructure:"Host"`
	Port       string `mapstructure:"Port"`
	Expiration int    `mapstructure:"Expiration"`
}

type DatabaseConfiguration struct {
	Type     string `mapstructure:"Type"`
	URL      string `mapstructure:"URL"`
	Host     string `mapstructure:"Host"`
	Port     string `mapstructure:"Port"`
	Name     string `mapstructure:"Name"`
	Username string `mapstructure:"Username"`
	Password string `mapstructure:"Password"`
}

type ApplicationConfiguration struct {
	Application   Application           `mapstructure:"Application"`
	Authorization Authorization         `mapstructure:"Authorization"`
	Cache         CacheConfiguration    `mapstructure:"Cache"`
	Database      DatabaseConfiguration `mapstructure:"Database"`
}

// ReadConfiguration - read file from the current directory and marshal into the conf config struct.
func ReadConfiguration() *ApplicationConfiguration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Configuration reading error", err)
	}

	configuration := &ApplicationConfiguration{
		Application: Application{
			Port:        "8080",
			ContextPath: "/api",
		},
		Authorization: Authorization{
			AuthorizationType:   constants.Cookies,
			AuthorizationHeader: "Authorization",
			TokenEncodingType:   constants.JWT,
			TokenExpiration:     3600,
			JWT: JWTConfiguration{
				TokenEncoding: "Wow, much safe",
			},
			Encryption: EncryptionConfiguration{
				Cost: bcrypt.DefaultCost,
			},
		},
		Cache: CacheConfiguration{
			Name:       constants.InMemory,
			Expiration: 3600,
		},
		Database: DatabaseConfiguration{
			Type: constants.InMemory,
		},
	}

	err = viper.Unmarshal(configuration)
	if err != nil {
		log.Fatal("Unable to decode into configuration:", err)
	}

	return configuration
}
