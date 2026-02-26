package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware aborts with 403 if the authenticated user is not an admin.
// Must be used after AuthOrAPIKeyMiddleware (which sets ContextKeyRole).
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error_code":  "FORBIDDEN",
				"description": "Admin access required",
			})
			return
		}
		c.Next()
	}
}
