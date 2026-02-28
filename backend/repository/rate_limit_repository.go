package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/valkey-io/valkey-go/valkeycompat"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/logger"
)

const (
	rateLimitLoginPrefix   = "rate_limit:login:"
	rateLimitIPPrefix      = "rate_limit:ip:"
	rateLimitLockoutPrefix = "rate_limit:lockout:"
	rateLimitViolations    = "rate_limit:violations:"
)

// RateLimitInfo contains information about the current rate limit status
type RateLimitInfo struct {
	Attempts    int
	MaxAttempts int
	IsLocked    bool
	LockoutEnds time.Time
	RetryAfter  time.Duration
	Violations  int // Number of times user has been locked out
}

type RateLimitRepositoryI interface {
	// Account-based rate limiting (by email)
	IncrementLoginAttempts(ctx context.Context, email string, window time.Duration) (int, error)
	GetLoginAttempts(ctx context.Context, email string) (int, error)
	ResetLoginAttempts(ctx context.Context, email string) error
	SetLoginLockout(ctx context.Context, email string, duration time.Duration) error
	GetLoginLockout(ctx context.Context, email string) (time.Duration, error)
	IsLoginLocked(ctx context.Context, email string) (bool, time.Duration, error)

	// IP-based rate limiting
	IncrementIPAttempts(ctx context.Context, ipAddress string, window time.Duration) (int, error)
	GetIPAttempts(ctx context.Context, ipAddress string) (int, error)
	ResetIPAttempts(ctx context.Context, ipAddress string) error
	SetIPLockout(ctx context.Context, ipAddress string, duration time.Duration) error
	GetIPLockout(ctx context.Context, ipAddress string) (time.Duration, error)
	IsIPLocked(ctx context.Context, ipAddress string) (bool, time.Duration, error)

	// Progressive lockout tracking
	IncrementViolations(ctx context.Context, key string) (int, error)
	GetViolations(ctx context.Context, key string) (int, error)
	ResetViolations(ctx context.Context, key string) error
}

type RateLimitRepository struct {
	cache valkeycompat.Cmdable
}

func NewRateLimitRepository(cache valkeycompat.Cmdable) RateLimitRepositoryI {
	return &RateLimitRepository{
		cache: cache,
	}
}

// Account-based rate limiting methods

func (r *RateLimitRepository) IncrementLoginAttempts(ctx context.Context, email string, window time.Duration) (int, error) {
	key := fmt.Sprintf("%s%s", rateLimitLoginPrefix, email)

	pipe := r.cache.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment login attempts",
			zap.String("email", email),
			zap.Error(err))
		return 0, err
	}

	count := int(incr.Val())
	logger.DebugCtx(ctx, "Login attempts incremented",
		zap.String("email", email),
		zap.Int("count", count))

	return count, nil
}

func (r *RateLimitRepository) GetLoginAttempts(ctx context.Context, email string) (int, error) {
	key := fmt.Sprintf("%s%s", rateLimitLoginPrefix, email)

	val, err := r.cache.Get(ctx, key).Result()
	if err == valkeycompat.Nil {
		return 0, nil
	}
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get login attempts",
			zap.String("email", email),
			zap.Error(err))
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RateLimitRepository) ResetLoginAttempts(ctx context.Context, email string) error {
	key := fmt.Sprintf("%s%s", rateLimitLoginPrefix, email)

	if err := r.cache.Del(ctx, key).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to reset login attempts",
			zap.String("email", email),
			zap.Error(err))
		return err
	}

	logger.DebugCtx(ctx, "Login attempts reset",
		zap.String("email", email))
	return nil
}

func (r *RateLimitRepository) SetLoginLockout(ctx context.Context, email string, duration time.Duration) error {
	key := fmt.Sprintf("%s%s", rateLimitLockoutPrefix, email)

	if err := r.cache.Set(ctx, key, "1", duration).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to set login lockout",
			zap.String("email", email),
			zap.Duration("duration", duration),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Login lockout set",
		zap.String("email", email),
		zap.Duration("duration", duration))
	return nil
}

func (r *RateLimitRepository) GetLoginLockout(ctx context.Context, email string) (time.Duration, error) {
	key := fmt.Sprintf("%s%s", rateLimitLockoutPrefix, email)

	ttl, err := r.cache.TTL(ctx, key).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get login lockout TTL",
			zap.String("email", email),
			zap.Error(err))
		return 0, err
	}

	if ttl < 0 {
		return 0, nil
	}

	return ttl, nil
}

