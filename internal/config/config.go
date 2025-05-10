package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort int    `mapstructure:"grpc_port"`
	RESTPort int    `mapstructure:"rest_port"`
	LogLevel string `mapstructure:"log_level"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
