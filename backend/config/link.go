package config

import (
	"context"

	"github.com/spf13/viper"
)

type LinkConfig struct {
	SlugLength              int
	SlugSecret              string
	DefaultRedirectType     int
	CacheTTLSeconds         int
	ClickQueueBatchSize     int
	ClickQueueFlushInterval int
	DemoRateLimitPerIP      int
}

func LoadLinkConfig(_ context.Context) *LinkConfig {
	viper.SetDefault("LINK_SLUG_LENGTH", 5)
	viper.SetDefault("LINK_SLUG_SECRET", "")
	viper.SetDefault("LINK_DEFAULT_REDIRECT_TYPE", 302)
	viper.SetDefault("LINK_CACHE_TTL_SECONDS", 86400)
	viper.SetDefault("CLICK_QUEUE_BATCH_SIZE", 100)
	viper.SetDefault("CLICK_QUEUE_FLUSH_INTERVAL_SECONDS", 5)
	viper.SetDefault("DEMO_RATE_LIMIT_PER_IP", 5)

	return &LinkConfig{
		SlugLength:              viper.GetInt("LINK_SLUG_LENGTH"),
		SlugSecret:              viper.GetString("LINK_SLUG_SECRET"),
		DefaultRedirectType:     viper.GetInt("LINK_DEFAULT_REDIRECT_TYPE"),
		CacheTTLSeconds:         viper.GetInt("LINK_CACHE_TTL_SECONDS"),
		ClickQueueBatchSize:     viper.GetInt("CLICK_QUEUE_BATCH_SIZE"),
		ClickQueueFlushInterval: viper.GetInt("CLICK_QUEUE_FLUSH_INTERVAL_SECONDS"),
		DemoRateLimitPerIP:      viper.GetInt("DEMO_RATE_LIMIT_PER_IP"),
	}
}