func (r *RateLimitRepository) IsLoginLocked(ctx context.Context, email string) (bool, time.Duration, error) {
	key := fmt.Sprintf("%s%s", rateLimitLockoutPrefix, email)

	ttl, err := r.cache.TTL(ctx, key).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check login lockout",
			zap.String("email", email),
			zap.Error(err))
		return false, 0, err
	}

	// TTL returns -2 if key doesn't exist, -1 if no expiry
	if ttl < 0 {
		return false, 0, nil
	}

	return true, ttl, nil
}

// IP-based rate limiting methods

func (r *RateLimitRepository) IncrementIPAttempts(ctx context.Context, ipAddress string, window time.Duration) (int, error) {
	key := fmt.Sprintf("%s%s", rateLimitIPPrefix, ipAddress)

	pipe := r.cache.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment IP attempts",
			zap.String("ip", ipAddress),
			zap.Error(err))
		return 0, err
	}

	count := int(incr.Val())
	logger.DebugCtx(ctx, "IP attempts incremented",
		zap.String("ip", ipAddress),
		zap.Int("count", count))

	return count, nil
}

func (r *RateLimitRepository) GetIPAttempts(ctx context.Context, ipAddress string) (int, error) {
	key := fmt.Sprintf("%s%s", rateLimitIPPrefix, ipAddress)

	val, err := r.cache.Get(ctx, key).Result()
	if err == valkeycompat.Nil {
		return 0, nil
	}
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get IP attempts",
			zap.String("ip", ipAddress),
			zap.Error(err))
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RateLimitRepository) ResetIPAttempts(ctx context.Context, ipAddress string) error {
	key := fmt.Sprintf("%s%s", rateLimitIPPrefix, ipAddress)

	if err := r.cache.Del(ctx, key).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to reset IP attempts",
			zap.String("ip", ipAddress),
			zap.Error(err))
		return err
	}

	logger.DebugCtx(ctx, "IP attempts reset",
		zap.String("ip", ipAddress))
	return nil
}

func (r *RateLimitRepository) SetIPLockout(ctx context.Context, ipAddress string, duration time.Duration) error {
	key := fmt.Sprintf("%sip:%s", rateLimitLockoutPrefix, ipAddress)

	if err := r.cache.Set(ctx, key, "1", duration).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to set IP lockout",
			zap.String("ip", ipAddress),
			zap.Duration("duration", duration),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "IP lockout set",
		zap.String("ip", ipAddress),
		zap.Duration("duration", duration))
	return nil
}

func (r *RateLimitRepository) GetIPLockout(ctx context.Context, ipAddress string) (time.Duration, error) {
	key := fmt.Sprintf("%sip:%s", rateLimitLockoutPrefix, ipAddress)

	ttl, err := r.cache.TTL(ctx, key).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get IP lockout TTL",
			zap.String("ip", ipAddress),
			zap.Error(err))
		return 0, err
	}

	if ttl < 0 {
		return 0, nil
	}

	return ttl, nil
}

func (r *RateLimitRepository) IsIPLocked(ctx context.Context, ipAddress string) (bool, time.Duration, error) {
	key := fmt.Sprintf("%sip:%s", rateLimitLockoutPrefix, ipAddress)

	ttl, err := r.cache.TTL(ctx, key).Result()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check IP lockout",
			zap.String("ip", ipAddress),
			zap.Error(err))
		return false, 0, err
	}

	if ttl < 0 {
		return false, 0, nil
	}

	return true, ttl, nil
}

// Progressive lockout tracking methods

func (r *RateLimitRepository) IncrementViolations(ctx context.Context, key string) (int, error) {
	fullKey := fmt.Sprintf("%s%s", rateLimitViolations, key)

	// Violations persist for 24 hours
	pipe := r.cache.Pipeline()
	incr := pipe.Incr(ctx, fullKey)
	pipe.Expire(ctx, fullKey, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment violations",
			zap.String("key", key),
			zap.Error(err))
		return 0, err
	}

	count := int(incr.Val())
	logger.DebugCtx(ctx, "Violations incremented",
		zap.String("key", key),
		zap.Int("count", count))

	return count, nil
}

func (r *RateLimitRepository) GetViolations(ctx context.Context, key string) (int, error) {
	fullKey := fmt.Sprintf("%s%s", rateLimitViolations, key)

	val, err := r.cache.Get(ctx, fullKey).Result()
	if err == valkeycompat.Nil {
		return 0, nil
	}
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get violations",
			zap.String("key", key),
			zap.Error(err))
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *RateLimitRepository) ResetViolations(ctx context.Context, key string) error {
	fullKey := fmt.Sprintf("%s%s", rateLimitViolations, key)

	if err := r.cache.Del(ctx, fullKey).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to reset violations",
			zap.String("key", key),
			zap.Error(err))
		return err
	}

	logger.DebugCtx(ctx, "Violations reset",
		zap.String("key", key))
	return nil
}
