package config

import (
	"context"

	"github.com/spf13/viper"
)

type SessionConfig struct {
	RememberMeEnabled           bool
	RememberMeRefreshExpiryDays int
	TrackActivity               bool
	ActivityUpdateIntervalMins  int
	NotifyOnNewSession          bool
	NotifyOnNewSessionEmail     bool
	NotifyOnSuspiciousActivity  bool
}

func LoadSessionConfig(_ context.Context) *SessionConfig {
	viper.SetDefault("SESSION_REMEMBER_ME_ENABLED", true)
	viper.SetDefault("SESSION_REMEMBER_ME_REFRESH_EXPIRY_DAYS", 30)
	viper.SetDefault("SESSION_TRACK_ACTIVITY", true)
	viper.SetDefault("SESSION_ACTIVITY_UPDATE_INTERVAL_MINS", 5)
	viper.SetDefault("SESSION_NOTIFY_ON_NEW_SESSION", true)
	viper.SetDefault("SESSION_NOTIFY_ON_NEW_SESSION_EMAIL", false)
	viper.SetDefault("SESSION_NOTIFY_ON_SUSPICIOUS_ACTIVITY", true)

	return &SessionConfig{
		RememberMeEnabled:           viper.GetBool("SESSION_REMEMBER_ME_ENABLED"),
		RememberMeRefreshExpiryDays: viper.GetInt("SESSION_REMEMBER_ME_REFRESH_EXPIRY_DAYS"),
		TrackActivity:               viper.GetBool("SESSION_TRACK_ACTIVITY"),
		ActivityUpdateIntervalMins:  viper.GetInt("SESSION_ACTIVITY_UPDATE_INTERVAL_MINS"),
		NotifyOnNewSession:          viper.GetBool("SESSION_NOTIFY_ON_NEW_SESSION"),
		NotifyOnNewSessionEmail:     viper.GetBool("SESSION_NOTIFY_ON_NEW_SESSION_EMAIL"),
		NotifyOnSuspiciousActivity:  viper.GetBool("SESSION_NOTIFY_ON_SUSPICIOUS_ACTIVITY"),
	}
}
