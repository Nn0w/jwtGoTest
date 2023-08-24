package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDBUri             string        `mapstructure:"MONGODB_URI"`
	MongoDBName            string        `mapstructure:"MONGODB_NAME"`
	Port                   string        `mapstructure:"PORT"`
	AccessTokenHSKey       string        `mapstructu:"ACCESS_TOKEN_HS_KEY"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenHSKey      string        `mapstructu:"REFRESH_TOKEN_HS_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func GetConfigOrPanic() (config Config) {
	config, err := LoadConfig(".")
	if err != nil {
		log.Panic("Could not load environment variables: ", err)
	}
	return config
}
