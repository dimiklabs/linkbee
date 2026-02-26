package config

import (
	"context"

	"github.com/spf13/viper"
)

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	TTL      int
}

func LoadCacheConfig(_ context.Context) *CacheConfig {
	viper.SetDefault("VALKEY_HOST", "localhost")
	viper.SetDefault("VALKEY_PORT", "6379")
	viper.SetDefault("VALKEY_PASSWORD", "")
	viper.SetDefault("VALKEY_DB", 0)
	viper.SetDefault("VALKEY_TTL", 3600)

	return &CacheConfig{
		Host:     viper.GetString("VALKEY_HOST"),
		Port:     viper.GetString("VALKEY_PORT"),
		Password: viper.GetString("VALKEY_PASSWORD"),
		DB:       viper.GetInt("VALKEY_DB"),
		TTL:      viper.GetInt("VALKEY_TTL"),
	}
}
