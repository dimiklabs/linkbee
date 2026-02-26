package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware(allowedOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowedOrigins, ",")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		for _, allowed := range origins {
			if strings.TrimSpace(allowed) == origin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Refresh-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
