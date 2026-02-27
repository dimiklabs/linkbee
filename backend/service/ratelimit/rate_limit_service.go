package ratelimit

import (
	"context"
	"math"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/response"
)

type RateLimitServiceI interface {
	// CheckLoginRateLimit checks if login is allowed for the given email and IP
	// Returns nil if allowed, or a ServiceError if rate limited
	CheckLoginRateLimit(ctx context.Context, email, ipAddress string) *dto.ServiceError

	// RecordFailedLogin records a failed login attempt for rate limiting
	RecordFailedLogin(ctx context.Context, email, ipAddress string) error

	// RecordSuccessfulLogin clears rate limit counters on successful login
	RecordSuccessfulLogin(ctx context.Context, email, ipAddress string) error

	// GetRateLimitStatus returns current rate limit status for email and IP
	GetRateLimitStatus(ctx context.Context, email, ipAddress string) (*response.RateLimitStatusResponse, error)

	// CheckForgotPasswordRateLimit checks if a forgot-password request is allowed for the given email and IP
	// Returns nil if allowed, or a ServiceError if rate limited (max 5 per 15 minutes per IP)
	CheckForgotPasswordRateLimit(ctx context.Context, email, ip string) *dto.ServiceError

	// CheckSignupRateLimit checks if a signup request is allowed for the given IP
	// Returns nil if allowed, or a ServiceError if rate limited (max 10 per hour per IP)
	CheckSignupRateLimit(ctx context.Context, ip string) *dto.ServiceError
}

type RateLimitService struct {
	rateLimitRepo repository.RateLimitRepositoryI
	config        *config.RateLimitConfig
}

func NewRateLimitService(
	rateLimitRepo repository.RateLimitRepositoryI,
	config *config.RateLimitConfig,
) RateLimitServiceI {
	return &RateLimitService{
		rateLimitRepo: rateLimitRepo,
		config:        config,
	}
}

func (s *RateLimitService) CheckLoginRateLimit(ctx context.Context, email, ipAddress string) *dto.ServiceError {
	if !s.config.Enabled {
		return nil
	}

	// Check IP lockout first (more severe)
	ipLocked, ipTTL, err := s.rateLimitRepo.IsIPLocked(ctx, ipAddress)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check IP lockout", zap.Error(err))
		// Don't block login on rate limit check errors
		return nil
	}

	if ipLocked {
		logger.WarnCtx(ctx, "IP is locked out",
			zap.String("ip", ipAddress),
			zap.Duration("remaining", ipTTL))
		return dto.NewServiceErrorWithData(
			constant.ErrCodeRateLimited,
			constant.ErrMsgIPRateLimited,
			http.StatusTooManyRequests,
			map[string]interface{}{
				"retry_after_seconds": int(ipTTL.Seconds()),
				"lockout_type":        "ip",
			},
		)
	}

	// Check account lockout
	loginLocked, loginTTL, err := s.rateLimitRepo.IsLoginLocked(ctx, email)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check login lockout", zap.Error(err))
		return nil
	}

	if loginLocked {
		logger.WarnCtx(ctx, "Account is locked out",
			zap.String("email", email),
			zap.Duration("remaining", loginTTL))
		return dto.NewServiceErrorWithData(
			constant.ErrCodeRateLimited,
			constant.ErrMsgAccountRateLimited,
			http.StatusTooManyRequests,
			map[string]interface{}{
				"retry_after_seconds": int(loginTTL.Seconds()),
				"lockout_type":        "account",
			},
		)
	}

	// Check if approaching limits (warning)
	loginAttempts, _ := s.rateLimitRepo.GetLoginAttempts(ctx, email)
	ipAttempts, _ := s.rateLimitRepo.GetIPAttempts(ctx, ipAddress)

	remainingLoginAttempts := s.config.LoginMaxAttempts - loginAttempts
	remainingIPAttempts := s.config.IPMaxAttempts - ipAttempts

	if remainingLoginAttempts <= 0 || remainingIPAttempts <= 0 {
		// Should have been locked, but wasn't - this is a race condition
		// Lock now and return error
		if remainingLoginAttempts <= 0 {
			_ = s.lockAccount(ctx, email)
			return dto.NewServiceErrorWithData(
				constant.ErrCodeRateLimited,
				constant.ErrMsgAccountRateLimited,
				http.StatusTooManyRequests,
				map[string]interface{}{
					"retry_after_seconds": s.config.LoginLockoutMinutes * 60,
					"lockout_type":        "account",
				},
			)
		}
		if remainingIPAttempts <= 0 {
			_ = s.lockIP(ctx, ipAddress)
			return dto.NewServiceErrorWithData(
				constant.ErrCodeRateLimited,
				constant.ErrMsgIPRateLimited,
				http.StatusTooManyRequests,
				map[string]interface{}{
					"retry_after_seconds": s.config.IPLockoutMinutes * 60,
					"lockout_type":        "ip",
				},
			)
		}
	}

	return nil
}

