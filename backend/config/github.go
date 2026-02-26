package config

import (
	"context"

	"github.com/spf13/viper"
)

type GitHubOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Enabled      bool
}

func LoadGitHubOAuthConfig(_ context.Context) *GitHubOAuthConfig {
	viper.SetDefault("GITHUB_OAUTH_ENABLED", false)
	viper.SetDefault("GITHUB_CLIENT_ID", "")
	viper.SetDefault("GITHUB_CLIENT_SECRET", "")
	viper.SetDefault("GITHUB_REDIRECT_URL", "http://localhost:8080/api/v1/auth/github/callback")

	return &GitHubOAuthConfig{
		Enabled:      viper.GetBool("GITHUB_OAUTH_ENABLED"),
		ClientID:     viper.GetString("GITHUB_CLIENT_ID"),
		ClientSecret: viper.GetString("GITHUB_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("GITHUB_REDIRECT_URL"),
	}
}
