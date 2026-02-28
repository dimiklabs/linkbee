package auth

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/request"
	"github.com/shafikshaon/shortlink/response"
	userSrv "github.com/shafikshaon/shortlink/service/user"
	"github.com/shafikshaon/shortlink/util"
)

type AuthServiceI interface {
	Signup(ctx context.Context, req *request.SignupRequest) (*model.User, *dto.ServiceError)
	Login(ctx context.Context, req *request.LoginRequest, userAgent, ipAddress string) (*response.LoginResponse, *response.SessionsListResponse, *dto.ServiceError)
	RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.LoginResponse, *dto.ServiceError)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *request.ChangePasswordRequest) *dto.ServiceError
	ForgotPassword(ctx context.Context, req *request.ForgotPasswordRequest) *dto.ServiceError
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) *dto.ServiceError
	GetProfile(ctx context.Context, userID uuid.UUID) (*response.ProfileResponse, *dto.ServiceError)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *request.UpdateProfileRequest) (*response.ProfileResponse, *dto.ServiceError)
	Logout(ctx context.Context, accessToken string, refreshToken string) *dto.ServiceError
	DeleteAccount(ctx context.Context, userID uuid.UUID) *dto.ServiceError
	ReactivateAccount(ctx context.Context, req *request.ReactivateAccountRequest, userAgent, ipAddress string) (*response.LoginResponse, *response.SessionsListResponse, *dto.ServiceError)
	GetSessions(ctx context.Context, userID uuid.UUID, currentJTI string) (*response.SessionsListResponse, *dto.ServiceError)
	DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) *dto.ServiceError
	LogoutAll(ctx context.Context, userID uuid.UUID) *dto.ServiceError
	VerifyEmail(ctx context.Context, token string) *dto.ServiceError
	ResendVerificationEmail(ctx context.Context, email string) *dto.ServiceError
	ValidateSession(ctx context.Context, req *request.RefreshTokenRequest) *dto.ServiceError
	GetJWTSecret() string
	GetJWTIssuer() string
	// TOTP
	GetTOTPStatus(ctx context.Context, userID uuid.UUID) (*response.TOTPStatusResponse, *dto.ServiceError)
	SetupTOTP(ctx context.Context, userID uuid.UUID) (*response.TOTPSetupResponse, *dto.ServiceError)
	ConfirmTOTP(ctx context.Context, userID uuid.UUID, code string) (*response.TOTPBackupCodesResponse, *dto.ServiceError)
	DisableTOTP(ctx context.Context, userID uuid.UUID, password string) *dto.ServiceError
	VerifyTOTPLogin(ctx context.Context, req *request.TOTPVerifyLoginRequest, userAgent, ipAddress string) (*response.LoginResponse, *dto.ServiceError)
}

type EmailSenderI interface {
	SendVerificationEmail(ctx context.Context, user *model.User) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, resetToken string) error
}

type AuthService struct {
	userService         userSrv.UserServiceI
	passwordResetRepo   repository.PasswordResetRepositoryI
	tokenBlacklistRepo  repository.TokenBlacklistRepositoryI
	sessionRepo         repository.SessionRepositoryI
	totpBackupCodeRepo  repository.TotpBackupCodeRepositoryI
	emailService        EmailSenderI
	cfg                 *config.AppConfig
	sessionCfg          *config.SessionConfig
}

func NewAuthService(
	userService userSrv.UserServiceI,
	passwordResetRepo repository.PasswordResetRepositoryI,
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI,
	sessionRepo repository.SessionRepositoryI,
	totpBackupCodeRepo repository.TotpBackupCodeRepositoryI,
	emailService EmailSenderI,
	cfg *config.AppConfig,
	sessionCfg *config.SessionConfig,
) AuthServiceI {
	return &AuthService{
		userService:        userService,
		passwordResetRepo:  passwordResetRepo,
		tokenBlacklistRepo: tokenBlacklistRepo,
		sessionRepo:        sessionRepo,
		totpBackupCodeRepo: totpBackupCodeRepo,
		emailService:       emailService,
		cfg:                cfg,
		sessionCfg:         sessionCfg,
	}
}

func (s *AuthService) GetJWTSecret() string {
	return s.cfg.JWTSecret
}

func (s *AuthService) GetJWTIssuer() string {
	return s.cfg.JWTIssuer
}