func (s *RateLimitService) RecordFailedLogin(ctx context.Context, email, ipAddress string) error {
	if !s.config.Enabled {
		return nil
	}

	loginWindow := time.Duration(s.config.LoginWindowMinutes) * time.Minute
	ipWindow := time.Duration(s.config.IPWindowMinutes) * time.Minute

	// Increment login attempts for this email
	loginAttempts, err := s.rateLimitRepo.IncrementLoginAttempts(ctx, email, loginWindow)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment login attempts", zap.Error(err))
		return err
	}

	// Increment IP attempts
	ipAttempts, err := s.rateLimitRepo.IncrementIPAttempts(ctx, ipAddress, ipWindow)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment IP attempts", zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Failed login attempt recorded",
		zap.String("email", email),
		zap.String("ip", ipAddress),
		zap.Int("login_attempts", loginAttempts),
		zap.Int("ip_attempts", ipAttempts),
		zap.Int("max_login_attempts", s.config.LoginMaxAttempts),
		zap.Int("max_ip_attempts", s.config.IPMaxAttempts))

	// Check if we need to lock the account
	if loginAttempts >= s.config.LoginMaxAttempts {
		if err := s.lockAccount(ctx, email); err != nil {
			logger.ErrorCtx(ctx, "Failed to lock account", zap.Error(err))
		}
	}

	// Check if we need to lock the IP
	if ipAttempts >= s.config.IPMaxAttempts {
		if err := s.lockIP(ctx, ipAddress); err != nil {
			logger.ErrorCtx(ctx, "Failed to lock IP", zap.Error(err))
		}
	}

	return nil
}

func (s *RateLimitService) RecordSuccessfulLogin(ctx context.Context, email, ipAddress string) error {
	if !s.config.Enabled {
		return nil
	}

	// Reset login attempts for this email on successful login
	if err := s.rateLimitRepo.ResetLoginAttempts(ctx, email); err != nil {
		logger.ErrorCtx(ctx, "Failed to reset login attempts", zap.Error(err))
		// Don't fail the login on this error
	}

	// Reset violations for this email (reward good behavior)
	if err := s.rateLimitRepo.ResetViolations(ctx, email); err != nil {
		logger.ErrorCtx(ctx, "Failed to reset violations", zap.Error(err))
	}

	logger.InfoCtx(ctx, "Login rate limits cleared on successful login",
		zap.String("email", email),
		zap.String("ip", ipAddress))

	return nil
}

func (s *RateLimitService) GetRateLimitStatus(ctx context.Context, email, ipAddress string) (*response.RateLimitStatusResponse, error) {
	loginAttempts, _ := s.rateLimitRepo.GetLoginAttempts(ctx, email)
	ipAttempts, _ := s.rateLimitRepo.GetIPAttempts(ctx, ipAddress)
	loginLocked, loginTTL, _ := s.rateLimitRepo.IsLoginLocked(ctx, email)
	ipLocked, ipTTL, _ := s.rateLimitRepo.IsIPLocked(ctx, ipAddress)
	violations, _ := s.rateLimitRepo.GetViolations(ctx, email)

	return &response.RateLimitStatusResponse{
		AccountAttempts:          loginAttempts,
		AccountMaxAttempts:       s.config.LoginMaxAttempts,
		AccountRemainingAttempts: max(0, s.config.LoginMaxAttempts-loginAttempts),
		AccountLocked:            loginLocked,
		AccountLockoutRemaining:  int(loginTTL.Seconds()),
		IPAttempts:               ipAttempts,
		IPMaxAttempts:            s.config.IPMaxAttempts,
		IPRemainingAttempts:      max(0, s.config.IPMaxAttempts-ipAttempts),
		IPLocked:                 ipLocked,
		IPLockoutRemaining:       int(ipTTL.Seconds()),
		Violations:               violations,
	}, nil
}

func (s *RateLimitService) CheckForgotPasswordRateLimit(ctx context.Context, email, ip string) *dto.ServiceError {
	if !s.config.Enabled {
		return nil
	}

	// Use a prefixed compound key so forgot-password attempts are namespaced
	// separately from login IP attempts in the cache store.
	scopedIP := "forgot_pwd:" + ip
	window := 15 * time.Minute
	const maxAttempts = 5

	locked, ttl, err := s.rateLimitRepo.IsIPLocked(ctx, scopedIP)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check forgot-password rate limit", zap.Error(err))
		return nil
	}
	if locked {
		logger.WarnCtx(ctx, "Forgot-password rate limited",
			zap.String("ip", ip),
			zap.Duration("remaining", ttl))
		return dto.NewServiceErrorWithData(
			constant.ErrCodeRateLimited,
			"Too many password reset requests. Please try again later.",
			http.StatusTooManyRequests,
			map[string]interface{}{
				"retry_after_seconds": int(ttl.Seconds()),
			},
		)
	}

	attempts, err := s.rateLimitRepo.IncrementIPAttempts(ctx, scopedIP, window)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment forgot-password attempts", zap.Error(err))
		return nil
	}

	if attempts >= maxAttempts {
		_ = s.rateLimitRepo.SetIPLockout(ctx, scopedIP, window)
		return dto.NewServiceErrorWithData(
			constant.ErrCodeRateLimited,
			"Too many password reset requests. Please try again later.",
			http.StatusTooManyRequests,
			map[string]interface{}{
				"retry_after_seconds": int(window.Seconds()),
			},
		)
	}

	return nil
}

