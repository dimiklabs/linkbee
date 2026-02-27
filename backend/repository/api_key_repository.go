package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type APIKeyRepositoryI interface {
	Create(ctx context.Context, key *model.APIKey) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetByPrefix(ctx context.Context, prefix string) (*model.APIKey, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}

type apiKeyRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewAPIKeyRepository(master, replica *gorm.DB) APIKeyRepositoryI {
	return &apiKeyRepository{master: master, replica: replica}
}

func (r *apiKeyRepository) Create(ctx context.Context, key *model.APIKey) error {
	return r.master.WithContext(ctx).Create(key).Error
}

func (r *apiKeyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
	var keys []model.APIKey
	err := r.replica.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&keys).Error
	return keys, err
}

func (r *apiKeyRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.replica.WithContext(ctx).Model(&model.APIKey{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *apiKeyRepository) GetByPrefix(ctx context.Context, prefix string) (*model.APIKey, error) {
	var key model.APIKey
	err := r.replica.WithContext(ctx).
		Where("key_prefix = ?", prefix).
		First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *apiKeyRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.master.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.APIKey{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.master.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("id = ?", id).
		Update("last_used_at", now).Error
}
