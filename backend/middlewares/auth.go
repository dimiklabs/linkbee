package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/util"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyEmail  = "email"
	ContextKeyRole   = "role"
	ContextKeyJTI    = "jti"
)

func AuthMiddleware(cfg *config.AppConfig, tokenBlacklistRepo repository.TokenBlacklistRepositoryI) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing authorization header",
				zap.String("request_id", requestID),
				zap.String("path", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_code":  "UNAUTHORIZED",
				"description": "Authorization header required",
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

		tokenStr := parts[1]

		claims, err := util.ValidateAccessToken(tokenStr, cfg.JWTSecret, cfg.JWTIssuer)
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
