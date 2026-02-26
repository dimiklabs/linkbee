package config

import (
	"context"

	"github.com/spf13/viper"
)

type RateLimitConfig struct {
	Enabled bool

	LoginMaxAttempts    int
	LoginWindowMinutes  int
	LoginLockoutMinutes int

	IPMaxAttempts    int
	IPWindowMinutes  int
	IPLockoutMinutes int

	ProgressiveEnabled       bool
	ProgressiveMultiplier    float64
	ProgressiveMaxMultiplier float64
}

func LoadRateLimitConfig(_ context.Context) *RateLimitConfig {
	viper.SetDefault("RATE_LIMIT_ENABLED", true)

	viper.SetDefault("RATE_LIMIT_LOGIN_MAX_ATTEMPTS", 5)
	viper.SetDefault("RATE_LIMIT_LOGIN_WINDOW_MINUTES", 15)
	viper.SetDefault("RATE_LIMIT_LOGIN_LOCKOUT_MINUTES", 15)

	viper.SetDefault("RATE_LIMIT_IP_MAX_ATTEMPTS", 20)
	viper.SetDefault("RATE_LIMIT_IP_WINDOW_MINUTES", 15)
	viper.SetDefault("RATE_LIMIT_IP_LOCKOUT_MINUTES", 30)

	viper.SetDefault("RATE_LIMIT_PROGRESSIVE_ENABLED", true)
	viper.SetDefault("RATE_LIMIT_PROGRESSIVE_MULTIPLIER", 2.0)
	viper.SetDefault("RATE_LIMIT_PROGRESSIVE_MAX_MULTIPLIER", 8.0)

	return &RateLimitConfig{
		Enabled: viper.GetBool("RATE_LIMIT_ENABLED"),

		LoginMaxAttempts:    viper.GetInt("RATE_LIMIT_LOGIN_MAX_ATTEMPTS"),
		LoginWindowMinutes:  viper.GetInt("RATE_LIMIT_LOGIN_WINDOW_MINUTES"),
		LoginLockoutMinutes: viper.GetInt("RATE_LIMIT_LOGIN_LOCKOUT_MINUTES"),

		IPMaxAttempts:    viper.GetInt("RATE_LIMIT_IP_MAX_ATTEMPTS"),
		IPWindowMinutes:  viper.GetInt("RATE_LIMIT_IP_WINDOW_MINUTES"),
		IPLockoutMinutes: viper.GetInt("RATE_LIMIT_IP_LOCKOUT_MINUTES"),

		ProgressiveEnabled:       viper.GetBool("RATE_LIMIT_PROGRESSIVE_ENABLED"),
		ProgressiveMultiplier:    viper.GetFloat64("RATE_LIMIT_PROGRESSIVE_MULTIPLIER"),
		ProgressiveMaxMultiplier: viper.GetFloat64("RATE_LIMIT_PROGRESSIVE_MAX_MULTIPLIER"),
	}
}