func (s *AuthService) Signup(ctx context.Context, req *request.SignupRequest) (*model.User, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing signup",
		zap.String("email", req.Email))

	createReq := &dto.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, svcErr := s.userService.CreateUser(ctx, createReq)
	if svcErr != nil {
		return nil, svcErr
	}

	// Send verification email
	if err := s.emailService.SendVerificationEmail(ctx, user); err != nil {
		logger.ErrorCtx(ctx, "Failed to send verification email",
			zap.String("user_id", user.ID.String()),
			zap.String("email", user.Email),
			zap.Error(err))
		// Don't fail signup if email fails, user can request resend
	}

	logger.InfoCtx(ctx, "Signup completed successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *request.LoginRequest, userAgent, ipAddress string) (*response.LoginResponse, *response.SessionsListResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing login",
		zap.String("email", req.Email),
		zap.Bool("remember_me", req.RememberMe))

	user, svcErr := s.userService.GetUserByEmail(ctx, req.Email)
	if svcErr != nil {
		logger.WarnCtx(ctx, "Login failed: user not found",
			zap.String("email", req.Email))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInvalidCredentials, constant.ErrMsgInvalidCredentials, http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.WarnCtx(ctx, "Login failed: invalid password",
			zap.String("email", req.Email))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInvalidCredentials, constant.ErrMsgInvalidCredentials, http.StatusUnauthorized)
	}

	// Check if user has scheduled deletion (in grace period)
	if user.Status == "inactive" && user.ScheduledDeletionAt != nil {
		// User is in grace period, require confirmation to reactivate
		logger.InfoCtx(ctx, "User login attempted during deletion grace period, confirmation required",
			zap.String("user_id", user.ID.String()))
		return nil, nil, dto.NewServiceError(constant.ErrCodeAccountPendingDeletion, constant.ErrMsgAccountPendingDeletion, http.StatusForbidden)
	} else if user.Status == "inactive" {
		logger.WarnCtx(ctx, "Login failed: user inactive",
			zap.String("user_id", user.ID.String()))
		return nil, nil, dto.NewServiceError(constant.ErrCodeUserInactive, constant.ErrMsgUserInactive, http.StatusForbidden)
	}

	// If TOTP is enabled, issue a short-lived TOTP session token instead of full login
	if user.TotpEnabled {
		jwtCfg := &util.JWTConfig{
			Secret: s.cfg.JWTSecret,
			Issuer: s.cfg.JWTIssuer,
		}
		totpToken, err := util.GenerateTOTPSessionToken(jwtCfg, user.ID.String())
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to generate TOTP session token",
				zap.String("user_id", user.ID.String()),
				zap.Error(err))
			return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}
		logger.InfoCtx(ctx, "TOTP required, returning TOTP session token",
			zap.String("user_id", user.ID.String()))
		return &response.LoginResponse{
			RequiresTOTP: true,
			TOTPSession:  totpToken,
		}, nil, nil
	}

	opts := &dto.LoginOptions{
		RememberMe:  req.RememberMe,
		LoginMethod: model.LoginMethodLocal,
		UserAgent:   userAgent,
		IPAddress:   ipAddress,
	}
	return s.completeLoginWithOptions(ctx, user, opts)
}

