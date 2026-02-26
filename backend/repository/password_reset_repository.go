package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
)

type PasswordResetRepositoryI interface {
	Create(ctx context.Context, reset *model.PasswordReset) error
	GetByToken(ctx context.Context, token string) (*model.PasswordReset, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type PasswordResetRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewPasswordResetRepository(masterDB, replicaDB *gorm.DB) PasswordResetRepositoryI {
	return &PasswordResetRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

func (r *PasswordResetRepository) Create(ctx context.Context, reset *model.PasswordReset) error {
	logger.DebugCtx(ctx, "Creating password reset token",
		zap.String("user_id", reset.UserID.String()))

	if err := r.masterDB.WithContext(ctx).Create(reset).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create password reset token",
			zap.String("user_id", reset.UserID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Password reset token created",
		zap.String("user_id", reset.UserID.String()))
	return nil
}

func (r *PasswordResetRepository) GetByToken(ctx context.Context, token string) (*model.PasswordReset, error) {
	logger.DebugCtx(ctx, "Fetching password reset by token")

	var reset model.PasswordReset
	if err := r.replicaDB.WithContext(ctx).
		Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).
		First(&reset).Error; err != nil {
		logger.DebugCtx(ctx, "Password reset token not found or expired",
			zap.Error(err))
		return nil, err
	}
	return &reset, nil
}

func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Marking password reset token as used",
		zap.String("id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.PasswordReset{}).
		Where("id = ?", id).
		Update("used", true).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to mark password reset token as used",
			zap.String("id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Password reset token marked as used",
		zap.String("id", id.String()))
	return nil
}

func (r *PasswordResetRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting password reset tokens for user",
		zap.String("user_id", userID.String()))

	if err := r.masterDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.PasswordReset{}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to delete password reset tokens",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Password reset tokens deleted for user",
		zap.String("user_id", userID.String()))
	return nil
}

func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	logger.DebugCtx(ctx, "Deleting expired password reset tokens")

	result := r.masterDB.WithContext(ctx).
		Where("expires_at < ? OR used = ?", time.Now(), true).
		Delete(&model.PasswordReset{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete expired password reset tokens",
			zap.Error(result.Error))
		return result.Error
	}

	logger.InfoCtx(ctx, "Expired password reset tokens deleted",
		zap.Int64("count", result.RowsAffected))
	return nil
}
