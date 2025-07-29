package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/pkg/auth"
)

// AuthMiddleware creates a gin middleware for JWT authentication.
func AuthMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header is missing"},
			)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header format must be Bearer {token}"},
			)
			return
		}

		tokenString := headerParts[1]

		claims, err := auth.ValidateToken(tokenString, cfg.SecretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set the user ID in the context for later use.
		c.Set("userID", claims.Subject)

		c.Next()
	}
}
