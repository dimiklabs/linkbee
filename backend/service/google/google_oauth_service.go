package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/response"
	"github.com/shafikshaon/linkbee/util"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// GoogleTokenResponse represents Google's token response
type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// GoogleUserInfo represents user info from Google
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type GoogleOAuthServiceI interface {
	// GetAuthorizationURL generates the Google OAuth authorization URL
	GetAuthorizationURL(ctx context.Context, ipAddress, userAgent, redirectURI string) (string, error)

	// HandleCallback processes the OAuth callback and returns login response
	HandleCallback(ctx context.Context, state, code, ipAddress, userAgent string) (*response.LoginResponse, *dto.ServiceError)
}

type GoogleOAuthService struct {
	googleConfig       *config.GoogleOAuthConfig
	appConfig          *config.AppConfig
	userRepo           repository.UserRepositoryI
	oauthStateRepo     repository.OAuthStateRepositoryI
	sessionRepo        repository.SessionRepositoryI
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI
	httpClient         *http.Client
}

func NewGoogleOAuthService(
	googleConfig *config.GoogleOAuthConfig,
	appConfig *config.AppConfig,
	userRepo repository.UserRepositoryI,
	oauthStateRepo repository.OAuthStateRepositoryI,
	sessionRepo repository.SessionRepositoryI,
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI,
) GoogleOAuthServiceI {
	return &GoogleOAuthService{
		googleConfig:       googleConfig,
		appConfig:          appConfig,
		userRepo:           userRepo,
		oauthStateRepo:     oauthStateRepo,
		sessionRepo:        sessionRepo,
		tokenBlacklistRepo: tokenBlacklistRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *GoogleOAuthService) GetAuthorizationURL(ctx context.Context, ipAddress, userAgent, redirectURI string) (string, error) {
	logger.InfoCtx(ctx, "Generating Google OAuth authorization URL")

	if !s.googleConfig.Enabled {
		logger.WarnCtx(ctx, "Google OAuth is not enabled")
		return "", fmt.Errorf("google OAuth is not enabled")
	}

	// Create OAuth state with PKCE
	stateData, err := s.oauthStateRepo.Create(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to create OAuth state", zap.Error(err))
		return "", err
	}

	// Generate PKCE code challenge
	codeChallenge := s.oauthStateRepo.GenerateCodeChallenge(stateData.CodeVerifier)

	// Build authorization URL with security parameters
	params := url.Values{}
	params.Set("client_id", s.googleConfig.ClientID)
	params.Set("redirect_uri", s.googleConfig.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", stateData.State)
	params.Set("nonce", stateData.Nonce)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	params.Set("access_type", "offline")   // Get refresh token
	params.Set("prompt", "select_account") // Always show an account selector

	authURL := fmt.Sprintf("%s?%s", googleAuthURL, params.Encode())

	logger.InfoCtx(ctx, "Google OAuth authorization URL generated",
		zap.String("state", stateData.State[:8]+"..."))

	return authURL, nil
}

func (s *GoogleOAuthService) HandleCallback(ctx context.Context, state, code, ipAddress, userAgent string) (*response.LoginResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing Google OAuth callback")

	if !s.googleConfig.Enabled {
		logger.WarnCtx(ctx, "Google OAuth is not enabled")
		return nil, dto.NewServiceError(constant.ErrCodeOAuthDisabled, constant.ErrMsgOAuthDisabled, http.StatusBadRequest)
	}

	// Retrieve and validate OAuth state
	stateData, err := s.oauthStateRepo.Get(ctx, state)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid OAuth state", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthInvalidState, constant.ErrMsgOAuthInvalidState, http.StatusBadRequest)
	}

	// Delete state immediately after retrieval (single use)
	_ = s.oauthStateRepo.Delete(ctx, state)

	// Verify client binding (optional but recommended)
	if stateData.IPAddress != ipAddress {
		logger.WarnCtx(ctx, "OAuth state IP mismatch",
			zap.String("expected_ip", stateData.IPAddress),
			zap.String("actual_ip", ipAddress))
		// Log warning but don't fail - IP can change
	}

	// Exchange authorization code for tokens
	tokens, err := s.exchangeCodeForTokens(ctx, code, stateData.CodeVerifier)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to exchange code for tokens", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthTokenExchange, constant.ErrMsgOAuthTokenExchange, http.StatusBadRequest)
	}

	// Get user info from Google
	userInfo, err := s.getUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user info from Google", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthProviderError, constant.ErrMsgOAuthProviderError, http.StatusBadGateway)
	}

	// Validate email is present and verified
	if userInfo.Email == "" {
		logger.WarnCtx(ctx, "No email in Google user info")
		return nil, dto.NewServiceError(constant.ErrCodeOAuthEmailNotFound, constant.ErrMsgOAuthEmailNotFound, http.StatusBadRequest)
	}

	if !userInfo.EmailVerified {
		logger.WarnCtx(ctx, "Google email not verified",
			zap.String("email", userInfo.Email))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthInvalidToken, "Google email is not verified", http.StatusBadRequest)
	}

	// Find or create user
	user, svcErr := s.findOrCreateUser(ctx, userInfo)
	if svcErr != nil {
		return nil, svcErr
	}

	// Check if user is active
	if user.Status != "active" {
		logger.WarnCtx(ctx, "User account is not active",
			zap.String("user_id", user.ID.String()),
			zap.String("status", user.Status))
		return nil, dto.NewServiceError(constant.ErrCodeUserInactive, constant.ErrMsgUserInactive, http.StatusForbidden)
	}

	// Complete login - generate tokens and create a session
	loginResp, svcErr := s.completeLogin(ctx, user, userAgent, ipAddress)
	if svcErr != nil {
		return nil, svcErr
	}

	logger.InfoCtx(ctx, "Google OAuth login successful",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	return loginResp, nil
}

