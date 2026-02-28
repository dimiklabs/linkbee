package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/response"
	githubSrv "github.com/shafikshaon/shortlink/service/github"
	googleSrv "github.com/shafikshaon/shortlink/service/google"
	"github.com/shafikshaon/shortlink/transport"
)

type OAuthHandler struct {
	googleService googleSrv.GoogleOAuthServiceI
	githubService githubSrv.GitHubOAuthServiceI
	frontendURL   string
}

func NewOAuthHandler(googleService googleSrv.GoogleOAuthServiceI, githubService githubSrv.GitHubOAuthServiceI, appCfg *config.AppConfig) *OAuthHandler {
	return &OAuthHandler{
		googleService: googleService,
		githubService: githubService,
		frontendURL:   appCfg.FrontendURL,
	}
}

// oauthErrorRedirect redirects the browser to the frontend login page with error details.
func (h *OAuthHandler) oauthErrorRedirect(c *gin.Context, errCode, description string) {
	q := url.Values{}
	q.Set("oauth_error", errCode)
	q.Set("oauth_message", description)
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?%s", h.frontendURL, q.Encode()))
}

// oauthSuccessRedirect redirects the browser to the frontend callback page with tokens.
func (h *OAuthHandler) oauthSuccessRedirect(c *gin.Context, login *response.LoginResponse) {
	q := url.Values{}
	q.Set("access_token", login.AccessToken)
	q.Set("refresh_token", login.RefreshToken)
	q.Set("token_type", login.TokenType)
	q.Set("expires_in", fmt.Sprintf("%d", login.ExpiresIn))
	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/auth/callback?%s", h.frontendURL, q.Encode()))
}

// GoogleLogin initiates the Google OAuth flow
// GET /api/v1/auth/google
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "Google OAuth login initiated")

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = "/"
	}

	authURL, err := h.googleService.GetAuthorizationURL(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate Google auth URL", zap.Error(err))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeOAuthDisabled, constant.ErrMsgOAuthDisabled)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles the OAuth callback from Google
// GET /api/v1/auth/google/callback
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "Google OAuth callback received")

	if errParam := c.Query("error"); errParam != "" {
		errDesc := c.Query("error_description")
		logger.WarnCtx(ctx, "Google OAuth error",
			zap.String("error", errParam),
			zap.String("description", errDesc))
		h.oauthErrorRedirect(c, constant.ErrCodeOAuthProviderError, errDesc)
		return
	}

	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		logger.WarnCtx(ctx, "Missing code or state in Google callback")
		h.oauthErrorRedirect(c, constant.ErrCodeOAuthInvalidState, "Missing authorization code or state")
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	loginResp, svcErr := h.googleService.HandleCallback(ctx, state, code, ipAddress, userAgent)
	if svcErr != nil {
		logger.WarnCtx(ctx, "Google OAuth callback failed",
			zap.String("error_code", svcErr.ErrorCode),
			zap.String("description", svcErr.Description))
		h.oauthErrorRedirect(c, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "Google OAuth login successful")
	h.oauthSuccessRedirect(c, loginResp)
}

// GitHubLogin initiates the GitHub OAuth flow
// GET /api/v1/auth/github
func (h *OAuthHandler) GitHubLogin(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "GitHub OAuth login initiated")

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = "/"
	}

	authURL, err := h.githubService.GetAuthorizationURL(ctx, ipAddress, userAgent, redirectURI)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate GitHub auth URL", zap.Error(err))
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeGitHubOAuthDisabled, constant.ErrMsgGitHubOAuthDisabled)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GitHubCallback handles the OAuth callback from GitHub
// GET /api/v1/auth/github/callback
func (h *OAuthHandler) GitHubCallback(c *gin.Context) {
	ctx := c.Request.Context()

	logger.InfoCtx(ctx, "GitHub OAuth callback received")

	if errParam := c.Query("error"); errParam != "" {
		errDesc := c.Query("error_description")
		logger.WarnCtx(ctx, "GitHub OAuth error",
			zap.String("error", errParam),
			zap.String("description", errDesc))
		h.oauthErrorRedirect(c, constant.ErrCodeOAuthProviderError, errDesc)
		return
	}

	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		logger.WarnCtx(ctx, "Missing code or state in GitHub callback")
		h.oauthErrorRedirect(c, constant.ErrCodeOAuthInvalidState, "Missing authorization code or state")
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	loginResp, svcErr := h.githubService.HandleCallback(ctx, state, code, ipAddress, userAgent)
	if svcErr != nil {
		logger.WarnCtx(ctx, "GitHub OAuth callback failed",
			zap.String("error_code", svcErr.ErrorCode),
			zap.String("description", svcErr.Description))
		h.oauthErrorRedirect(c, svcErr.ErrorCode, svcErr.Description)
		return
	}

	logger.InfoCtx(ctx, "GitHub OAuth login successful")
	h.oauthSuccessRedirect(c, loginResp)
}
