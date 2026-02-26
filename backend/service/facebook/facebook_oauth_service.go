package facebook

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

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/response"
	"github.com/shafikshaon/shortlink/util"
)

const (
	facebookAuthURL     = "https://www.facebook.com/v18.0/dialog/oauth"
	facebookTokenURL    = "https://graph.facebook.com/v18.0/oauth/access_token"
	facebookUserInfoURL = "https://graph.facebook.com/v18.0/me"
)

// FacebookTokenResponse represents Facebook's token response
type FacebookTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// FacebookUserInfo represents user info from Facebook
type FacebookUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Picture   struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

type FacebookOAuthServiceI interface {
	// GetAuthorizationURL generates the Facebook OAuth authorization URL
	GetAuthorizationURL(ctx context.Context, ipAddress, userAgent, redirectURI string) (string, error)

	// HandleCallback processes the OAuth callback and returns login response
	HandleCallback(ctx context.Context, state, code, ipAddress, userAgent string) (*response.LoginResponse, *dto.ServiceError)
}

type FacebookOAuthService struct {
	facebookConfig     *config.FacebookOAuthConfig
	appConfig          *config.AppConfig
	userRepo           repository.UserRepositoryI
	oauthStateRepo     repository.OAuthStateRepositoryI
	sessionRepo        repository.SessionRepositoryI
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI
	httpClient         *http.Client
}

