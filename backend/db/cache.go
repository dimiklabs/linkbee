package db

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeycompat"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/logger"
)

func NewValkeyClient(ctx context.Context, cfg *config.CacheConfig) (valkeycompat.Cmdable, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	logger.Debug("Attempting to connect to Valkey",
		zap.String("address", addr),
		zap.Int("db", cfg.DB))

	clientOpts := valkey.ClientOption{
		InitAddress:  []string{addr},
		SelectDB:     cfg.DB,
		DisableCache: true,
	}

	if cfg.Password != "" {
		clientOpts.Password = cfg.Password
	}

	client, err := valkey.NewClient(clientOpts)
	if err != nil {
		logger.Error("Failed to create Valkey client",
			zap.String("address", addr),
			zap.Error(err))
		return nil, fmt.Errorf("failed to create Valkey client: %w", err)
	}

	compat := valkeycompat.NewAdapter(client)

	if err := compat.Ping(ctx).Err(); err != nil {
		client.Close()
		logger.Error("Failed to connect to Valkey",
			zap.String("address", addr),
			zap.Error(err))
		return nil, fmt.Errorf("failed to connect to Valkey: %w", err)
	}

	logger.Debug("Valkey connection established",
		zap.String("address", addr),
		zap.Int("db", cfg.DB))

	return compat, nil
}
