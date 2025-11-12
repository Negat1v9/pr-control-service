package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	AppConfig
	WebConfig
	PostgresConfig
}

type AppConfig struct {
	Env string
}

type WebConfig struct {
	ListenAddress string
	ReadTimeout   int64
	WriteTimeout  int64
}

type PostgresConfig struct {
	DbHost     string
	DbPort     int
	DbName     string
	DbUser     string
	DbPassword string
	DbSslMode  string
}

func parseCfg(fileName string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(fileName)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

func LoadConfig(fileName string) (*Config, error) {
	v, err := parseCfg(fileName)
	if err != nil {
		return nil, err
	}
	var cfg *Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
