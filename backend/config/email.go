package config

import (
	"context"

	"github.com/spf13/viper"
)

type EmailConfig struct {
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	FromEmail      string
	FromName       string
	BaseURL        string
	VerifyTokenTTL int
}

func LoadEmailConfig(_ context.Context) *EmailConfig {
	viper.SetDefault("SMTP_HOST", "localhost")
	viper.SetDefault("SMTP_PORT", 1025)
	viper.SetDefault("SMTP_USER", "")
	viper.SetDefault("SMTP_PASSWORD", "")
	viper.SetDefault("EMAIL_FROM_ADDRESS", "noreply@shortlink.io")
	viper.SetDefault("EMAIL_FROM_NAME", "Shortlink")
	viper.SetDefault("APP_BASE_URL", "http://localhost:3000")
	viper.SetDefault("EMAIL_VERIFY_TOKEN_TTL_HOURS", 24)

	return &EmailConfig{
		SMTPHost:       viper.GetString("SMTP_HOST"),
		SMTPPort:       viper.GetInt("SMTP_PORT"),
		SMTPUser:       viper.GetString("SMTP_USER"),
		SMTPPassword:   viper.GetString("SMTP_PASSWORD"),
		FromEmail:      viper.GetString("EMAIL_FROM_ADDRESS"),
		FromName:       viper.GetString("EMAIL_FROM_NAME"),
		BaseURL:        viper.GetString("APP_BASE_URL"),
		VerifyTokenTTL: viper.GetInt("EMAIL_VERIFY_TOKEN_TTL_HOURS"),
	}
}
