package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	githubAuthURL      = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserURL      = "https://api.github.com/user"
	githubUserEmailURL = "https://api.github.com/user/emails"
)

// GitHubTokenResponse represents GitHub's token response
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHubUserInfo represents user info from GitHub
type GitHubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubEmail represents an email from GitHub's email API
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type GitHubOAuthServiceI interface {
	// GetAuthorizationURL generates the GitHub OAuth authorization URL
	GetAuthorizationURL(ctx context.Context, ipAddress, userAgent, redirectURI string) (string, error)

	// HandleCallback processes the OAuth callback and returns login response
	HandleCallback(ctx context.Context, state, code, ipAddress, userAgent string) (*response.LoginResponse, *dto.ServiceError)
}

type GitHubOAuthService struct {
	githubConfig       *config.GitHubOAuthConfig
	appConfig          *config.AppConfig
	userRepo           repository.UserRepositoryI
	oauthStateRepo     repository.OAuthStateRepositoryI
	sessionRepo        repository.SessionRepositoryI
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI
	httpClient         *http.Client
}

func NewGitHubOAuthService(
	githubConfig *config.GitHubOAuthConfig,
	appConfig *config.AppConfig,
	userRepo repository.UserRepositoryI,
	oauthStateRepo repository.OAuthStateRepositoryI,
	sessionRepo repository.SessionRepositoryI,
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI,
) GitHubOAuthServiceI {
	return &GitHubOAuthService{
		githubConfig:       githubConfig,
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

func (s *GitHubOAuthService) GetAuthorizationURL(ctx context.Context, ipAddress, userAgent, redirectURI string) (string, error) {
	logger.InfoCtx(ctx, "Generating GitHub OAuth authorization URL")

	if !s.githubConfig.Enabled {
		logger.WarnCtx(ctx, "GitHub OAuth is not enabled")
		return "", fmt.Errorf("github OAuth is not enabled")
	}

	// Create OAuth state (GitHub doesn't support PKCE, but we still use state for CSRF protection)
	stateData, err := s.oauthStateRepo.Create(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to create OAuth state", zap.Error(err))
		return "", err
	}

	// Build authorization URL
	params := url.Values{}
	params.Set("client_id", s.githubConfig.ClientID)
	params.Set("redirect_uri", s.githubConfig.RedirectURL)
	params.Set("scope", "user:email read:user")
	params.Set("state", stateData.State)
	params.Set("allow_signup", "true")

	authURL := fmt.Sprintf("%s?%s", githubAuthURL, params.Encode())

	logger.InfoCtx(ctx, "GitHub OAuth authorization URL generated",
		zap.String("state", stateData.State[:8]+"..."))

	return authURL, nil
}

func (s *GitHubOAuthService) HandleCallback(ctx context.Context, state, code, ipAddress, userAgent string) (*response.LoginResponse, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Processing GitHub OAuth callback")

	if !s.githubConfig.Enabled {
		logger.WarnCtx(ctx, "GitHub OAuth is not enabled")
		return nil, dto.NewServiceError(constant.ErrCodeGitHubOAuthDisabled, constant.ErrMsgGitHubOAuthDisabled, http.StatusBadRequest)
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

	// Get user info from GitHub
	userInfo, err := s.getUserInfo(ctx, tokens.AccessToken)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user info from GitHub", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthProviderError, constant.ErrMsgOAuthProviderError, http.StatusBadGateway)
	}

	// Get user's primary verified email if not provided in user info
	if userInfo.Email == "" {
		email, err := s.getPrimaryEmail(ctx, tokens.AccessToken)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to get user email from GitHub", zap.Error(err))
			return nil, dto.NewServiceError(constant.ErrCodeOAuthEmailNotFound, constant.ErrMsgOAuthEmailNotFound, http.StatusBadRequest)
		}
		userInfo.Email = email
	}

	// Validate email is present
	if userInfo.Email == "" {
		logger.WarnCtx(ctx, "No email in GitHub user info")
		return nil, dto.NewServiceError(constant.ErrCodeOAuthEmailNotFound, constant.ErrMsgOAuthEmailNotFound, http.StatusBadRequest)
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

	logger.InfoCtx(ctx, "GitHub OAuth login successful",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	return loginResp, nil
}

func (s *GitHubOAuthService) exchangeCodeForTokens(ctx context.Context, code string) (*GitHubTokenResponse, error) {
	logger.DebugCtx(ctx, "Exchanging authorization code for tokens")

	data := url.Values{}
	data.Set("client_id", s.githubConfig.ClientID)
	data.Set("client_secret", s.githubConfig.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.githubConfig.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, githubTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

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

	var tokens GitHubTokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokens.AccessToken == "" {
		return nil, fmt.Errorf("no access token in response")
	}

	return &tokens, nil
}

func (s *GitHubOAuthService) getUserInfo(ctx context.Context, accessToken string) (*GitHubUserInfo, error) {
	logger.DebugCtx(ctx, "Fetching user info from GitHub")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

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

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return &userInfo, nil
}

func (s *GitHubOAuthService) getPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	logger.DebugCtx(ctx, "Fetching primary email from GitHub")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubUserEmailURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create email request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send email request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read email response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorCtx(ctx, "Failed to get user emails",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)))
		return "", fmt.Errorf("email request failed with status %d", resp.StatusCode)
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("failed to parse email response: %w", err)
	}

	// Find primary verified email
	for _, email := range emails {
		if email.Primary && email.Verified {
			return email.Email, nil
		}
	}

	// Fall back to any verified email
	for _, email := range emails {
		if email.Verified {
			return email.Email, nil
		}
	}

	return "", fmt.Errorf("no verified email found")
}

