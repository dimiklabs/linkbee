package config

import (
	"context"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Env                string
	ServerAddress      string
	CorsAllowedOrigins string
	LogFilePath        string
	JWTSecret          string
	JWTIssuer          string
	JWTAccessExpiry    int
	JWTRefreshExpiry   int
	MaxSessions        int
	BaseDomain         string
	GeoDBPath          string
}

func LoadAppConfig(_ context.Context) *AppConfig {
	viper.SetDefault("ENV", "development")
	viper.SetDefault("SERVER_ADDRESS", ":8080")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	viper.SetDefault("LOG_FILE_PATH", "")
	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_ISSUER", "shortlink-auth")
	viper.SetDefault("JWT_ACCESS_EXPIRY_MINUTES", 15)
	viper.SetDefault("JWT_REFRESH_EXPIRY_DAYS", 7)
	viper.SetDefault("MAX_SESSIONS", 5)
	viper.SetDefault("BASE_DOMAIN", "http://localhost:8080")
	viper.SetDefault("GEO_DB_PATH", "")

	return &AppConfig{
		Env:                viper.GetString("ENV"),
		ServerAddress:      viper.GetString("SERVER_ADDRESS"),
		CorsAllowedOrigins: viper.GetString("CORS_ALLOWED_ORIGINS"),
		LogFilePath:        viper.GetString("LOG_FILE_PATH"),
		JWTSecret:          viper.GetString("JWT_SECRET"),
		JWTIssuer:          viper.GetString("JWT_ISSUER"),
		JWTAccessExpiry:    viper.GetInt("JWT_ACCESS_EXPIRY_MINUTES"),
		JWTRefreshExpiry:   viper.GetInt("JWT_REFRESH_EXPIRY_DAYS"),
		MaxSessions:        viper.GetInt("MAX_SESSIONS"),
		BaseDomain:         viper.GetString("BASE_DOMAIN"),
		GeoDBPath:          viper.GetString("GEO_DB_PATH"),
	}
}