func (s *AuthService) completeLoginWithOptions(ctx context.Context, user *model.User, opts *dto.LoginOptions) (*response.LoginResponse, *response.SessionsListResponse, *dto.ServiceError) {
	// Check session count before creating tokens
	sessionCount, err := s.sessionRepo.CountActiveByUserID(ctx, user.ID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to count active sessions",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	maxSessions := s.cfg.MaxSessions
	if maxSessions < 1 {
		maxSessions = 1
	}

	if sessionCount >= int64(maxSessions) {
		logger.InfoCtx(ctx, "Max sessions reached, auto-logout oldest session",
			zap.String("user_id", user.ID.String()),
			zap.Int64("current", sessionCount),
			zap.Int("max", maxSessions))

		// Get oldest session to logout
		oldestSession, err := s.sessionRepo.GetOldestActiveByUserID(ctx, user.ID)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to get oldest active session",
				zap.String("user_id", user.ID.String()),
				zap.Error(err))
			return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}

		// Blacklist the refresh token of the oldest session
		remainingTime := time.Until(oldestSession.ExpiresAt)
		if remainingTime > 0 {
			if err := s.tokenBlacklistRepo.Add(ctx, oldestSession.RefreshTokenJTI, remainingTime); err != nil {
				logger.WarnCtx(ctx, "Failed to blacklist refresh token for auto-logout session",
					zap.String("jti", oldestSession.RefreshTokenJTI),
					zap.Error(err))
			}
		}

		// Delete the oldest session
		if err := s.sessionRepo.Delete(ctx, oldestSession.ID); err != nil {
			logger.ErrorCtx(ctx, "Failed to delete oldest session",
				zap.String("session_id", oldestSession.ID.String()),
				zap.Error(err))
			return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}

		logger.InfoCtx(ctx, "Auto-logged out oldest session",
			zap.String("session_id", oldestSession.ID.String()),
			zap.String("user_id", user.ID.String()))
	}

	// Determine refresh token expiry based on remember_me option
	refreshExpiryDays := s.cfg.JWTRefreshExpiry
	if opts.RememberMe && s.sessionCfg != nil && s.sessionCfg.RememberMeEnabled {
		refreshExpiryDays = s.sessionCfg.RememberMeRefreshExpiryDays
		logger.InfoCtx(ctx, "Using extended session expiry for remember_me",
			zap.String("user_id", user.ID.String()),
			zap.Int("expiry_days", refreshExpiryDays))
	}

	jwtCfg := &util.JWTConfig{
		Secret:              s.cfg.JWTSecret,
		Issuer:              s.cfg.JWTIssuer,
		AccessExpiryMinutes: s.cfg.JWTAccessExpiry,
		RefreshExpiryDays:   refreshExpiryDays,
	}

	tokenPair, err := util.GenerateTokenPair(jwtCfg, user.ID.String(), user.Email, user.Role)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate token pair",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Extract refresh token JTI for session
	refreshClaims, err := util.ValidateRefreshToken(tokenPair.RefreshToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to validate newly created refresh token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Parse user agent for device info
	deviceInfo := util.ParseUserAgent(opts.UserAgent)

	// Determine login method
	loginMethod := opts.LoginMethod
	if loginMethod == "" {
		loginMethod = model.LoginMethodLocal
	}

	// Create session with enhanced fields
	session := &model.Session{
		UserID:          user.ID,
		RefreshTokenJTI: refreshClaims.ID,
		UserAgent:       opts.UserAgent,
		IPAddress:       opts.IPAddress,
		CreatedAt:       time.Now(),
		ExpiresAt:       refreshClaims.ExpiresAt.Time,
		LastActivityAt:  time.Now(),
		RememberMe:      opts.RememberMe,
		DeviceName:      deviceInfo.DeviceName,
		DeviceType:      deviceInfo.DeviceType,
		Browser:         deviceInfo.Browser,
		OS:              deviceInfo.OS,
		LoginMethod:     loginMethod,
		LastActivityIP:  opts.IPAddress,
		ActivityCount:   0,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.ErrorCtx(ctx, "Failed to create session",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	if updateErr := s.userService.UpdateLastLogin(ctx, user.ID); updateErr != nil {
		logger.WarnCtx(ctx, "Failed to update last login",
			zap.String("user_id", user.ID.String()))
	}

	logger.InfoCtx(ctx, "Login successful",
		zap.String("user_id", user.ID.String()),
		zap.String("session_id", session.ID.String()),
		zap.String("email", user.Email),
		zap.Bool("remember_me", opts.RememberMe),
		zap.String("login_method", loginMethod),
		zap.String("device_type", deviceInfo.DeviceType))

	return &response.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User: response.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *request.RefreshTokenRequest) (*response.LoginResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing token refresh")

	claims, err := util.ValidateRefreshToken(req.RefreshToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid refresh token",
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, constant.ErrMsgInvalidRefreshToken, http.StatusUnauthorized)
	}

	// Check if refresh token is blacklisted
	isBlacklisted, err := s.tokenBlacklistRepo.IsBlacklisted(ctx, claims.ID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check token blacklist",
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	if isBlacklisted {
		logger.WarnCtx(ctx, "Refresh token is blacklisted",
			zap.String("jti", claims.ID))
		return nil, dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, constant.ErrMsgInvalidRefreshToken, http.StatusUnauthorized)
	}

	// Find session by refresh token JTI
	session, err := s.sessionRepo.GetByRefreshTokenJTI(ctx, claims.ID)
	if err != nil {
		logger.WarnCtx(ctx, "Session not found for refresh token",
			zap.String("jti", claims.ID))
		return nil, dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, constant.ErrMsgInvalidRefreshToken, http.StatusUnauthorized)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in refresh token",
			zap.String("subject", claims.Subject),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, constant.ErrMsgInvalidRefreshToken, http.StatusUnauthorized)
	}

	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		logger.WarnCtx(ctx, "User not found for refresh token",
			zap.String("user_id", userID.String()))
		return nil, dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, constant.ErrMsgInvalidRefreshToken, http.StatusUnauthorized)
	}

	if user.Status == "inactive" {
		logger.WarnCtx(ctx, "Token refresh failed: user inactive",
			zap.String("user_id", user.ID.String()))
		return nil, dto.NewServiceError(constant.ErrCodeUserInactive, constant.ErrMsgUserInactive, http.StatusForbidden)
	}

	jwtCfg := &util.JWTConfig{
		Secret:              s.cfg.JWTSecret,
		Issuer:              s.cfg.JWTIssuer,
		AccessExpiryMinutes: s.cfg.JWTAccessExpiry,
		RefreshExpiryDays:   s.cfg.JWTRefreshExpiry,
	}

	tokenPair, err := util.GenerateTokenPair(jwtCfg, user.ID.String(), user.Email, user.Role)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate token pair",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Extract the new refresh token's JTI for session update
	newRefreshClaims, err := util.ValidateRefreshToken(tokenPair.RefreshToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to validate newly created refresh token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Blacklist the old refresh token so it can't be reused
	oldTokenRemaining := time.Until(session.ExpiresAt)
	if oldTokenRemaining > 0 {
		if err := s.tokenBlacklistRepo.Add(ctx, claims.ID, oldTokenRemaining); err != nil {
			logger.WarnCtx(ctx, "Failed to blacklist old refresh token",
				zap.String("jti", claims.ID),
				zap.Error(err))
		}
	}

	// Update the session's refresh token JTI and expiry to match the new token
	if err := s.sessionRepo.UpdateRefreshTokenJTI(ctx, session.ID, newRefreshClaims.ID, newRefreshClaims.ExpiresAt.Time); err != nil {
		logger.ErrorCtx(ctx, "Failed to update session refresh token JTI",
			zap.String("session_id", session.ID.String()),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Update session's last activity with details
	if s.sessionCfg != nil && s.sessionCfg.TrackActivity {
		if err := s.sessionRepo.UpdateActivityWithDetails(ctx, session.ID, ""); err != nil {
			logger.WarnCtx(ctx, "Failed to update session activity",
				zap.String("session_id", session.ID.String()))
		}
	}

	logger.InfoCtx(ctx, "Token refresh successful",
		zap.String("user_id", user.ID.String()),
		zap.String("session_id", session.ID.String()))

	return &response.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User: response.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *request.ChangePasswordRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing password change",
		zap.String("user_id", userID.String()))

	changeReq := &dto.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	if svcErr := s.userService.ChangePassword(ctx, userID, changeReq); svcErr != nil {
		return svcErr
	}

	logger.InfoCtx(ctx, "Password changed successfully",
		zap.String("user_id", userID.String()))

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req *request.ForgotPasswordRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing forgot password",
		zap.String("email", req.Email))

	user, svcErr := s.userService.GetUserByEmail(ctx, req.Email)
	if svcErr != nil {
		// Return success even if user not found to prevent email enumeration
		logger.DebugCtx(ctx, "Forgot password: user not found, returning success to prevent enumeration",
			zap.String("email", req.Email))
		return nil
	}

	if user.Status == "inactive" {
		// Return success even if user inactive to prevent enumeration
		logger.DebugCtx(ctx, "Forgot password: user inactive, returning success to prevent enumeration",
			zap.String("email", req.Email))
		return nil
	}

	// Delete any existing password reset tokens for this user
	if err := s.passwordResetRepo.DeleteByUserID(ctx, user.ID); err != nil {
		logger.WarnCtx(ctx, "Failed to delete existing password reset tokens",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
	}

	// Generate secure reset token
	token, err := util.GenerateSecureToken(32)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate password reset token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Token expires in 1 hour
	passwordReset := &model.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := s.passwordResetRepo.Create(ctx, passwordReset); err != nil {
		logger.ErrorCtx(ctx, "Failed to create password reset token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	if err := s.emailService.SendPasswordResetEmail(ctx, user.Email, token); err != nil {
		logger.WarnCtx(ctx, "Failed to send password reset email — token saved but email not delivered",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		// Do not surface the email error to the caller: the token is persisted and
		// the response must not reveal whether the address exists (enumeration protection).
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing password reset")

	passwordReset, err := s.passwordResetRepo.GetByToken(ctx, req.Token)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid or expired password reset token")
		return dto.NewServiceError(constant.ErrCodeInvalidResetToken, constant.ErrMsgInvalidResetToken, http.StatusBadRequest)
	}

	// Verify user exists and is active
	user, svcErr := s.userService.GetUserByID(ctx, passwordReset.UserID)
	if svcErr != nil {
		logger.WarnCtx(ctx, "User not found for password reset",
			zap.String("user_id", passwordReset.UserID.String()))
		return dto.NewServiceError(constant.ErrCodeInvalidResetToken, constant.ErrMsgInvalidResetToken, http.StatusBadRequest)
	}

	if user.Status == "inactive" {
		logger.WarnCtx(ctx, "Password reset failed: user inactive",
			zap.String("user_id", user.ID.String()))
		return dto.NewServiceError(constant.ErrCodeUserInactive, constant.ErrMsgUserInactive, http.StatusForbidden)
	}

	// Set new password
	if svcErr := s.userService.SetPassword(ctx, user.ID, req.NewPassword); svcErr != nil {
		return svcErr
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(ctx, passwordReset.ID); err != nil {
		logger.WarnCtx(ctx, "Failed to mark password reset token as used",
			zap.String("id", passwordReset.ID.String()),
			zap.Error(err))
	}

	// Delete all password reset tokens for this user
	if err := s.passwordResetRepo.DeleteByUserID(ctx, user.ID); err != nil {
		logger.WarnCtx(ctx, "Failed to delete password reset tokens after reset",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
	}

	logger.InfoCtx(ctx, "Password reset successful",
		zap.String("user_id", user.ID.String()))

	return nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*response.ProfileResponse, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Getting user profile",
		zap.String("user_id", userID.String()))

	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		return nil, svcErr
	}

	logger.DebugCtx(ctx, "Profile retrieved successfully",
		zap.String("user_id", userID.String()))

	return &response.ProfileResponse{
		ID:                   user.ID.String(),
		Email:                user.Email,
		Role:                 user.Role,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Phone:                user.Phone,
		ProfilePicture:       user.ProfilePicture,
		ProfilePictureSource: user.ProfilePictureSource,
	}, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *request.UpdateProfileRequest) (*response.ProfileResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Updating user profile",
		zap.String("user_id", userID.String()))

	updateReq := &dto.UpdateUserRequest{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		ProfilePicture: req.ProfilePicture,
	}

	user, svcErr := s.userService.UpdateUser(ctx, userID, updateReq)
	if svcErr != nil {
		return nil, svcErr
	}

	logger.InfoCtx(ctx, "Profile updated successfully",
		zap.String("user_id", userID.String()))

	return &response.ProfileResponse{
		ID:                   user.ID.String(),
		Email:                user.Email,
		Role:                 user.Role,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Phone:                user.Phone,
		ProfilePicture:       user.ProfilePicture,
		ProfilePictureSource: user.ProfilePictureSource,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string, refreshToken string) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing logout")

	// Validate and parse the access token to get the JTI and expiration
	claims, err := util.ValidateAccessToken(accessToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid access token for logout",
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized, http.StatusUnauthorized)
	}

	// Calculate remaining time until token expiration
	expiresAt := claims.ExpiresAt.Time
	remainingTime := time.Until(expiresAt)
	if remainingTime > 0 {
		// Add token JTI to blacklist with remaining expiration time
		if err := s.tokenBlacklistRepo.Add(ctx, claims.ID, remainingTime); err != nil {
			logger.ErrorCtx(ctx, "Failed to blacklist access token",
				zap.String("jti", claims.ID),
				zap.Error(err))
			return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}
	}

	// If refresh token provided, delete session and blacklist refresh token
	if refreshToken != "" {
		refreshClaims, err := util.ValidateRefreshToken(refreshToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
		if err == nil {
			// Delete session by refresh token JTI
			if err := s.sessionRepo.DeleteByRefreshTokenJTI(ctx, refreshClaims.ID); err != nil {
				logger.WarnCtx(ctx, "Failed to delete session on logout",
					zap.Error(err))
			}

			// Blacklist the refresh token
			refreshExpiresAt := refreshClaims.ExpiresAt.Time
			refreshRemainingTime := time.Until(refreshExpiresAt)
			if refreshRemainingTime > 0 {
				if err := s.tokenBlacklistRepo.Add(ctx, refreshClaims.ID, refreshRemainingTime); err != nil {
					logger.WarnCtx(ctx, "Failed to blacklist refresh token",
						zap.String("jti", refreshClaims.ID),
						zap.Error(err))
				}
			}
		}
	}

	logger.InfoCtx(ctx, "Logout successful",
		zap.String("user_id", claims.UserID),
		zap.String("jti", claims.ID))

	return nil
}

func (s *AuthService) DeleteAccount(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing account deletion request",
		zap.String("user_id", userID.String()))

	// Schedule deletion for 30 days from now
	if svcErr := s.userService.ScheduleDeletion(ctx, userID, 30); svcErr != nil {
		return svcErr
	}

	logger.InfoCtx(ctx, "Account scheduled for deletion",
		zap.String("user_id", userID.String()),
		zap.Int("grace_period_days", 30))

	return nil
}

func (s *AuthService) ReactivateAccount(ctx context.Context, req *request.ReactivateAccountRequest, userAgent, ipAddress string) (*response.LoginResponse, *response.SessionsListResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing account reactivation",
		zap.String("email", req.Email))

	user, svcErr := s.userService.GetUserByEmail(ctx, req.Email)
	if svcErr != nil {
		logger.WarnCtx(ctx, "Reactivation failed: user not found",
			zap.String("email", req.Email))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInvalidCredentials, constant.ErrMsgInvalidCredentials, http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.WarnCtx(ctx, "Reactivation failed: invalid password",
			zap.String("email", req.Email))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInvalidCredentials, constant.ErrMsgInvalidCredentials, http.StatusUnauthorized)
	}

	// Check if user is in grace period
	if user.Status != "inactive" || user.ScheduledDeletionAt == nil {
		logger.WarnCtx(ctx, "Reactivation failed: account not pending deletion",
			zap.String("user_id", user.ID.String()))
		return nil, nil, dto.NewServiceError(constant.ErrCodeBadRequest, "Account is not pending deletion", http.StatusBadRequest)
	}

	// Reactivate the account
	if svcErr := s.userService.CancelScheduledDeletion(ctx, user.ID); svcErr != nil {
		logger.ErrorCtx(ctx, "Failed to reactivate user account",
			zap.String("user_id", user.ID.String()))
		return nil, nil, svcErr
	}

	logger.InfoCtx(ctx, "User account reactivated",
		zap.String("user_id", user.ID.String()))

	// Check session count before creating tokens
	sessionCount, err := s.sessionRepo.CountActiveByUserID(ctx, user.ID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to count active sessions",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	maxSessions := s.cfg.MaxSessions
	if maxSessions < 1 {
		maxSessions = 1
	}

	if sessionCount >= int64(maxSessions) {
		logger.InfoCtx(ctx, "Max sessions reached during reactivation, auto-logout oldest session",
			zap.String("user_id", user.ID.String()),
			zap.Int64("current", sessionCount),
			zap.Int("max", maxSessions))

		oldestSession, err := s.sessionRepo.GetOldestActiveByUserID(ctx, user.ID)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to get oldest active session",
				zap.String("user_id", user.ID.String()),
				zap.Error(err))
			return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}

		remainingTime := time.Until(oldestSession.ExpiresAt)
		if remainingTime > 0 {
			if err := s.tokenBlacklistRepo.Add(ctx, oldestSession.RefreshTokenJTI, remainingTime); err != nil {
				logger.WarnCtx(ctx, "Failed to blacklist refresh token for auto-logout session",
					zap.String("jti", oldestSession.RefreshTokenJTI),
					zap.Error(err))
			}
		}

		if err := s.sessionRepo.Delete(ctx, oldestSession.ID); err != nil {
			logger.ErrorCtx(ctx, "Failed to delete oldest session",
				zap.String("session_id", oldestSession.ID.String()),
				zap.Error(err))
			return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}

		logger.InfoCtx(ctx, "Auto-logged out oldest session during reactivation",
			zap.String("session_id", oldestSession.ID.String()),
			zap.String("user_id", user.ID.String()))
	}

	// Generate tokens and login
	jwtCfg := &util.JWTConfig{
		Secret:              s.cfg.JWTSecret,
		Issuer:              s.cfg.JWTIssuer,
		AccessExpiryMinutes: s.cfg.JWTAccessExpiry,
		RefreshExpiryDays:   s.cfg.JWTRefreshExpiry,
	}

	tokenPair, err := util.GenerateTokenPair(jwtCfg, user.ID.String(), user.Email, user.Role)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate token pair",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	refreshClaims, err := util.ValidateRefreshToken(tokenPair.RefreshToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to validate newly created refresh token",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	session := &model.Session{
		UserID:          user.ID,
		RefreshTokenJTI: refreshClaims.ID,
		UserAgent:       userAgent,
		IPAddress:       ipAddress,
		CreatedAt:       time.Now(),
		ExpiresAt:       refreshClaims.ExpiresAt.Time,
		LastActivityAt:  time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.ErrorCtx(ctx, "Failed to create session",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return nil, nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	if updateErr := s.userService.UpdateLastLogin(ctx, user.ID); updateErr != nil {
		logger.WarnCtx(ctx, "Failed to update last login",
			zap.String("user_id", user.ID.String()))
	}

	logger.InfoCtx(ctx, "Account reactivated and login successful",
		zap.String("user_id", user.ID.String()),
		zap.String("session_id", session.ID.String()),
		zap.String("email", user.Email))

	return &response.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    tokenPair.ExpiresIn,
		User: response.UserResponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil, nil
}

func (s *AuthService) GetSessions(ctx context.Context, userID uuid.UUID, currentJTI string) (*response.SessionsListResponse, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Getting active sessions",
		zap.String("user_id", userID.String()))

	sessions, err := s.sessionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get active sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	maxSessions := s.cfg.MaxSessions
	if maxSessions < 1 {
		maxSessions = 1
	}

	sessionResponses := make([]response.SessionResponse, len(sessions))
	for i, sess := range sessions {
		sessionResponses[i] = response.SessionResponse{
			ID:             sess.ID.String(),
			UserAgent:      sess.UserAgent,
			IPAddress:      sess.IPAddress,
			CreatedAt:      sess.CreatedAt,
			LastActivityAt: sess.LastActivityAt,
			IsCurrent:      sess.RefreshTokenJTI == currentJTI,
			RememberMe:     sess.RememberMe,
			DeviceName:     sess.DeviceName,
			DeviceType:     sess.DeviceType,
			Browser:        sess.Browser,
			OS:             sess.OS,
			Location:       sess.Location,
			LoginMethod:    sess.LoginMethod,
			LastActivityIP: sess.LastActivityIP,
			ActivityCount:  sess.ActivityCount,
		}
	}

	return &response.SessionsListResponse{
		Sessions:   sessionResponses,
		Count:      len(sessionResponses),
		MaxAllowed: maxSessions,
	}, nil
}

