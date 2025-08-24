package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/moriverse/45-server/internal/domain/user"
)

const (
	lastActiveCacheKeyPrefix = "last-active"
	lastActiveCacheTTL       = 5 * time.Minute
)

// Service is the application service for user-related operations.
type Service struct {
	userRepo    user.Repository
	redisClient *redis.Client
	logger      *slog.Logger
}

// NewService creates a new instance of the user service.
func NewService(
	userRepo user.Repository,
	redisClient *redis.Client,
	logger *slog.Logger,
) *Service {
	return &Service{
		userRepo:    userRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// UpdateLastActive updates a user's last active time if the configured TTL has passed.
// It uses Redis for caching to avoid hitting the database on every request.
func (s *Service) UpdateLastActive(ctx context.Context, userID user.UserID) {
	key := fmt.Sprintf("%s:%s", lastActiveCacheKeyPrefix, userID)

	// SetNX returns true if the key was set, false if it already existed.
	wasSet, err := s.redisClient.SetNX(ctx, key, "active", lastActiveCacheTTL).Result()
	if err != nil {
		s.logger.Error(
			"Failed to set last active cache key",
			"user_id", userID,
			"error", err,
		)
		return
	}

	// If the key was set, it means this is the first request in the TTL window,
	// so we should update the database.
	if wasSet {
		go func() {
			// We use a background context because the original request context might be cancelled.
			if err := s.userRepo.UpdateLastActiveAt(
				context.Background(),
				userID,
				time.Now(),
			); err != nil {
				s.logger.Error(
					"Failed to update last active time in database",
					"user_id", userID,
					"error", err,
				)
			}
		}()
	}
}
