package session

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/response"
)

type SessionNotificationServiceI interface {
	// NotifyNewSession sends notifications about a new session to the user
	NotifyNewSession(ctx context.Context, user *model.User, newSession *model.Session) error
	// CheckAndNotifyConcurrentSessions checks for concurrent sessions and returns alert if needed
	CheckAndNotifyConcurrentSessions(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) (*response.ConcurrentSessionsAlert, error)
	// ShouldNotifyNewSession determines if notification should be sent for this session
	ShouldNotifyNewSession(ctx context.Context, newSession *model.Session, existingSessions []model.Session) bool
}

type EmailNotifierI interface {
	SendNewSessionNotification(ctx context.Context, user *model.User, session *model.Session) error
}

type SessionNotificationService struct {
	sessionRepo   repository.SessionRepositoryI
	emailNotifier EmailNotifierI
	cfg           *config.SessionConfig
}

func NewSessionNotificationService(
	sessionRepo repository.SessionRepositoryI,
	emailNotifier EmailNotifierI,
	cfg *config.SessionConfig,
) SessionNotificationServiceI {
	return &SessionNotificationService{
		sessionRepo:   sessionRepo,
		emailNotifier: emailNotifier,
		cfg:           cfg,
	}
}

func (s *SessionNotificationService) NotifyNewSession(ctx context.Context, user *model.User, newSession *model.Session) error {
	if !s.cfg.NotifyOnNewSession {
		return nil
	}

	logger.InfoCtx(ctx, "Sending new session notification",
		zap.String("user_id", user.ID.String()),
		zap.String("session_id", newSession.ID.String()))

	// Send email notification if enabled
	if s.cfg.NotifyOnNewSessionEmail && s.emailNotifier != nil {
		if err := s.emailNotifier.SendNewSessionNotification(ctx, user, newSession); err != nil {
			logger.WarnCtx(ctx, "Failed to send new session email notification",
				zap.String("user_id", user.ID.String()),
				zap.Error(err))
			// Don't fail the login if email notification fails
		}
	}

	// Mark session as notified
	if err := s.sessionRepo.MarkSessionNotified(ctx, newSession.ID); err != nil {
		logger.WarnCtx(ctx, "Failed to mark session as notified",
			zap.String("session_id", newSession.ID.String()),
			zap.Error(err))
	}

	return nil
}

func (s *SessionNotificationService) CheckAndNotifyConcurrentSessions(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) (*response.ConcurrentSessionsAlert, error) {
	if !s.cfg.NotifyOnNewSession {
		return nil, nil
	}

	// Get other active sessions
	otherSessions, err := s.sessionRepo.GetOtherActiveSessions(ctx, userID, currentSessionID)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to get other active sessions for notification",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}

	if len(otherSessions) == 0 {
		// No other sessions, no alert needed
		return nil, nil
	}

	// Get current session
	currentSession, err := s.sessionRepo.GetByID(ctx, currentSessionID)
	if err != nil {
		logger.WarnCtx(ctx, "Failed to get current session for notification",
			zap.String("session_id", currentSessionID.String()),
			zap.Error(err))
		return nil, err
	}

	// Build alert response
	existingSessionNotifications := make([]response.NewSessionNotification, 0, len(otherSessions))
	for _, sess := range otherSessions {
		existingSessionNotifications = append(existingSessionNotifications, response.NewSessionNotification{
			SessionID:   sess.ID.String(),
			DeviceName:  sess.DeviceName,
			DeviceType:  sess.DeviceType,
			Browser:     sess.Browser,
			OS:          sess.OS,
			IPAddress:   sess.IPAddress,
			Location:    sess.Location,
			LoginMethod: sess.LoginMethod,
			CreatedAt:   sess.CreatedAt,
		})
	}

	alert := &response.ConcurrentSessionsAlert{
		Message: "You have logged in from a new device. Other active sessions were detected.",
		NewSession: response.NewSessionNotification{
			SessionID:   currentSession.ID.String(),
			DeviceName:  currentSession.DeviceName,
			DeviceType:  currentSession.DeviceType,
			Browser:     currentSession.Browser,
			OS:          currentSession.OS,
			IPAddress:   currentSession.IPAddress,
			Location:    currentSession.Location,
			LoginMethod: currentSession.LoginMethod,
			CreatedAt:   currentSession.CreatedAt,
		},
		ExistingSessions: existingSessionNotifications,
		TotalSessions:    len(otherSessions) + 1,
	}

	logger.InfoCtx(ctx, "Concurrent sessions alert generated",
		zap.String("user_id", userID.String()),
		zap.Int("total_sessions", alert.TotalSessions))

	return alert, nil
}

func (s *SessionNotificationService) ShouldNotifyNewSession(ctx context.Context, newSession *model.Session, existingSessions []model.Session) bool {
	if !s.cfg.NotifyOnNewSession {
		return false
	}

	// Always notify if suspicious activity detection is enabled
	if s.cfg.NotifyOnSuspiciousActivity {
		// Check if this is a new location or device
		for _, existing := range existingSessions {
			// Different IP address could indicate new location
			if existing.IPAddress != newSession.IPAddress {
				logger.InfoCtx(ctx, "New session from different IP detected",
					zap.String("session_id", newSession.ID.String()),
					zap.String("existing_ip", existing.IPAddress),
					zap.String("new_ip", newSession.IPAddress))
				return true
			}

			// Different device type or browser
			if existing.DeviceType != newSession.DeviceType || existing.Browser != newSession.Browser {
				logger.InfoCtx(ctx, "New session from different device detected",
					zap.String("session_id", newSession.ID.String()))
				return true
			}
		}
	}

	// Notify if there are other existing sessions
	return len(existingSessions) > 0
}
