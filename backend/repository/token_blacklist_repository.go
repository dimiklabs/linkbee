package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go/valkeycompat"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/logger"
)

const tokenBlacklistPrefix = "token_blacklist:"

type TokenBlacklistRepositoryI interface {
	Add(ctx context.Context, jti string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

type TokenBlacklistRepository struct {
	cache valkeycompat.Cmdable
}

func NewTokenBlacklistRepository(cache valkeycompat.Cmdable) TokenBlacklistRepositoryI {
	return &TokenBlacklistRepository{
		cache: cache,
	}
}

func (r *TokenBlacklistRepository) Add(ctx context.Context, jti string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", tokenBlacklistPrefix, jti)

	logger.DebugCtx(ctx, "Adding token to blacklist",
		zap.String("jti", jti),
		zap.Duration("expiration", expiration))

	if err := r.cache.Set(ctx, key, "1", expiration).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to add token to blacklist",
			zap.String("jti", jti),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Token added to blacklist",
		zap.String("jti", jti))
	return nil
}

func (r *TokenBlacklistRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("%s%s", tokenBlacklistPrefix, jti)

	exists, err := r.cache.Exists(ctx, key).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check token blacklist",
			zap.String("jti", jti),
			zap.Error(err))
		return false, err
	}

	return exists > 0, nil
}
