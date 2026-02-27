package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type WebhookDeliveryRepositoryI interface {
	Create(ctx context.Context, delivery *model.WebhookDelivery) error
	ListByWebhookID(ctx context.Context, webhookID, userID uuid.UUID, page, limit int) ([]*model.WebhookDelivery, int64, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*model.WebhookDelivery, error)
}

type webhookDeliveryRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewWebhookDeliveryRepository(master, replica *gorm.DB) WebhookDeliveryRepositoryI {
	return &webhookDeliveryRepository{master: master, replica: replica}
}

func (r *webhookDeliveryRepository) Create(ctx context.Context, delivery *model.WebhookDelivery) error {
	return r.master.WithContext(ctx).Create(delivery).Error
}

func (r *webhookDeliveryRepository) ListByWebhookID(ctx context.Context, webhookID, userID uuid.UUID, page, limit int) ([]*model.WebhookDelivery, int64, error) {
	var deliveries []*model.WebhookDelivery
	var total int64

	q := r.replica.WithContext(ctx).Model(&model.WebhookDelivery{}).
		Where("webhook_id = ? AND user_id = ?", webhookID, userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&deliveries).Error
	return deliveries, total, err
}

func (r *webhookDeliveryRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.WebhookDelivery, error) {
	var delivery model.WebhookDelivery
	err := r.replica.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&delivery).Error
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}