func (s *AuthService) DeleteSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Deleting session",
		zap.String("user_id", userID.String()),
		zap.String("session_id", sessionID.String()))

	// Get session to verify ownership and get JTI for blacklisting
	session, err := s.sessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		logger.WarnCtx(ctx, "Session not found",
			zap.String("session_id", sessionID.String()))
		return dto.NewServiceError(constant.ErrCodeSessionNotFound, constant.ErrMsgSessionNotFound, http.StatusNotFound)
	}

	// Verify session belongs to user
	if session.UserID != userID {
		logger.WarnCtx(ctx, "Session does not belong to user",
			zap.String("session_id", sessionID.String()),
			zap.String("user_id", userID.String()))
		return dto.NewServiceError(constant.ErrCodeForbidden, constant.ErrMsgForbidden, http.StatusForbidden)
	}

	// Blacklist the refresh token
	remainingTime := time.Until(session.ExpiresAt)
	if remainingTime > 0 {
		if err := s.tokenBlacklistRepo.Add(ctx, session.RefreshTokenJTI, remainingTime); err != nil {
			logger.WarnCtx(ctx, "Failed to blacklist refresh token for deleted session",
				zap.String("jti", session.RefreshTokenJTI),
				zap.Error(err))
		}
	}

	// Delete session
	if err := s.sessionRepo.Delete(ctx, sessionID); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete session",
			zap.String("session_id", sessionID.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "Session deleted successfully",
		zap.String("session_id", sessionID.String()))

	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Logging out from all devices",
		zap.String("user_id", userID.String()))

	// Get all active sessions to blacklist their refresh tokens
	sessions, err := s.sessionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get active sessions for logout all",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Blacklist all refresh tokens
	for _, session := range sessions {
		remainingTime := time.Until(session.ExpiresAt)
		if remainingTime > 0 {
			if err := s.tokenBlacklistRepo.Add(ctx, session.RefreshTokenJTI, remainingTime); err != nil {
				logger.WarnCtx(ctx, "Failed to blacklist refresh token during logout all",
					zap.String("jti", session.RefreshTokenJTI),
					zap.Error(err))
			}
		}
	}

	// Delete all sessions for the user
	if err := s.sessionRepo.DeleteByUserID(ctx, userID); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete all sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "Logged out from all devices successfully",
		zap.String("user_id", userID.String()),
		zap.Int("sessions_invalidated", len(sessions)))

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing email verification")

	if err := s.emailService.VerifyEmail(ctx, token); err != nil {
		logger.WarnCtx(ctx, "Email verification failed",
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInvalidVerificationToken, constant.ErrMsgInvalidVerificationToken, http.StatusBadRequest)
	}

	logger.InfoCtx(ctx, "Email verified successfully")
	return nil
}

