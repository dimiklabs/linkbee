package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/request"
	"github.com/shafikshaon/shortlink/response"
	auditSvc "github.com/shafikshaon/shortlink/service/audit"
	authSrv "github.com/shafikshaon/shortlink/service/auth"
	rateLimitSrv "github.com/shafikshaon/shortlink/service/ratelimit"
	"github.com/shafikshaon/shortlink/transport"
	"github.com/shafikshaon/shortlink/util"
)

type AuthHandler struct {
	authService      authSrv.AuthServiceI
	rateLimitService rateLimitSrv.RateLimitServiceI
	auditService     auditSvc.AuditServiceI
}

func NewAuthHandler(authService authSrv.AuthServiceI, rateLimitService rateLimitSrv.RateLimitServiceI, auditService auditSvc.AuditServiceI) *AuthHandler {
	return &AuthHandler{
		authService:      authService,
		rateLimitService: rateLimitService,
		auditService:     auditService,
	}
}

// Signup godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account. A verification email is sent after signup.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.SignupRequest	true	"Signup details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		409		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := c.Request.Context()

	ipAddress := c.ClientIP()

	if h.rateLimitService != nil {
		if rateLimitErr := h.rateLimitService.CheckSignupRateLimit(ctx, ipAddress); rateLimitErr != nil {
			logger.WarnCtx(ctx, "Signup blocked by rate limiter",
				zap.String("ip", ipAddress))
			if rateLimitErr.Data != nil {
				transport.RespondWithErrorData(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description, rateLimitErr.Data)
			} else {
				transport.RespondWithError(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description)
			}
			return
		}
	}

	var req request.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid signup request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Signup request received",
		zap.String("email", req.Email))

	user, svcErr := h.authService.Signup(ctx, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	signupResp := response.SignupResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Status:    user.Status,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	logger.InfoCtx(ctx, "User signed up successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: user.ID, Action: model.AuditActionUserSignup,
		ResourceType: model.AuditResourceUser, ResourceID: user.ID.String(),
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusCreated, "User registered successfully", signupResp)
}

// Login godoc
//
//	@Summary		Login
//	@Description	Authenticates a user and returns JWT access + refresh tokens.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		429		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid login request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	// Extract User-Agent and IP Address for session tracking
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	logger.InfoCtx(ctx, "Login request received",
		zap.String("email", req.Email),
		zap.String("ip", ipAddress))

	// Check rate limits before attempting login
	if h.rateLimitService != nil {
		if rateLimitErr := h.rateLimitService.CheckLoginRateLimit(ctx, req.Email, ipAddress); rateLimitErr != nil {
			logger.WarnCtx(ctx, "Login blocked by rate limiter",
				zap.String("email", req.Email),
				zap.String("ip", ipAddress),
				zap.String("error_code", rateLimitErr.ErrorCode))
			if rateLimitErr.Data != nil {
				transport.RespondWithErrorData(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description, rateLimitErr.Data)
			} else {
				transport.RespondWithError(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description)
			}
			return
		}
	}

	loginResp, sessionsResp, svcErr := h.authService.Login(ctx, &req, userAgent, ipAddress)
	if svcErr != nil {
		// Record failed login attempt for rate limiting
		if h.rateLimitService != nil {
			// Only record for invalid credentials errors
			if svcErr.ErrorCode == constant.ErrCodeInvalidCredentials {
				_ = h.rateLimitService.RecordFailedLogin(ctx, req.Email, ipAddress)
			}
		}

		// Check for max sessions exceeded - return sessions list
		if svcErr.ErrorCode == constant.ErrCodeMaxSessionsExceeded && sessionsResp != nil {
			transport.RespondWithErrorData(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description, sessionsResp)
			return
		}
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	// Login successful - clear rate limits
	if h.rateLimitService != nil {
		_ = h.rateLimitService.RecordSuccessfulLogin(ctx, req.Email, ipAddress)
	}

	logger.InfoCtx(ctx, "User logged in successfully",
		zap.String("user_id", loginResp.User.ID),
		zap.String("email", loginResp.User.Email))

	if uid, err := uuid.Parse(loginResp.User.ID); err == nil {
		h.auditService.LogAsync(auditSvc.LogEntry{
			UserID: uid, Action: model.AuditActionUserLogin,
			ResourceType: model.AuditResourceUser, ResourceID: loginResp.User.ID,
			IPAddress: ipAddress, UserAgent: userAgent,
		})
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Login successful", loginResp)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Exchanges a refresh token for a new access + refresh token pair.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.RefreshTokenRequest	true	"Refresh token"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid refresh token request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Refresh token request received")

	tokenResp, svcErr := h.authService.RefreshToken(ctx, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Token refreshed successfully",
		zap.String("user_id", tokenResp.User.ID))

	transport.RespondWithSuccess(c, http.StatusOK, "Token refreshed successfully", tokenResp)
}

