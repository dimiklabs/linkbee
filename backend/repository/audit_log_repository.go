package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/model"
)

type AuditLogRepositoryI interface {
	Create(ctx context.Context, log *model.AuditLog) error
	ListByUserID(ctx context.Context, userID uuid.UUID, page, limit int, action, resourceType string, from, to *time.Time) ([]model.AuditLog, int64, error)
	ListAllByUserID(ctx context.Context, userID uuid.UUID, action string, from, to *time.Time, limit int) ([]model.AuditLog, error)
}

type auditLogRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewAuditLogRepository(master, replica *gorm.DB) AuditLogRepositoryI {
	return &auditLogRepository{master: master, replica: replica}
}

func (r *auditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.master.WithContext(ctx).Create(log).Error
}

func (r *auditLogRepository) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int,
	action, resourceType string,
	from, to *time.Time,
) ([]model.AuditLog, int64, error) {
	offset := (page - 1) * limit

	q := r.replica.WithContext(ctx).
		Model(&model.AuditLog{}).
		Where("user_id = ?", userID)

	if action != "" {
		q = q.Where("action = ?", action)
	}
	if resourceType != "" {
		q = q.Where("resource_type = ?", resourceType)
	}
	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []model.AuditLog
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *auditLogRepository) ListAllByUserID(
	ctx context.Context,
	userID uuid.UUID,
	action string,
	from, to *time.Time,
	limit int,
) ([]model.AuditLog, error) {
	q := r.replica.WithContext(ctx).
		Model(&model.AuditLog{}).
		Where("user_id = ?", userID)

	if action != "" {
		q = q.Where("action = ?", action)
	}
	if from != nil {
		q = q.Where("created_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("created_at <= ?", *to)
	}

	var logs []model.AuditLog
	if err := q.Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
