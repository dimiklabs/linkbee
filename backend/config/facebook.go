package config

import (
	"context"

	"github.com/spf13/viper"
)

type FacebookOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Enabled      bool
}

func LoadFacebookOAuthConfig(_ context.Context) *FacebookOAuthConfig {
	viper.SetDefault("FACEBOOK_OAUTH_ENABLED", false)
	viper.SetDefault("FACEBOOK_CLIENT_ID", "")
	viper.SetDefault("FACEBOOK_CLIENT_SECRET", "")
	viper.SetDefault("FACEBOOK_REDIRECT_URL", "http://localhost:8080/api/v1/auth/facebook/callback")

	return &FacebookOAuthConfig{
		Enabled:      viper.GetBool("FACEBOOK_OAUTH_ENABLED"),
		ClientID:     viper.GetString("FACEBOOK_CLIENT_ID"),
		ClientSecret: viper.GetString("FACEBOOK_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("FACEBOOK_REDIRECT_URL"),
	}
}
