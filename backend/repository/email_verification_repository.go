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

type EmailVerificationRepositoryI interface {
	Create(ctx context.Context, verification *model.EmailVerification) error
	GetByToken(ctx context.Context, token string) (*model.EmailVerification, error)
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type EmailVerificationRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewEmailVerificationRepository(masterDB, replicaDB *gorm.DB) EmailVerificationRepositoryI {
	return &EmailVerificationRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

func (r *EmailVerificationRepository) Create(ctx context.Context, verification *model.EmailVerification) error {
	logger.DebugCtx(ctx, "Creating email verification token",
		zap.String("user_id", verification.UserID.String()))

	if err := r.masterDB.WithContext(ctx).Create(verification).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create email verification token",
			zap.String("user_id", verification.UserID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Email verification token created",
		zap.String("user_id", verification.UserID.String()))
	return nil
}

func (r *EmailVerificationRepository) GetByToken(ctx context.Context, token string) (*model.EmailVerification, error) {
	logger.DebugCtx(ctx, "Fetching email verification by token")

	var verification model.EmailVerification
	if err := r.replicaDB.WithContext(ctx).
		Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).
		First(&verification).Error; err != nil {
		logger.DebugCtx(ctx, "Email verification token not found or expired",
			zap.Error(err))
		return nil, err
	}
	return &verification, nil
}

func (r *EmailVerificationRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Marking email verification token as used",
		zap.String("id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.EmailVerification{}).
		Where("id = ?", id).
		Update("used", true).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to mark email verification token as used",
			zap.String("id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Email verification token marked as used",
		zap.String("id", id.String()))
	return nil
}

func (r *EmailVerificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting email verification tokens for user",
		zap.String("user_id", userID.String()))

	if err := r.masterDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.EmailVerification{}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to delete email verification tokens",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Email verification tokens deleted for user",
		zap.String("user_id", userID.String()))
	return nil
}

func (r *EmailVerificationRepository) DeleteExpired(ctx context.Context) error {
	logger.DebugCtx(ctx, "Deleting expired email verification tokens")

	result := r.masterDB.WithContext(ctx).
		Where("expires_at < ? OR used = ?", time.Now(), true).
		Delete(&model.EmailVerification{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete expired email verification tokens",
			zap.Error(result.Error))
		return result.Error
	}

	logger.InfoCtx(ctx, "Expired email verification tokens deleted",
		zap.Int64("count", result.RowsAffected))
	return nil
}
