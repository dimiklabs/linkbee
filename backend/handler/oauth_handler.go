package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/logger"
	githubSrv "github.com/shafikshaon/shortlink/service/github"
	googleSrv "github.com/shafikshaon/shortlink/service/google"
	"github.com/shafikshaon/shortlink/transport"
)

type OAuthHandler struct {
	googleService googleSrv.GoogleOAuthServiceI
	githubService githubSrv.GitHubOAuthServiceI
}

func NewOAuthHandler(googleService googleSrv.GoogleOAuthServiceI, githubService githubSrv.GitHubOAuthServiceI) *OAuthHandler {
	return &OAuthHandler{
		googleService: googleService,
		githubService: githubService,
	}
}

// GoogleLogin initiates the Google OAuth flow
// GET /api/v1/auth/google
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "Google OAuth login initiated")

	// Get client info for binding
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Get optional redirect URI from query params (for frontend to handle post-login)
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = "/"
	}

	authURL, err := h.googleService.GetAuthorizationURL(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate Google auth URL",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeOAuthDisabled, constant.ErrMsgOAuthDisabled)
		return
	}

	// Redirect to Google
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles the OAuth callback from Google
// GET /api/v1/auth/google/callback
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "Google OAuth callback received")

	// Check for error from Google
	if errParam := c.Query("error"); errParam != "" {
		errDesc := c.Query("error_description")
		logger.WarnCtx(ctx, "Google OAuth error",
			zap.String("error", errParam),
			zap.String("description", errDesc))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeOAuthProviderError, errDesc)
		return
	}

	// Get authorization code and state
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		logger.WarnCtx(ctx, "Missing code or state in callback",
			zap.Bool("has_code", code != ""),
			zap.Bool("has_state", state != ""))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeOAuthInvalidState, "Missing authorization code or state")
		return
	}

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Process the callback
	loginResp, svcErr := h.googleService.HandleCallback(ctx, state, code, ipAddress, userAgent)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Google OAuth login successful")
	transport.RespondWithSuccess(c, http.StatusOK, "Login successful", loginResp)
}

// GitHubLogin initiates the GitHub OAuth flow
// GET /api/v1/auth/github
func (h *OAuthHandler) GitHubLogin(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "GitHub OAuth login initiated")

	// Get client info for binding
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Get optional redirect URI from query params (for frontend to handle post-login)
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = "/"
	}

	authURL, err := h.githubService.GetAuthorizationURL(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate GitHub auth URL",
			zap.Error(err))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeGitHubOAuthDisabled, constant.ErrMsgGitHubOAuthDisabled)
		return
	}

	// Redirect to GitHub
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback handles the OAuth callback from GitHub
// GET /api/v1/auth/github/callback
func (h *OAuthHandler) GitHubCallback(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "GitHub OAuth callback received")

	// Check for error from GitHub
	if errParam := c.Query("error"); errParam != "" {
		errDesc := c.Query("error_description")
		logger.WarnCtx(ctx, "GitHub OAuth error",
			zap.String("error", errParam),
			zap.String("description", errDesc))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeOAuthProviderError, errDesc)
		return
	}

	// Get authorization code and state
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		logger.WarnCtx(ctx, "Missing code or state in callback",
			zap.Bool("has_code", code != ""),
			zap.Bool("has_state", state != ""))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeOAuthInvalidState, "Missing authorization code or state")
		return
	}

	// Get client info
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Process the callback
	loginResp, svcErr := h.githubService.HandleCallback(ctx, state, code, ipAddress, userAgent)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "GitHub OAuth login successful")
	transport.RespondWithSuccess(c, http.StatusOK, "Login successful", loginResp)
}