func NewFacebookOAuthService(
	facebookConfig *config.FacebookOAuthConfig,
	appConfig *config.AppConfig,
	userRepo repository.UserRepositoryI,
	oauthStateRepo repository.OAuthStateRepositoryI,
	sessionRepo repository.SessionRepositoryI,
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI,
) FacebookOAuthServiceI {
	return &FacebookOAuthService{
		facebookConfig:     facebookConfig,
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

func (s *FacebookOAuthService) GetAuthorizationURL(ctx context.Context, ipAddress, userAgent, redirectURI string) (string, error) {
	logger.InfoCtx(ctx, "Generating Facebook OAuth authorization URL")

	if !s.facebookConfig.Enabled {
		logger.WarnCtx(ctx, "Facebook OAuth is not enabled")
		return "", fmt.Errorf("facebook OAuth is not enabled")
	}

	// Create OAuth state for CSRF protection
	stateData, err := s.oauthStateRepo.Create(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to create OAuth state", zap.Error(err))
		return "", err
	}

	// Build authorization URL
	params := url.Values{}
	params.Set("client_id", s.facebookConfig.ClientID)
	params.Set("redirect_uri", s.facebookConfig.RedirectURL)
	params.Set("scope", "email public_profile")
	params.Set("state", stateData.State)
	params.Set("response_type", "code")

	authURL := fmt.Sprintf("%s?%s", facebookAuthURL, params.Encode())

	logger.InfoCtx(ctx, "Facebook OAuth authorization URL generated",
		zap.String("state", stateData.State[:8]+"..."))

	return authURL, nil
}

func (s *FacebookOAuthService) HandleCallback(ctx context.Context, state, code, ipAddress, userAgent string) (*response.LoginResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing Facebook OAuth callback")

	if !s.facebookConfig.Enabled {
		logger.WarnCtx(ctx, "Facebook OAuth is not enabled")
		return nil, dto.NewServiceError(constant.ErrCodeFacebookOAuthDisabled, constant.ErrMsgFacebookOAuthDisabled, http.StatusBadRequest)
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
	tokens, err := s.exchangeCodeForTokens(ctx, code)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to exchange code for tokens", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthTokenExchange, constant.ErrMsgOAuthTokenExchange, http.StatusBadRequest)
	}

	// Get user info from Facebook
	userInfo, err := s.getUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user info from Facebook", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthProviderError, constant.ErrMsgOAuthProviderError, http.StatusBadGateway)
	}

	// Validate email is present
	if userInfo.Email == "" {
		logger.WarnCtx(ctx, "No email in Facebook user info")
		return nil, dto.NewServiceError(constant.ErrCodeOAuthEmailNotFound, constant.ErrMsgOAuthEmailNotFound, http.StatusBadRequest)
	}

	// Find or create a user
	user, svcErr := s.findOrCreateUser(ctx, userInfo)
	if svcErr != nil {
		return nil, svcErr
	}

	// Check if the user is active
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

	logger.InfoCtx(ctx, "Facebook OAuth login successful",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	return loginResp, nil
}

func (s *FacebookOAuthService) exchangeCodeForTokens(ctx context.Context, code string) (*FacebookTokenResponse, error) {
	logger.DebugCtx(ctx, "Exchanging authorization code for tokens")

	params := url.Values{}
	params.Set("client_id", s.facebookConfig.ClientID)
	params.Set("client_secret", s.facebookConfig.ClientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", s.facebookConfig.RedirectURL)

	reqURL := fmt.Sprintf("%s?%s", facebookTokenURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

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

	var tokens FacebookTokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokens.AccessToken == "" {
		return nil, fmt.Errorf("no access token in response")
	}

	return &tokens, nil
}

func (s *FacebookOAuthService) getUserInfo(ctx context.Context, accessToken string) (*FacebookUserInfo, error) {
	logger.DebugCtx(ctx, "Fetching user info from Facebook")

	params := url.Values{}
	params.Set("fields", "id,email,name,first_name,last_name,picture")
	params.Set("access_token", accessToken)

	reqURL := fmt.Sprintf("%s?%s", facebookUserInfoURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

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

	var userInfo FacebookUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return &userInfo, nil
}

func (s *FacebookOAuthService) findOrCreateUser(ctx context.Context, facebookUser *FacebookUserInfo) (*model.User, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Finding or creating user for Facebook login",
		zap.String("facebook_id", facebookUser.ID),
		zap.String("email", facebookUser.Email))

	// First, try to find a user by Facebook ID
	user, err := s.userRepo.GetByFacebookID(ctx, facebookUser.ID)
	if err == nil {
		// User found by Facebook ID
		logger.DebugCtx(ctx, "User found by Facebook ID",
			zap.String("user_id", user.ID.String()))

		// Update last login
		_ = s.userRepo.UpdateLastLogin(ctx, user.ID)
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ErrorCtx(ctx, "Error finding user by Facebook ID", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// User not found by Facebook ID, try to find by email
	user, err = s.userRepo.GetByEmail(ctx, facebookUser.Email)
	if err == nil {
		// User exists with this email but Facebook is not linked
		// Require explicit linking from profile
		logger.WarnCtx(ctx, "User exists but Facebook account not linked",
			zap.String("email", facebookUser.Email),
			zap.String("user_id", user.ID.String()))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthNotLinked, constant.ErrMsgOAuthLoginNotLinked, http.StatusForbidden)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ErrorCtx(ctx, "Error finding user by email", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Create a new user
	logger.InfoCtx(ctx, "Creating new user from Facebook OAuth",
		zap.String("email", facebookUser.Email))

	now := time.Now()

	// Parse name - Facebook provides first_name and last_name directly
	firstName := facebookUser.FirstName
	lastName := facebookUser.LastName
	if firstName == "" && facebookUser.Name != "" {
		firstName, lastName = parseFacebookName(facebookUser.Name)
	}

	newUser := &model.User{
		Email:                facebookUser.Email,
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       facebookUser.Picture.Data.URL,
		ProfilePictureSource: model.AuthProviderFacebook,
		Status:               "active",
		Role:                 "user",
		AuthProvider:         model.AuthProviderFacebook,
		FacebookID:           &facebookUser.ID,
		EmailVerified:        true, // Facebook verified the email
		EmailVerifiedAt:      &now,
		LastLogin:            &now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		logger.ErrorCtx(ctx, "Failed to create user", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User created from Facebook OAuth",
		zap.String("user_id", newUser.ID.String()),
		zap.String("email", newUser.Email))

	return newUser, nil
}

func parseFacebookName(name string) (firstName, lastName string) {
	parts := strings.SplitN(name, " ", 2)
	firstName = parts[0]
	if len(parts) > 1 {
		lastName = parts[1]
	}
	return
}

func (s *FacebookOAuthService) completeLogin(ctx context.Context, user *model.User, userAgent, ipAddress string) (*response.LoginResponse, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Completing login for user",
		zap.String("user_id", user.ID.String()))

	// Check session limit and auto-logout the oldest if at max
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
				// Continue anyway - the session will be deleted
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
