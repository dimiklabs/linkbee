package config

import (
	"context"

	"github.com/spf13/viper"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Enabled      bool
}

func LoadGoogleOAuthConfig(_ context.Context) *GoogleOAuthConfig {
	viper.SetDefault("GOOGLE_OAUTH_ENABLED", false)
	viper.SetDefault("GOOGLE_CLIENT_ID", "")
	viper.SetDefault("GOOGLE_CLIENT_SECRET", "")
	viper.SetDefault("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback")

	return &GoogleOAuthConfig{
		Enabled:      viper.GetBool("GOOGLE_OAUTH_ENABLED"),
		ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
	}
}
