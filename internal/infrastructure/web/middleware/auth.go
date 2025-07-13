package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/pkg/auth"
)

// AuthMiddleware creates a gin middleware for JWT authentication.
func AuthMiddleware() gin.HandlerFunc {
	// In a real app, you'd load the config once at startup.
	// For simplicity, we'll load it here, but be aware of the performance implications.
	cfg, err := config.LoadConfig() // This might need adjustment based on your config loading strategy
	if err != nil {
		// A more robust solution would be to panic or log fatally
		// if the config can't be loaded at startup.
		panic("Could not load configuration for middleware")
	}

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

		claims, err := auth.ValidateToken(tokenString, cfg.JWT.SecretKey)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid token"},
			)
			return
		}

		// Set the user ID in the context for later use.
		c.Set("userID", claims.Subject)

		c.Next()
	}
}
