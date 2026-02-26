package config

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/logger"
)

var (
	globalConfig *Config
	configOnce   sync.Once
)

type Config struct {
	App       *AppConfig
	DB        *DBConfig
	Cache     *CacheConfig
	Email     *EmailConfig
	Google    *GoogleOAuthConfig
	GitHub    *GitHubOAuthConfig
	Facebook  *FacebookOAuthConfig
	RateLimit *RateLimitConfig
	Session   *SessionConfig
	Link      *LinkConfig
}

func initViper() string {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return ""
	}
	return viper.ConfigFileUsed()
}

func NewConfig(ctx context.Context) *Config {
	configOnce.Do(func() {
		configFile := initViper()

		appConfig := LoadAppConfig(ctx)

		logger.Init(&logger.Config{
			Env:         appConfig.Env,
			LogFilePath: appConfig.LogFilePath,
		})

		if configFile != "" {
			logger.Info("Configuration loaded from file", zap.String("file", configFile))
		} else {
			logger.Warn("No config file found, using environment variables and defaults")
		}

		globalConfig = &Config{
			App:       appConfig,
			DB:        LoadDBConfig(ctx),
			Cache:     LoadCacheConfig(ctx),
			Email:     LoadEmailConfig(ctx),
			Google:    LoadGoogleOAuthConfig(ctx),
			GitHub:    LoadGitHubOAuthConfig(ctx),
			Facebook:  LoadFacebookOAuthConfig(ctx),
			RateLimit: LoadRateLimitConfig(ctx),
			Session:   LoadSessionConfig(ctx),
			Link:      LoadLinkConfig(ctx),
		}

		logger.Info("All configurations loaded successfully")
	})
	return globalConfig
}
