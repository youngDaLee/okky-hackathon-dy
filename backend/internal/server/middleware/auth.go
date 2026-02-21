package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth is a stub middleware that reads the Authorization header.
// Until JWT middleware is implemented (#2), it accepts any non-empty Bearer token
// and injects its value as "userID" into the context.
// In production this will be replaced with full JWT validation.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
			return
		}
		// TODO(#2): validate JWT and extract real user ID
		c.Set("userID", token)
		c.Next()
	}
}