func (s *RateLimitService) CheckSignupRateLimit(ctx context.Context, ip string) *dto.ServiceError {
	if !s.config.Enabled {
		return nil
	}

	// Use a prefixed compound key so signup attempts are namespaced
	// separately from login IP attempts in the cache store.
	scopedIP := "signup:" + ip
	window := time.Hour
	const maxAttempts = 10

	locked, ttl, err := s.rateLimitRepo.IsIPLocked(ctx, scopedIP)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check signup rate limit", zap.Error(err))
		return nil
	}
	if locked {
		logger.WarnCtx(ctx, "Signup rate limited",
			zap.String("ip", ip),
			zap.Duration("remaining", ttl))
		return dto.NewServiceErrorWithData(
			constant.ErrCodeRateLimited,
			"Too many signup attempts. Please try again later.",
			http.StatusTooManyRequests,
			map[string]interface{}{
				"retry_after_seconds": int(ttl.Seconds()),
			},
		)
	}

	attempts, err := s.rateLimitRepo.IncrementIPAttempts(ctx, scopedIP, window)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to increment signup attempts", zap.Error(err))
		return nil
	}

	if attempts >= maxAttempts {
		_ = s.rateLimitRepo.SetIPLockout(ctx, scopedIP, window)
		return dto.NewServiceErrorWithData(
			constant.ErrCodeRateLimited,
			"Too many signup attempts. Please try again later.",
			http.StatusTooManyRequests,
			map[string]interface{}{
				"retry_after_seconds": int(window.Seconds()),
			},
		)
	}

	return nil
}

func (s *RateLimitService) lockAccount(ctx context.Context, email string) error {
	lockoutDuration := s.calculateLockoutDuration(ctx, email)

	// Increment violations
	_, _ = s.rateLimitRepo.IncrementViolations(ctx, email)

	// Set lockout
	if err := s.rateLimitRepo.SetLoginLockout(ctx, email, lockoutDuration); err != nil {
		return err
	}

	// Reset attempt counter (lockout handles the blocking now)
	_ = s.rateLimitRepo.ResetLoginAttempts(ctx, email)

	logger.WarnCtx(ctx, "Account locked due to too many failed login attempts",
		zap.String("email", email),
		zap.Duration("lockout_duration", lockoutDuration))

	return nil
}

func (s *RateLimitService) lockIP(ctx context.Context, ipAddress string) error {
	lockoutDuration := s.calculateIPLockoutDuration(ctx, ipAddress)

	// Increment violations for IP
	_, _ = s.rateLimitRepo.IncrementViolations(ctx, "ip:"+ipAddress)

	// Set lockout
	if err := s.rateLimitRepo.SetIPLockout(ctx, ipAddress, lockoutDuration); err != nil {
		return err
	}

	// Reset attempt counter
	_ = s.rateLimitRepo.ResetIPAttempts(ctx, ipAddress)

	logger.WarnCtx(ctx, "IP locked due to too many failed login attempts",
		zap.String("ip", ipAddress),
		zap.Duration("lockout_duration", lockoutDuration))

	return nil
}

func (s *RateLimitService) calculateLockoutDuration(ctx context.Context, email string) time.Duration {
	baseDuration := time.Duration(s.config.LoginLockoutMinutes) * time.Minute

	if !s.config.ProgressiveEnabled {
		return baseDuration
	}

	violations, err := s.rateLimitRepo.GetViolations(ctx, email)
	if err != nil {
		return baseDuration
	}

	// Calculate multiplier: base * (multiplier ^ violations)
	multiplier := math.Pow(s.config.ProgressiveMultiplier, float64(violations))
	if multiplier > s.config.ProgressiveMaxMultiplier {
		multiplier = s.config.ProgressiveMaxMultiplier
	}

	progressiveDuration := time.Duration(float64(baseDuration) * multiplier)

	logger.DebugCtx(ctx, "Calculated progressive lockout duration",
		zap.String("email", email),
		zap.Int("violations", violations),
		zap.Float64("multiplier", multiplier),
		zap.Duration("duration", progressiveDuration))

	return progressiveDuration
}

func (s *RateLimitService) calculateIPLockoutDuration(ctx context.Context, ipAddress string) time.Duration {
	baseDuration := time.Duration(s.config.IPLockoutMinutes) * time.Minute

	if !s.config.ProgressiveEnabled {
		return baseDuration
	}

	violations, err := s.rateLimitRepo.GetViolations(ctx, "ip:"+ipAddress)
	if err != nil {
		return baseDuration
	}

	multiplier := math.Pow(s.config.ProgressiveMultiplier, float64(violations))
	if multiplier > s.config.ProgressiveMaxMultiplier {
		multiplier = s.config.ProgressiveMaxMultiplier
	}

	return time.Duration(float64(baseDuration) * multiplier)
}
