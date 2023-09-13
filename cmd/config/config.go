package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	ServerAddress string        `mapstructure:"server_address"`
	ReadTimeout   time.Duration `mapstructure:"read_timeout"`
	WriteTimeout  time.Duration `mapstructure:"write_timeout"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &config, nil

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
