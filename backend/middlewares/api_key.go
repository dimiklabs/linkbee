package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/repository"
	apiKeySvc "github.com/shafikshaon/shortlink/service/apikey"
	"github.com/shafikshaon/shortlink/util"
)

// AuthOrAPIKeyMiddleware authenticates via JWT Bearer token OR X-API-Key header.
// If the X-API-Key header is present it takes precedence over the Authorization header.
func AuthOrAPIKeyMiddleware(
	cfg *config.AppConfig,
	tokenBlacklistRepo repository.TokenBlacklistRepositoryI,
	apiKeyService apiKeySvc.APIKeyServiceI,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c)

		// ── API key path ──────────────────────────────────────────────────────
		if apiKeyHeader := c.GetHeader("X-API-Key"); apiKeyHeader != "" {
			validated, svcErr := apiKeyService.Validate(c.Request.Context(), apiKeyHeader)
			if svcErr != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error_code":  "UNAUTHORIZED",
					"description": svcErr.Description,
					"request_id":  requestID,
				})
				return
			}
			c.Set(ContextKeyUserID, validated.UserID.String())
			c.Next()
			return
		}

		// ── Query-param token path (for <img src> endpoints like /qr) ─────────
		if tokenParam := c.Query("token"); tokenParam != "" {
			claims, err := util.ValidateAccessToken(tokenParam, cfg.JWTSecret, cfg.JWTIssuer)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error_code":  "UNAUTHORIZED",
					"description": "Invalid or expired token",
					"request_id":  requestID,
				})
				return
			}
			if tokenBlacklistRepo != nil {
				blacklisted, bErr := tokenBlacklistRepo.IsBlacklisted(c.Request.Context(), claims.ID)
				if bErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error_code":  "INTERNAL_SERVER_ERROR",
						"description": "Failed to validate token",
						"request_id":  requestID,
					})
					return
				}
				if blacklisted {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error_code":  "UNAUTHORIZED",
						"description": "Token has been revoked",
						"request_id":  requestID,
					})
					return
				}
			}
			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyEmail, claims.Email)
			c.Set(ContextKeyRole, claims.Role)
			c.Set(ContextKeyJTI, claims.ID)
			c.Next()
			return
		}

		// ── JWT Bearer path ───────────────────────────────────────────────────
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header",
				zap.String("request_id", requestID),
				zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":  "UNAUTHORIZED",
				"description": "Authorization header or X-API-Key required",
				"request_id":  requestID,
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":  "UNAUTHORIZED",
				"description": "Invalid authorization header format. Expected: Bearer <token>",
				"request_id":  requestID,
			})
			return
		}

		claims, err := util.ValidateAccessToken(parts[1], cfg.JWTSecret, cfg.JWTIssuer)
		if err != nil {
			logger.Warn("Token validation failed",
				zap.String("request_id", requestID),
				zap.String("path", c.Request.URL.Path),
				zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":  "UNAUTHORIZED",
				"description": "Invalid or expired token",
				"request_id":  requestID,
			})
			return
		}

		if tokenBlacklistRepo != nil {
			blacklisted, err := tokenBlacklistRepo.IsBlacklisted(c.Request.Context(), claims.ID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error_code":  "INTERNAL_SERVER_ERROR",
					"description": "Failed to validate token",
					"request_id":  requestID,
				})
				return
			}
			if blacklisted {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error_code":  "UNAUTHORIZED",
					"description": "Token has been revoked",
					"request_id":  requestID,
				})
				return
			}
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Set(ContextKeyRole, claims.Role)
		c.Set(ContextKeyJTI, claims.ID)

		c.Next()
	}
}
