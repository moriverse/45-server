package context

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/moriverse/45-server/internal/domain/user"
)

func SetUserID(c *gin.Context, userID user.UserID) {
	c.Set(string(UserIDKey), userID)
}

func GetUserID(c *gin.Context) (user.UserID, bool) {
	val, exists := c.Get(string(UserIDKey))
	if !exists {
		return "", false
	}
	userID, ok := val.(user.UserID)
	return userID, ok
}

func SetLogger(c *gin.Context, logger *slog.Logger) {
	c.Set(string(LoggerKey), logger)
}

func GetLogger(c *gin.Context) (*slog.Logger, bool) {
	val, exists := c.Get(string(LoggerKey))
	if !exists {
		return nil, false
	}
	logger, ok := val.(*slog.Logger)
	return logger, ok
}