func (s *GitHubOAuthService) findOrCreateUser(ctx context.Context, githubUser *GitHubUserInfo) (*model.User, *dto.ServiceError) {
	githubIDStr := strconv.FormatInt(githubUser.ID, 10)

	logger.DebugCtx(ctx, "Finding or creating user for GitHub login",
		zap.String("github_id", githubIDStr),
		zap.String("email", githubUser.Email))

	// First, try to find a user by GitHub ID
	user, err := s.userRepo.GetByGitHubID(ctx, githubIDStr)
	if err == nil {
		// User found by GitHub ID
		logger.DebugCtx(ctx, "User found by GitHub ID",
			zap.String("user_id", user.ID.String()))

		// Update last login
		_ = s.userRepo.UpdateLastLogin(ctx, user.ID)
		return user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ErrorCtx(ctx, "Error finding user by GitHub ID", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// User not found by GitHub ID, try to find by email
	user, err = s.userRepo.GetByEmail(ctx, githubUser.Email)
	if err == nil {
		// User exists with this email but GitHub is not linked
		// Require explicit linking from profile
		logger.WarnCtx(ctx, "User exists but GitHub account not linked",
			zap.String("email", githubUser.Email),
			zap.String("user_id", user.ID.String()))
		return nil, dto.NewServiceError(constant.ErrCodeOAuthNotLinked, constant.ErrMsgOAuthLoginNotLinked, http.StatusForbidden)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.ErrorCtx(ctx, "Error finding user by email", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	// Create a new user
	logger.InfoCtx(ctx, "Creating new user from GitHub OAuth",
		zap.String("email", githubUser.Email))

	now := time.Now()

	// Parse name into first and last name
	firstName, lastName := parseGitHubName(githubUser.Name, githubUser.Login)

	newUser := &model.User{
		Email:                githubUser.Email,
		FirstName:            firstName,
		LastName:             lastName,
		ProfilePicture:       githubUser.AvatarURL,
		ProfilePictureSource: model.AuthProviderGitHub,
		Status:               "active",
		Role:                 "user",
		AuthProvider:         model.AuthProviderGitHub,
		GitHubID:             &githubIDStr,
		EmailVerified:        true, // GitHub verified the email
		EmailVerifiedAt:      &now,
		LastLogin:            &now,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		logger.ErrorCtx(ctx, "Failed to create user", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User created from GitHub OAuth",
		zap.String("user_id", newUser.ID.String()),
		zap.String("email", newUser.Email))

	return newUser, nil
}

func parseGitHubName(name, login string) (firstName, lastName string) {
	if name == "" {
		return login, ""
	}

	parts := strings.SplitN(name, " ", 2)
	firstName = parts[0]
	if len(parts) > 1 {
		lastName = parts[1]
	}
	return
}

func (s *GitHubOAuthService) completeLogin(ctx context.Context, user *model.User, userAgent, ipAddress string) (*response.LoginResponse, *dto.ServiceError) {
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
