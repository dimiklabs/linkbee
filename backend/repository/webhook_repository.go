package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type WebhookRepositoryI interface {
	Create(ctx context.Context, webhook *model.Webhook) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Webhook, error)
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Webhook, error)
	Update(ctx context.Context, webhook *model.Webhook) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	GetActiveByEvent(ctx context.Context, userID uuid.UUID, event string) ([]model.Webhook, error)
}

type webhookRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewWebhookRepository(master, replica *gorm.DB) WebhookRepositoryI {
	return &webhookRepository{master: master, replica: replica}
}

func (r *webhookRepository) Create(ctx context.Context, webhook *model.Webhook) error {
	return r.master.WithContext(ctx).Create(webhook).Error
}

func (r *webhookRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	err := r.replica.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&webhooks).Error
	return webhooks, err
}

func (r *webhookRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.replica.WithContext(ctx).Model(&model.Webhook{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *webhookRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Webhook, error) {
	var webhook model.Webhook
	err := r.replica.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&webhook).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (r *webhookRepository) Update(ctx context.Context, webhook *model.Webhook) error {
	return r.master.WithContext(ctx).Save(webhook).Error
}

func (r *webhookRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.master.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Webhook{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetActiveByEvent returns all active webhooks for a user that are subscribed to the given event.
// Uses PostgreSQL array overlap operator (@>) to check if the event is in the events array.
func (r *webhookRepository) GetActiveByEvent(ctx context.Context, userID uuid.UUID, event string) ([]model.Webhook, error) {
	var webhooks []model.Webhook
	err := r.replica.WithContext(ctx).
		Where("user_id = ? AND is_active = true AND events @> ARRAY[?]::text[]", userID, event).
		Find(&webhooks).Error
	return webhooks, err
}
