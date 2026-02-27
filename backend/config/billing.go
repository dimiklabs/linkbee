package config

import "github.com/spf13/viper"

type BillingConfig struct {
	LemonSqueezyWebhookSecret string
	LemonSqueezyStoreID       string
	ProVariantID              string
	BusinessVariantID         string
	ProCheckoutURL            string
	BusinessCheckoutURL       string
}

func LoadBillingConfig() *BillingConfig {
	viper.SetDefault("LEMONSQUEEZY_WEBHOOK_SECRET", "")
	viper.SetDefault("LEMONSQUEEZY_STORE_ID", "")
	viper.SetDefault("LEMONSQUEEZY_PRO_VARIANT_ID", "")
	viper.SetDefault("LEMONSQUEEZY_BUSINESS_VARIANT_ID", "")
	viper.SetDefault("LEMONSQUEEZY_PRO_CHECKOUT_URL", "")
	viper.SetDefault("LEMONSQUEEZY_BUSINESS_CHECKOUT_URL", "")

	return &BillingConfig{
		LemonSqueezyWebhookSecret: viper.GetString("LEMONSQUEEZY_WEBHOOK_SECRET"),
		LemonSqueezyStoreID:       viper.GetString("LEMONSQUEEZY_STORE_ID"),
		ProVariantID:              viper.GetString("LEMONSQUEEZY_PRO_VARIANT_ID"),
		BusinessVariantID:         viper.GetString("LEMONSQUEEZY_BUSINESS_VARIANT_ID"),
		ProCheckoutURL:            viper.GetString("LEMONSQUEEZY_PRO_CHECKOUT_URL"),
		BusinessCheckoutURL:       viper.GetString("LEMONSQUEEZY_BUSINESS_CHECKOUT_URL"),
	}
}
