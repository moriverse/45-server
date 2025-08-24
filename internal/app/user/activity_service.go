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

// ActivityService is responsible for tracking user activity.
type ActivityService struct {
	redisClient *redis.Client
	userRepo    user.Repository
	logger      *slog.Logger
}

// NewActivityService creates a new instance of ActivityService.
func NewActivityService(
	redisClient *redis.Client,
	userRepo user.Repository,
	logger *slog.Logger,
) *ActivityService {
	return &ActivityService{
		redisClient: redisClient,
		userRepo:    userRepo,
		logger:      logger,
	}
}

// UpdateLastActive updates a user's last active time if the configured TTL has passed.
// It uses Redis for caching to avoid hitting the database on every request.
func (s *ActivityService) UpdateLastActive(ctx context.Context, userID user.UserID) {
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
			bgCtx := context.Background()
			now := time.Now()
			if err := s.userRepo.UpdateLastActiveAt(bgCtx, userID, now); err != nil {
				s.logger.Error(
					"Failed to update last active time in database",
					"user_id", userID,
					"error", err,
				)
			}
		}()
	}
}
