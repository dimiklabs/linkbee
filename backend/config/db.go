package config

import (
	"context"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	MaxOpenConnection int
	MaxIdleConnection int

	ReadHost     string
	ReadPort     string
	ReadUser     string
	ReadPassword string
}

func LoadDBConfig(_ context.Context) *DBConfig {
	viper.SetDefault("DB_HOST_WRITE", "localhost")
	viper.SetDefault("DB_PORT_WRITE", "5432")
	viper.SetDefault("DB_USER_WRITE", "linkbee")
	viper.SetDefault("DB_PASSWORD_WRITE", "linkbee")
	viper.SetDefault("DB_NAME", "linkbee")
	viper.SetDefault("DB_MAX_OPEN_CONNECTION", 10)
	viper.SetDefault("DB_MAX_IDLE_CONNECTION", 5)

	viper.SetDefault("DB_HOST_READ", "localhost")
	viper.SetDefault("DB_PORT_READ", "5432")
	viper.SetDefault("DB_USER_READ", "linkbee")
	viper.SetDefault("DB_PASSWORD_READ", "linkbee")

	return &DBConfig{
		Host:              viper.GetString("DB_HOST_WRITE"),
		Port:              viper.GetString("DB_PORT_WRITE"),
		User:              viper.GetString("DB_USER_WRITE"),
		Password:          viper.GetString("DB_PASSWORD_WRITE"),
		Name:              viper.GetString("DB_NAME"),
		MaxOpenConnection: viper.GetInt("DB_MAX_OPEN_CONNECTION"),
		MaxIdleConnection: viper.GetInt("DB_MAX_IDLE_CONNECTION"),

		ReadHost:     viper.GetString("DB_HOST_READ"),
		ReadPort:     viper.GetString("DB_PORT_READ"),
		ReadUser:     viper.GetString("DB_USER_READ"),
		ReadPassword: viper.GetString("DB_PASSWORD_READ"),
	}
}
