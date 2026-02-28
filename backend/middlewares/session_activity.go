package middlewares

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/repository"
)

// sessionActivityTracker tracks when sessions were last updated to avoid too frequent updates
type sessionActivityTracker struct {
	mu             sync.RWMutex
	lastUpdateTime map[uuid.UUID]time.Time
	updateInterval time.Duration
}

func newSessionActivityTracker(updateIntervalMins int) *sessionActivityTracker {
	return &sessionActivityTracker{
		lastUpdateTime: make(map[uuid.UUID]time.Time),
		updateInterval: time.Duration(updateIntervalMins) * time.Minute,
	}
}

func (t *sessionActivityTracker) shouldUpdate(sessionID uuid.UUID) bool {
	t.mu.RLock()
	lastUpdate, exists := t.lastUpdateTime[sessionID]
	t.mu.RUnlock()

	if !exists {
		return true
	}
	return time.Since(lastUpdate) >= t.updateInterval
}

func (t *sessionActivityTracker) markUpdated(sessionID uuid.UUID) {
	t.mu.Lock()
	t.lastUpdateTime[sessionID] = time.Now()
	t.mu.Unlock()
}

// cleanup removes entries for sessions that haven't been seen in a while
func (t *sessionActivityTracker) cleanup() {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := time.Now().Add(-24 * time.Hour)
	for sessionID, lastUpdate := range t.lastUpdateTime {
		if lastUpdate.Before(cutoff) {
			delete(t.lastUpdateTime, sessionID)
		}
	}
}

var (
	activityTracker     *sessionActivityTracker
	activityTrackerOnce sync.Once
)

// SessionActivityMiddleware creates a middleware that tracks session activity
func SessionActivityMiddleware(sessionRepo repository.SessionRepositoryI, cfg *config.SessionConfig) gin.HandlerFunc {
	// Only track activity if enabled
	if cfg == nil || !cfg.TrackActivity {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Initialize tracker once
	activityTrackerOnce.Do(func() {
		activityTracker = newSessionActivityTracker(cfg.ActivityUpdateIntervalMins)

		// Start cleanup goroutine
		go func() {
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				activityTracker.cleanup()
			}
		}()
	})

	// Cache for JTI to session ID mapping to avoid repeated lookups
	jtiToSessionID := make(map[string]uuid.UUID)
	var cacheMu sync.RWMutex

	return func(c *gin.Context) {
		// Get JTI from context (set by auth middleware)
		jti, exists := c.Get(ContextKeyJTI)
		if !exists {
			c.Next()
			return
		}

		jtiStr, ok := jti.(string)
		if !ok || jtiStr == "" {
			c.Next()
			return
		}

		// Try to get session ID from cache first
		cacheMu.RLock()
		sessionID, cached := jtiToSessionID[jtiStr]
		cacheMu.RUnlock()

		if !cached {
			// Look up session by refresh token JTI
			// Note: We're using access token JTI here, which is different from refresh token JTI
			// The auth middleware uses access token JTI, so we need to look up session differently
			// For now, skip caching and let the handler deal with session activity
			c.Next()
			return
		}

		// Check if we should update (rate-limited to avoid too many DB writes)
		if activityTracker.shouldUpdate(sessionID) {
			// Get IP address
			ipAddress := c.ClientIP()

			// Update activity asynchronously to not block the request
			go func(sessID uuid.UUID, ip string) {
				ctx := c.Request.Context()
				if err := sessionRepo.UpdateActivityWithDetails(ctx, sessID, ip); err != nil {
					logger.WarnCtx(ctx, "Failed to update session activity",
						zap.String("session_id", sessID.String()),
						zap.Error(err))
				} else {
					activityTracker.markUpdated(sessID)
					logger.DebugCtx(ctx, "Session activity updated",
						zap.String("session_id", sessID.String()))
				}
			}(sessionID, ipAddress)
		}

		c.Next()
	}
}

// SetSessionIDInContext is a helper to set session ID in context after looking it up
// This can be called from handlers that already have the session
func SetSessionIDInContext(c *gin.Context, sessionID uuid.UUID) {
	c.Set("session_id", sessionID)
}
