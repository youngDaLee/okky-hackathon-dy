package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TokenValidator is implemented by AuthService.ValidateAccessToken.
type TokenValidator interface {
	ValidateAccessToken(tokenStr string) (string, error)
}

// Auth returns a JWT validation middleware.
// Pass a TokenValidator to enable real JWT validation.
// Without a validator it falls back to treating the token value as userID (dev stub).
func Auth(validator ...TokenValidator) gin.HandlerFunc {
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

		if len(validator) > 0 && validator[0] != nil {
			userID, err := validator[0].ValidateAccessToken(token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
				return
			}
			c.Set("userID", userID)
		} else {
			c.Set("userID", token)
		}
		c.Next()
	}
}
