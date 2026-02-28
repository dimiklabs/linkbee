package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/model"
)

type TotpBackupCodeRepositoryI interface {
	CreateBatch(ctx context.Context, codes []*model.TotpBackupCode) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TotpBackupCode, error)
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	MarkUsed(ctx context.Context, id uuid.UUID) error
}

type TotpBackupCodeRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewTotpBackupCodeRepository(masterDB, replicaDB *gorm.DB) TotpBackupCodeRepositoryI {
	return &TotpBackupCodeRepository{masterDB: masterDB, replicaDB: replicaDB}
}

func (r *TotpBackupCodeRepository) CreateBatch(ctx context.Context, codes []*model.TotpBackupCode) error {
	return r.masterDB.WithContext(ctx).Create(codes).Error
}

func (r *TotpBackupCodeRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TotpBackupCode, error) {
	var codes []*model.TotpBackupCode
	err := r.replicaDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&codes).Error
	return codes, err
}

func (r *TotpBackupCodeRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.masterDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.TotpBackupCode{}).Error
}

func (r *TotpBackupCodeRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	return r.masterDB.WithContext(ctx).
		Model(&model.TotpBackupCode{}).
		Where("id = ?", id).
		Update("used", true).Error
}