// ChangePassword godoc
//
//	@Summary		Change password
//	@Description	Changes the authenticated user's password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		request.ChangePasswordRequest	true	"Current and new password"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid change password request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Change password request received",
		zap.String("user_id", userID.String()))

	if svcErr := h.authService.ChangePassword(ctx, userID, &req); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Password changed successfully",
		zap.String("user_id", userID.String()))

	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionPasswordChanged,
		ResourceType: model.AuditResourceUser, ResourceID: userID.String(),
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusOK, "Password changed successfully", nil)
}

// ForgotPassword godoc
//
//	@Summary		Forgot password
//	@Description	Sends a password reset email if the address is registered.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.ForgotPasswordRequest	true	"Email address"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid forgot password request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	ipAddress := c.ClientIP()

	if h.rateLimitService != nil {
		if rateLimitErr := h.rateLimitService.CheckForgotPasswordRateLimit(ctx, req.Email, ipAddress); rateLimitErr != nil {
			logger.WarnCtx(ctx, "Forgot-password blocked by rate limiter",
				zap.String("email", req.Email),
				zap.String("ip", ipAddress))
			if rateLimitErr.Data != nil {
				transport.RespondWithErrorData(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description, rateLimitErr.Data)
			} else {
				transport.RespondWithError(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description)
			}
			return
		}
	}

	logger.InfoCtx(ctx, "Forgot password request received",
		zap.String("email", req.Email))

	if svcErr := h.authService.ForgotPassword(ctx, &req); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Forgot password processed",
		zap.String("email", req.Email))

	transport.RespondWithSuccess(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
}

// ResetPassword godoc
//
//	@Summary		Reset password
//	@Description	Resets a user's password using the token received by email.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.ResetPasswordRequest	true	"Reset token and new password"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid reset password request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Reset password request received")

	if svcErr := h.authService.ResetPassword(ctx, &req); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Password reset successful")

	transport.RespondWithSuccess(c, http.StatusOK, "Password has been reset successfully", nil)
}

// GetProfile godoc
//
//	@Summary		Get user profile
//	@Description	Returns the authenticated user's profile information.
//	@Tags			auth
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	logger.DebugCtx(ctx, "Get profile request received",
		zap.String("user_id", userID.String()))

	profile, svcErr := h.authService.GetProfile(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.DebugCtx(ctx, "Profile retrieved successfully",
		zap.String("user_id", userID.String()))

	transport.RespondWithSuccess(c, http.StatusOK, "Profile retrieved successfully", profile)
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Updates the authenticated user's display name and other profile fields.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		request.UpdateProfileRequest	true	"Profile fields to update"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid update profile request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Update profile request received",
		zap.String("user_id", userID.String()))

	profile, svcErr := h.authService.UpdateProfile(ctx, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Profile updated successfully",
		zap.String("user_id", userID.String()))

	transport.RespondWithSuccess(c, http.StatusOK, "Profile updated successfully", profile)
}

// Logout godoc
//
//	@Summary		Logout
//	@Description	Revokes the current access token and optionally the refresh token.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		request.LogoutRequest	false	"Optional refresh token to revoke"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		logger.WarnCtx(ctx, "Missing or invalid authorization header")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	accessToken := authHeader[7:]

	// Optionally get refresh token from request body
	var req request.LogoutRequest
	_ = c.ShouldBindJSON(&req) // Ignore error - refresh token is optional

	logger.InfoCtx(ctx, "Logout request received")

	if svcErr := h.authService.Logout(ctx, accessToken, req.RefreshToken); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Logout successful")

	if userIDStr, ok := c.Get(middlewares.ContextKeyUserID); ok {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			h.auditService.LogAsync(auditSvc.LogEntry{
				UserID: uid, Action: model.AuditActionUserLogout,
				ResourceType: model.AuditResourceUser, ResourceID: uid.String(),
				IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
			})
		}
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Logged out successfully", nil)
}

