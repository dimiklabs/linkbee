package config

import "github.com/spf13/viper"

type BillingConfig struct {
	PaddleAPIKey          string
	PaddleWebhookSecret   string
	PaddleEnvironment     string // "production" or "sandbox"
	PaddleProPriceID      string
	PaddleGrowthPriceID   string
}

func LoadBillingConfig() *BillingConfig {
	viper.SetDefault("PADDLE_API_KEY", "")
	viper.SetDefault("PADDLE_WEBHOOK_SECRET", "")
	viper.SetDefault("PADDLE_ENVIRONMENT", "sandbox")
	viper.SetDefault("PADDLE_PRO_PRICE_ID", "")
	viper.SetDefault("PADDLE_GROWTH_PRICE_ID", "")

	return &BillingConfig{
		PaddleAPIKey:        viper.GetString("PADDLE_API_KEY"),
		PaddleWebhookSecret: viper.GetString("PADDLE_WEBHOOK_SECRET"),
		PaddleEnvironment:   viper.GetString("PADDLE_ENVIRONMENT"),
		PaddleProPriceID:    viper.GetString("PADDLE_PRO_PRICE_ID"),
		PaddleGrowthPriceID: viper.GetString("PADDLE_GROWTH_PRICE_ID"),
	}
}
