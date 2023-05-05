package util

import (
	"github.com/spf13/viper"
)

// Config is a struct that holds the configuration for the application.
type Config struct {
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	TokenSymmetricKey string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
}

func LoadConfig(configPath string) (Config, error) {
	var config Config
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return config, err
	}

	if err := v.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
