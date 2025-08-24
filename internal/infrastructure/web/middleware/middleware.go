package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appUser "github.com/moriverse/45-server/internal/app/user"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/utils"
)

const (
	LoggerKey = "logger"
)

// Middleware encapsulates all middleware logic and dependencies.
type Middleware struct {
	userService *appUser.Service
	jwtConfig   config.JWTConfig
	logger      *slog.Logger
}

// NewMiddleware creates a new Middleware instance.
func NewMiddleware(
	userService *appUser.Service,
	jwtConfig config.JWTConfig,
	logger *slog.Logger,
) *Middleware {
	return &Middleware{
		userService: userService,
		jwtConfig:   jwtConfig,
		logger:      logger,
	}
}

// LoggingMiddleware creates a request-specific logger with a request_id
// and stores it in the context.
func (m *Middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		requestLogger := m.logger.With("request_id", requestID)

		c.Set(LoggerKey, requestLogger)

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		requestLogger.Info(
			"Request handled",
			"status_code", statusCode,
			"latency", latency,
			"client_ip", clientIP,
			"method", method,
			"path", path,
			"body_size", bodySize,
			"error_message", errorMessage,
		)
	}
}

// AuthMiddleware is a Gin middleware for JWT authentication.
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header is missing"},
			)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header format is Bearer {token}"},
			)
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, m.jwtConfig.SecretKey)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		// Set user ID in context for downstream handlers
		c.Set("userID", claims.Subject)

		// Update user's last active time
		userID := user.UserID(claims.Subject)
		m.userService.UpdateLastActive(c.Request.Context(), userID)

		c.Next()
	}
}
