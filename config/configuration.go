package config

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