func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) *dto.ServiceError {
	logger.InfoCtx(ctx, "Processing resend verification email",
		zap.String("email", email))

	user, svcErr := s.userService.GetUserByEmail(ctx, email)
	if svcErr != nil {
		// Return success even if user not found to prevent email enumeration
		logger.DebugCtx(ctx, "Resend verification: user not found, returning success to prevent enumeration",
			zap.String("email", email))
		return nil
	}

	if user.EmailVerified {
		// Return success even if already verified to prevent enumeration
		logger.DebugCtx(ctx, "Resend verification: email already verified, returning success to prevent enumeration",
			zap.String("email", email))
		return nil
	}

	if err := s.emailService.ResendVerificationEmail(ctx, user.ID, user.Email); err != nil {
		logger.ErrorCtx(ctx, "Failed to resend verification email",
			zap.String("email", email),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeEmailSendFailed, constant.ErrMsgEmailSendFailed, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "Verification email resent successfully",
		zap.String("email", email))
	return nil
}

func (s *AuthService) GetTOTPStatus(ctx context.Context, userID uuid.UUID) (*response.TOTPStatusResponse, *dto.ServiceError) {
	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		return nil, svcErr
	}
	return &response.TOTPStatusResponse{Enabled: user.TotpEnabled}, nil
}