func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	logger.InfoCtx(ctx, "Delete account request received",
		zap.String("user_id", userID.String()))

	if svcErr := h.authService.DeleteAccount(ctx, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Account deletion scheduled",
		zap.String("user_id", userID.String()))

	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionAccountDeleted,
		ResourceType: model.AuditResourceUser, ResourceID: userID.String(),
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusOK, "Account scheduled for deletion. You have 30 days to login and reactivate your account.", nil)
}

func (h *AuthHandler) ReactivateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ReactivateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid reactivate account request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	// Extract User-Agent and IP Address for session tracking
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	logger.InfoCtx(ctx, "Reactivate account request received",
		zap.String("email", req.Email),
		zap.String("ip", ipAddress))

	// Check rate limits before attempting reactivation (uses same limits as login)
	if h.rateLimitService != nil {
		if rateLimitErr := h.rateLimitService.CheckLoginRateLimit(ctx, req.Email, ipAddress); rateLimitErr != nil {
			logger.WarnCtx(ctx, "Reactivation blocked by rate limiter",
				zap.String("email", req.Email),
				zap.String("ip", ipAddress))
			if rateLimitErr.Data != nil {
				transport.RespondWithErrorData(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description, rateLimitErr.Data)
			} else {
				transport.RespondWithError(c, rateLimitErr.StatusCode, rateLimitErr.ErrorCode, rateLimitErr.Description)
			}
			return
		}
	}

	loginResp, sessionsResp, svcErr := h.authService.ReactivateAccount(ctx, &req, userAgent, ipAddress)
	if svcErr != nil {
		// Record failed attempt for rate limiting (if invalid credentials)
		if h.rateLimitService != nil && svcErr.ErrorCode == constant.ErrCodeInvalidCredentials {
			_ = h.rateLimitService.RecordFailedLogin(ctx, req.Email, ipAddress)
		}

		// Check for max sessions exceeded - return sessions list
		if svcErr.ErrorCode == constant.ErrCodeMaxSessionsExceeded && sessionsResp != nil {
			transport.RespondWithErrorData(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description, sessionsResp)
			return
		}
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	// Reactivation successful - clear rate limits
	if h.rateLimitService != nil {
		_ = h.rateLimitService.RecordSuccessfulLogin(ctx, req.Email, ipAddress)
	}

	logger.InfoCtx(ctx, "Account reactivated and logged in successfully",
		zap.String("user_id", loginResp.User.ID),
		zap.String("email", loginResp.User.Email))

	transport.RespondWithSuccess(c, http.StatusOK, "Account reactivated successfully", loginResp)
}

func (h *AuthHandler) GetSessions(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	// Extract the refresh token JTI from the X-Refresh-Token header.
	// The access token JTI (from auth middleware) differs from the refresh token JTI
	// stored in sessions, so we need the actual refresh token to identify the current session.
	currentRefreshJTI := ""
	if refreshToken := c.GetHeader("X-Refresh-Token"); refreshToken != "" {
		if claims, err := util.ValidateRefreshToken(refreshToken, h.authService.GetJWTSecret(), h.authService.GetJWTIssuer()); err == nil {
			currentRefreshJTI = claims.ID
		}
	}

	logger.DebugCtx(ctx, "Get sessions request received",
		zap.String("user_id", userID.String()))

	sessionsResp, svcErr := h.authService.GetSessions(ctx, userID, currentRefreshJTI)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.DebugCtx(ctx, "Sessions retrieved successfully",
		zap.String("user_id", userID.String()),
		zap.Int("count", sessionsResp.Count))

	transport.RespondWithSuccess(c, http.StatusOK, "Sessions retrieved successfully", sessionsResp)
}

