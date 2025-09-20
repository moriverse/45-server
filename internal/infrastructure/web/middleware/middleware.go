package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/config"
	"github.com/moriverse/45-server/internal/infrastructure/web/context"
	"github.com/moriverse/45-server/internal/infrastructure/web/response"
	"github.com/moriverse/45-server/internal/utils"
)

type Middleware struct {
	jwtConfig config.JWTConfig
	logger    *slog.Logger
}

func NewMiddleware(jwtConfig config.JWTConfig, logger *slog.Logger) *Middleware {
	return &Middleware{
		jwtConfig: jwtConfig,
		logger:    logger,
	}
}

func (m *Middleware) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		requestLogger := m.logger.With(
			"method", c.Request.Method,
			"path", path,
			"latency", time.Since(start),
			"clientIP", c.ClientIP(),
			"statusCode", c.Writer.Status(),
		)
		if raw != "" {
			requestLogger = requestLogger.With("rawQuery", raw)
		}

		context.SetLogger(c, requestLogger)

		if c.Writer.Status() >= http.StatusInternalServerError {
			requestLogger.Error(
				"Request completed with server error",
				"errors",
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		} else if c.Writer.Status() >= http.StatusBadRequest {
			requestLogger.Warn(
				"Request completed with client error",
				"errors",
				c.Errors.ByType(gin.ErrorTypePrivate).String(),
			)
		} else {
			requestLogger.Info("Request completed successfully")
		}
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, response.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Authorization header is required.",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, response.APIError{
				Code:    "UNAUTHORIZED",
				Message: "Authorization header format must be Bearer {token}.",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, m.jwtConfig.SecretKey)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.APIError{
				Code:    "INVALID_TOKEN",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		context.SetUserID(c, user.UserID(claims.Subject))
		c.Next()
	}
}