func (s *GoogleOAuthService) exchangeCodeForTokens(ctx context.Context, code, codeVerifier string) (*GoogleTokenResponse, error) {
	logger.DebugCtx(ctx, "Exchanging authorization code for tokens")

	data := url.Values{}
	data.Set("client_id", s.googleConfig.ClientID)
	data.Set("client_secret", s.googleConfig.ClientSecret)
	data.Set("code", code)
	data.Set("code_verifier", codeVerifier)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", s.googleConfig.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorCtx(ctx, "Token exchange failed",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)))
		return nil, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	var tokens GoogleTokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokens, nil
}

func (s *GoogleOAuthService) getUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	logger.DebugCtx(ctx, "Fetching user info from Google")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send userinfo request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorCtx(ctx, "Failed to get user info",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)))
		return nil, fmt.Errorf("userinfo request failed with status %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return &userInfo, nil
}

func (s *GoogleOAuthService) findOrCreateUser(ctx context.Context, googleUser *GoogleUserInfo) (*model.User, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Finding or creating user for Google login",
		zap.String("google_id", googleUser.Sub),
		zap.String("email", googleUser.Email))

	// First, try to find a user by Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleUser.Sub)
	if err == nil {
		// User found by Google ID
		logger.DebugCtx(ctx, "User found by Google ID",
			zap.String("user_id", user.ID.String()))

		// Update last login
		_ = s.userRepo.UpdateLastLogin(ctx, user.ID)
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ErrorCtx(ctx, "Error finding user by Google ID", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// User not found by Google ID, try to find by email
	normalizedEmail := util.NormalizeEmail(googleUser.Email)
	user, err = s.userRepo.GetByEmail(ctx, normalizedEmail)
	if err == nil {
		// An account with this email exists but is not linked to Google.
		// Tell the user to sign in with their existing method.
		logger.WarnCtx(ctx, "Google OAuth: email already registered with another provider",
			zap.String("email", normalizedEmail),
			zap.String("user_id", user.ID.String()))
		return nil, dto.NewServiceError(constant.ErrCodeEmailAlreadyExists, constant.ErrMsgOAuthEmailAlreadyExists, http.StatusConflict)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ErrorCtx(ctx, "Error finding user by email", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Create a new user
	logger.InfoCtx(ctx, "Creating new user from Google OAuth",
		zap.String("email", normalizedEmail))

	now := time.Now()
	newUser := &model.User{
		Email:                normalizedEmail,
		FirstName:            googleUser.GivenName,
		LastName:             googleUser.FamilyName,
		ProfilePicture:       googleUser.Picture,
		ProfilePictureSource: model.AuthProviderGoogle,
		Status:               "active",
		Role:                 "user",
		AuthProvider:         model.AuthProviderGoogle,
		GoogleID:             &googleUser.Sub,
		EmailVerified:        true, // Google verified the email
		EmailVerifiedAt:      &now,
		LastLogin:            &now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		logger.ErrorCtx(ctx, "Failed to create user", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User created from Google OAuth",
		zap.String("user_id", newUser.ID.String()),
		zap.String("email", newUser.Email))

	return newUser, nil
}

func (s *GoogleOAuthService) completeLogin(ctx context.Context, user *model.User, userAgent, ipAddress string) (*response.LoginResponse, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Completing login for user",
		zap.String("user_id", user.ID.String()))

	// Check session limit and auto-logout oldest if at max
	activeCount, err := s.sessionRepo.CountActiveByUserID(ctx, user.ID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to count active sessions", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	if activeCount >= int64(s.appConfig.MaxSessions) {
		logger.InfoCtx(ctx, "Max sessions reached, logging out oldest session",
			zap.String("user_id", user.ID.String()),
			zap.Int64("active_sessions", activeCount),
			zap.Int("max_sessions", s.appConfig.MaxSessions))

		// Get oldest session
		oldestSession, err := s.sessionRepo.GetOldestActiveByUserID(ctx, user.ID)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to get oldest session", zap.Error(err))
			return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}

		// Blacklist the refresh token
		refreshExpiry := time.Until(oldestSession.ExpiresAt)
		if refreshExpiry > 0 {
			if err := s.tokenBlacklistRepo.Add(ctx, oldestSession.RefreshTokenJTI, refreshExpiry); err != nil {
				logger.ErrorCtx(ctx, "Failed to blacklist refresh token", zap.Error(err))
				// Continue anyway - session will be deleted
			}
		}

		// Delete the oldest session
		if err := s.sessionRepo.Delete(ctx, oldestSession.ID); err != nil {
			logger.ErrorCtx(ctx, "Failed to delete oldest session", zap.Error(err))
			return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}

		logger.InfoCtx(ctx, "Oldest session logged out automatically",
			zap.String("user_id", user.ID.String()),
			zap.String("session_id", oldestSession.ID.String()))
	}

	// Generate tokens
	jwtConfig := &util.JWTConfig{
		Secret:              s.appConfig.JWTSecret,
		Issuer:              s.appConfig.JWTIssuer,
		AccessExpiryMinutes: s.appConfig.JWTAccessExpiry,
		RefreshExpiryDays:   s.appConfig.JWTRefreshExpiry,
	}

	tokenPair, err := util.GenerateTokenPair(jwtConfig, user.ID.String(), user.Email, user.Role)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate token pair", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Extract refresh token JTI for session
	refreshClaims, err := util.ValidateRefreshToken(tokenPair.RefreshToken, jwtConfig.Secret, jwtConfig.Issuer)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to validate generated refresh token", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Create session
	session := &model.Session{
		UserID:          user.ID,
		RefreshTokenJTI: refreshClaims.ID,
		UserAgent:       userAgent,
		IPAddress:       ipAddress,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(time.Duration(s.appConfig.JWTRefreshExpiry) * 24 * time.Hour),
		LastActivityAt:  time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		logger.ErrorCtx(ctx, "Failed to create session", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

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
