package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
)

type SessionRepositoryI interface {
	// Create
	Create(ctx context.Context, session *model.Session) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*model.Session, error)
	GetByRefreshTokenJTI(ctx context.Context, jti string) (*model.Session, error)
	GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]model.Session, error)
	CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetOldestActiveByUserID(ctx context.Context, userID uuid.UUID) (*model.Session, error)
	GetOtherActiveSessions(ctx context.Context, userID uuid.UUID, excludeSessionID uuid.UUID) ([]model.Session, error)

	// Update
	UpdateLastActivity(ctx context.Context, id uuid.UUID) error
	UpdateActivityWithDetails(ctx context.Context, id uuid.UUID, ipAddress string) error
	UpdateRefreshTokenJTI(ctx context.Context, id uuid.UUID, newJTI string, newExpiresAt time.Time) error
	MarkSessionNotified(ctx context.Context, id uuid.UUID) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByRefreshTokenJTI(ctx context.Context, jti string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type SessionRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewSessionRepository(masterDB, replicaDB *gorm.DB) SessionRepositoryI {
	return &SessionRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

// --------------- Create ---------------

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	logger.DebugCtx(ctx, "Creating session",
		zap.String("user_id", session.UserID.String()))

	if err := r.masterDB.WithContext(ctx).Create(session).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create session",
			zap.String("user_id", session.UserID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Session created",
		zap.String("session_id", session.ID.String()),
		zap.String("user_id", session.UserID.String()))
	return nil
}

// --------------- Read ---------------

func (r *SessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Session, error) {
	logger.DebugCtx(ctx, "Fetching session by ID",
		zap.String("session_id", id.String()))

	var session model.Session
	if err := r.replicaDB.WithContext(ctx).
		Where("id = ? AND expires_at > ?", id, time.Now()).
		First(&session).Error; err != nil {
		logger.DebugCtx(ctx, "Session not found or expired",
			zap.String("session_id", id.String()),
			zap.Error(err))
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetByRefreshTokenJTI(ctx context.Context, jti string) (*model.Session, error) {
	logger.DebugCtx(ctx, "Fetching session by refresh token JTI")

	var session model.Session
	if err := r.replicaDB.WithContext(ctx).
		Where("refresh_token_jti = ? AND expires_at > ?", jti, time.Now()).
		First(&session).Error; err != nil {
		logger.DebugCtx(ctx, "Session not found for JTI",
			zap.Error(err))
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) ([]model.Session, error) {
	logger.DebugCtx(ctx, "Fetching active sessions for user",
		zap.String("user_id", userID.String()))

	var sessions []model.Session
	if err := r.replicaDB.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("last_activity_at DESC").
		Find(&sessions).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch active sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) CountActiveByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	logger.DebugCtx(ctx, "Counting active sessions for user",
		zap.String("user_id", userID.String()))

	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.Session{}).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count active sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *SessionRepository) GetOldestActiveByUserID(ctx context.Context, userID uuid.UUID) (*model.Session, error) {
	logger.DebugCtx(ctx, "Fetching oldest active session for user",
		zap.String("user_id", userID.String()))

	var session model.Session
	if err := r.replicaDB.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Order("created_at ASC").
		First(&session).Error; err != nil {
		logger.DebugCtx(ctx, "No active sessions found for user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) GetOtherActiveSessions(ctx context.Context, userID uuid.UUID, excludeSessionID uuid.UUID) ([]model.Session, error) {
	logger.DebugCtx(ctx, "Fetching other active sessions for user",
		zap.String("user_id", userID.String()),
		zap.String("exclude_session_id", excludeSessionID.String()))

	var sessions []model.Session
	if err := r.replicaDB.WithContext(ctx).
		Where("user_id = ? AND expires_at > ? AND id != ?", userID, time.Now(), excludeSessionID).
		Order("last_activity_at DESC").
		Find(&sessions).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch other active sessions",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}
	return sessions, nil
}

// --------------- Update ---------------

func (r *SessionRepository) UpdateLastActivity(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Updating session last activity",
		zap.String("session_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.Session{}).
		Where("id = ?", id).
		Update("last_activity_at", time.Now()).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update session last activity",
			zap.String("session_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *SessionRepository) UpdateActivityWithDetails(ctx context.Context, id uuid.UUID, ipAddress string) error {
	logger.DebugCtx(ctx, "Updating session activity with details",
		zap.String("session_id", id.String()))

	updates := map[string]interface{}{
		"last_activity_at": time.Now(),
		"activity_count":   gorm.Expr("activity_count + 1"),
	}
	if ipAddress != "" {
		updates["last_activity_ip"] = ipAddress
	}

	if err := r.masterDB.WithContext(ctx).
		Model(&model.Session{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update session activity with details",
			zap.String("session_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *SessionRepository) UpdateRefreshTokenJTI(ctx context.Context, id uuid.UUID, newJTI string, newExpiresAt time.Time) error {
	logger.DebugCtx(ctx, "Updating session refresh token JTI",
		zap.String("session_id", id.String()))

	updates := map[string]interface{}{
		"refresh_token_jti": newJTI,
		"expires_at":        newExpiresAt,
		"last_activity_at":  time.Now(),
	}

	if err := r.masterDB.WithContext(ctx).
		Model(&model.Session{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update session refresh token JTI",
			zap.String("session_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *SessionRepository) MarkSessionNotified(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Marking session as notified",
		zap.String("session_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.Session{}).
		Where("id = ?", id).
		Update("notified_at", new(time.Now())).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to mark session as notified",
			zap.String("session_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

// --------------- Delete ---------------

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting session",
		zap.String("session_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Session{}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to delete session",
			zap.String("session_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Session deleted",
		zap.String("session_id", id.String()))
	return nil
}

func (r *SessionRepository) DeleteByRefreshTokenJTI(ctx context.Context, jti string) error {
	logger.DebugCtx(ctx, "Deleting session by refresh token JTI")

	if err := r.masterDB.WithContext(ctx).
		Where("refresh_token_jti = ?", jti).
		Delete(&model.Session{}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to delete session by JTI",
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Session deleted by JTI")
	return nil
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting all sessions for user",
		zap.String("user_id", userID.String()))

	result := r.masterDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Session{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete sessions for user",
			zap.String("user_id", userID.String()),
			zap.Error(result.Error))
		return result.Error
	}

	logger.InfoCtx(ctx, "Sessions deleted for user",
		zap.String("user_id", userID.String()),
		zap.Int64("count", result.RowsAffected))
	return nil
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	logger.DebugCtx(ctx, "Deleting expired sessions")

	result := r.masterDB.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&model.Session{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete expired sessions",
			zap.Error(result.Error))
		return 0, result.Error
	}

	logger.InfoCtx(ctx, "Expired sessions deleted",
		zap.Int64("count", result.RowsAffected))
	return result.RowsAffected, nil
}