func (s *AuthService) SetupTOTP(ctx context.Context, userID uuid.UUID) (*response.TOTPSetupResponse, *dto.ServiceError) {
	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		return nil, svcErr
	}

	if user.TotpEnabled {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPAlreadyEnabled, constant.ErrMsgTOTPAlreadyEnabled, http.StatusConflict)
	}

	// Generate a random 20-byte secret, base32-encode it (RFC 4648)
	secretBytes := make([]byte, 20)
	if _, err := rand.Read(secretBytes); err != nil {
		logger.ErrorCtx(ctx, "Failed to generate TOTP secret", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	secret := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secretBytes)

	// Temporarily store the pending secret (not yet enabled)
	if err := s.userService.UpdateFields(ctx, userID, map[string]interface{}{"totp_secret": secret}); err != nil {
		logger.ErrorCtx(ctx, "Failed to store TOTP secret", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Build otpauth:// URI for the QR code
	qrURL := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		s.cfg.JWTIssuer, user.Email, secret, s.cfg.JWTIssuer)

	logger.InfoCtx(ctx, "TOTP setup initiated", zap.String("user_id", userID.String()))
	return &response.TOTPSetupResponse{Secret: secret, QRCodeURL: qrURL}, nil
}

func (s *AuthService) ConfirmTOTP(ctx context.Context, userID uuid.UUID, code string) (*response.TOTPBackupCodesResponse, *dto.ServiceError) {
	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		return nil, svcErr
	}

	if user.TotpEnabled {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPAlreadyEnabled, constant.ErrMsgTOTPAlreadyEnabled, http.StatusConflict)
	}

	if user.TotpSecret == "" {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPNotEnabled, "TOTP setup has not been initiated", http.StatusBadRequest)
	}

	// Validate the TOTP code against the stored (pending) secret
	valid := totp.Validate(code, user.TotpSecret)
	if !valid {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPInvalidCode, constant.ErrMsgTOTPInvalidCode, http.StatusUnauthorized)
	}

	// Enable TOTP
	if err := s.userService.UpdateFields(ctx, userID, map[string]interface{}{"totp_enabled": true}); err != nil {
		logger.ErrorCtx(ctx, "Failed to enable TOTP", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Generate 10 single-use backup codes
	backupPlain, backupHashed, genErr := generateBackupCodes(10)
	if genErr != nil {
		logger.ErrorCtx(ctx, "Failed to generate backup codes", zap.Error(genErr))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Delete any existing backup codes and store fresh ones
	_ = s.totpBackupCodeRepo.DeleteByUserID(ctx, userID)
	models := make([]*model.TotpBackupCode, len(backupHashed))
	for i, h := range backupHashed {
		models[i] = &model.TotpBackupCode{UserID: userID, CodeHash: h}
	}
	if err := s.totpBackupCodeRepo.CreateBatch(ctx, models); err != nil {
		logger.ErrorCtx(ctx, "Failed to store backup codes", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "TOTP confirmed and enabled", zap.String("user_id", userID.String()))
	return &response.TOTPBackupCodesResponse{BackupCodes: backupPlain}, nil
}

func (s *AuthService) DisableTOTP(ctx context.Context, userID uuid.UUID, password string) *dto.ServiceError {
	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		return svcErr
	}

	if !user.TotpEnabled {
		return dto.NewServiceError(constant.ErrCodeTOTPNotEnabled, constant.ErrMsgTOTPNotEnabled, http.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return dto.NewServiceError(constant.ErrCodeInvalidCredentials, constant.ErrMsgInvalidCredentials, http.StatusUnauthorized)
	}

	if err := s.userService.UpdateFields(ctx, userID, map[string]interface{}{
		"totp_enabled": false,
		"totp_secret":  "",
	}); err != nil {
		logger.ErrorCtx(ctx, "Failed to disable TOTP", zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	_ = s.totpBackupCodeRepo.DeleteByUserID(ctx, userID)

	logger.InfoCtx(ctx, "TOTP disabled", zap.String("user_id", userID.String()))
	return nil
}

func (s *AuthService) VerifyTOTPLogin(ctx context.Context, req *request.TOTPVerifyLoginRequest, userAgent, ipAddress string) (*response.LoginResponse, *dto.ServiceError) {
	jwtCfg := &util.JWTConfig{
		Secret: s.cfg.JWTSecret,
		Issuer: s.cfg.JWTIssuer,
	}

	claims, err := util.ValidateTOTPSessionToken(req.TOTPSession, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid TOTP session token", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeTOTPInvalidSession, constant.ErrMsgTOTPInvalidSession, http.StatusUnauthorized)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPInvalidSession, constant.ErrMsgTOTPInvalidSession, http.StatusUnauthorized)
	}

	user, svcErr := s.userService.GetUserByID(ctx, userID)
	if svcErr != nil {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPInvalidSession, constant.ErrMsgTOTPInvalidSession, http.StatusUnauthorized)
	}

	if !user.TotpEnabled {
		return nil, dto.NewServiceError(constant.ErrCodeTOTPNotEnabled, constant.ErrMsgTOTPNotEnabled, http.StatusBadRequest)
	}

	// Try TOTP code first
	codeValid := false
	if len(req.Code) == 6 {
		codeValid = totp.Validate(req.Code, user.TotpSecret)
	}

	// If TOTP code invalid, try backup codes
	if !codeValid {
		codes, listErr := s.totpBackupCodeRepo.ListByUserID(ctx, userID)
		if listErr == nil {
			for _, c := range codes {
				if !c.Used && bcrypt.CompareHashAndPassword([]byte(c.CodeHash), []byte(req.Code)) == nil {
					_ = s.totpBackupCodeRepo.MarkUsed(ctx, c.ID)
					codeValid = true
					break
				}
			}
		}
	}

	if !codeValid {
		logger.WarnCtx(ctx, "Invalid TOTP code", zap.String("user_id", userID.String()))
		return nil, dto.NewServiceError(constant.ErrCodeTOTPInvalidCode, constant.ErrMsgTOTPInvalidCode, http.StatusUnauthorized)
	}

	_ = jwtCfg // used above for ValidateTOTPSessionToken
	opts := &dto.LoginOptions{
		LoginMethod: model.LoginMethodLocal,
		UserAgent:   userAgent,
		IPAddress:   ipAddress,
	}
	loginResp, _, svcErr := s.completeLoginWithOptions(ctx, user, opts)
	if svcErr != nil {
		return nil, svcErr
	}
	return loginResp, nil
}

// generateBackupCodes creates n random 8-char alphanumeric codes and returns (plain, hashed).
func generateBackupCodes(n int) ([]string, []string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	plain := make([]string, n)
	hashed := make([]string, n)
	buf := make([]byte, 8)
	for i := 0; i < n; i++ {
		if _, err := rand.Read(buf); err != nil {
			return nil, nil, err
		}
		code := make([]byte, 8)
		for j, b := range buf {
			code[j] = charset[int(b)%len(charset)]
		}
		plain[i] = string(code)
		h, err := bcrypt.GenerateFromPassword([]byte(plain[i]), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, err
		}
		hashed[i] = string(h)
	}
	return plain, hashed, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, req *request.RefreshTokenRequest) *dto.ServiceError {
	logger.DebugCtx(ctx, "Validating session")

	// Parse and validate the refresh token
	claims, err := util.ValidateRefreshToken(req.RefreshToken, s.cfg.JWTSecret, s.cfg.JWTIssuer)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid refresh token for session validation",
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, constant.ErrMsgInvalidRefreshToken, http.StatusUnauthorized)
	}

	// Check if the refresh token JTI is blacklisted
	isBlacklisted, err := s.tokenBlacklistRepo.IsBlacklisted(ctx, claims.ID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check token blacklist for session validation",
			zap.String("jti", claims.ID),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	if isBlacklisted {
		logger.InfoCtx(ctx, "Session validation failed: refresh token is blacklisted",
			zap.String("jti", claims.ID))
		return dto.NewServiceError(constant.ErrCodeInvalidRefreshToken, "Session has been revoked", http.StatusUnauthorized)
	}

	// Verify the session still exists in the database
	_, err = s.sessionRepo.GetByRefreshTokenJTI(ctx, claims.ID)
	if err != nil {
		logger.InfoCtx(ctx, "Session validation failed: session not found",
			zap.String("jti", claims.ID))
		return dto.NewServiceError(constant.ErrCodeSessionNotFound, constant.ErrMsgSessionNotFound, http.StatusUnauthorized)
	}

	logger.DebugCtx(ctx, "Session validation successful",
		zap.String("jti", claims.ID))
	return nil
}