func (h *AuthHandler) DeleteSession(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	sessionIDStr := c.Param("id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid session ID",
			zap.String("session_id", sessionIDStr),
			zap.Error(err))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid session ID")
		return
	}

	logger.InfoCtx(ctx, "Delete session request received",
		zap.String("user_id", userID.String()),
		zap.String("session_id", sessionID.String()))

	if svcErr := h.authService.DeleteSession(ctx, userID, sessionID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Session deleted successfully",
		zap.String("user_id", userID.String()),
		zap.String("session_id", sessionID.String()))

	transport.RespondWithSuccess(c, http.StatusOK, "Session logged out successfully", nil)
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		logger.WarnCtx(ctx, "User ID not found in context")
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		logger.WarnCtx(ctx, "Invalid user ID in context",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	logger.InfoCtx(ctx, "Logout all devices request received",
		zap.String("user_id", userID.String()))

	if svcErr := h.authService.LogoutAll(ctx, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Logged out from all devices successfully",
		zap.String("user_id", userID.String()))

	transport.RespondWithSuccess(c, http.StatusOK, "Logged out from all devices successfully", nil)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid verify email request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Verify email request received")

	if svcErr := h.authService.VerifyEmail(ctx, req.Token); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Email verified successfully")

	transport.RespondWithSuccess(c, http.StatusOK, "Email verified successfully", nil)
}

func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid resend verification request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	logger.InfoCtx(ctx, "Resend verification email request received",
		zap.String("email", req.Email))

	if svcErr := h.authService.ResendVerificationEmail(ctx, req.Email); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Resend verification email processed",
		zap.String("email", req.Email))

	transport.RespondWithSuccess(c, http.StatusOK, "If the email exists and is not verified, a verification email has been sent", nil)
}

func (h *AuthHandler) GetRateLimitStatus(c *gin.Context) {
	ctx := c.Request.Context()

	email := c.Query("email")
	if email == "" {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Email is required")
		return
	}

	ipAddress := c.ClientIP()

	if h.rateLimitService == nil {
		transport.RespondWithSuccess(c, http.StatusOK, "Rate limiting is disabled", &response.RateLimitStatusResponse{
			AccountMaxAttempts:       0,
			AccountRemainingAttempts: -1, // -1 indicates unlimited
			IPMaxAttempts:            0,
			IPRemainingAttempts:      -1,
		})
		return
	}

	status, err := h.rateLimitService.GetRateLimitStatus(ctx, email, ipAddress)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get rate limit status", zap.Error(err))
		transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Rate limit status retrieved", status)
}

func (h *AuthHandler) GetTOTPStatus(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	resp, svcErr := h.authService.GetTOTPStatus(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "TOTP status retrieved", resp)
}

func (h *AuthHandler) SetupTOTP(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	resp, svcErr := h.authService.SetupTOTP(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "TOTP setup initiated. Scan the QR code and confirm with a code.", resp)
}

func (h *AuthHandler) ConfirmTOTP(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	var req request.TOTPConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}
	resp, svcErr := h.authService.ConfirmTOTP(ctx, userID, req.Code)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Two-factor authentication enabled. Save your backup codes.", resp)
}

func (h *AuthHandler) DisableTOTP(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, exists := c.Get(middlewares.ContextKeyUserID)
	if !exists {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}
	var req request.TOTPDisableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}
	if svcErr := h.authService.DisableTOTP(ctx, userID, req.Password); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Two-factor authentication disabled", nil)
}

func (h *AuthHandler) VerifyTOTPLogin(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.TOTPVerifyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}
	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	loginResp, svcErr := h.authService.VerifyTOTPLogin(ctx, &req, userAgent, ipAddress)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	if uid, err := uuid.Parse(loginResp.User.ID); err == nil {
		h.auditService.LogAsync(auditSvc.LogEntry{
			UserID: uid, Action: model.AuditActionUserLogin,
			ResourceType: model.AuditResourceUser, ResourceID: loginResp.User.ID,
			IPAddress: ipAddress, UserAgent: userAgent,
		})
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Login successful", loginResp)
}

func (h *AuthHandler) ValidateSession(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnCtx(ctx, "Invalid validate session request",
			zap.Error(err))
		errCode, errMsg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, errCode, errMsg)
		return
	}

	if svcErr := h.authService.ValidateSession(ctx, &req); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Session is valid", nil)
}
