package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type CustomDomainRepositoryI interface {
	Create(ctx context.Context, d *model.CustomDomain) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.CustomDomain, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*model.CustomDomain, error)
	GetByDomain(ctx context.Context, domain string) (*model.CustomDomain, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}

type customDomainRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewCustomDomainRepository(master, replica *gorm.DB) CustomDomainRepositoryI {
	return &customDomainRepository{master: master, replica: replica}
}

func (r *customDomainRepository) Create(ctx context.Context, d *model.CustomDomain) error {
	return r.master.WithContext(ctx).Create(d).Error
}

func (r *customDomainRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.CustomDomain, error) {
	var domains []model.CustomDomain
	err := r.replica.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&domains).Error
	return domains, err
}

func (r *customDomainRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.CustomDomain, error) {
	var d model.CustomDomain
	err := r.replica.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *customDomainRepository) GetByDomain(ctx context.Context, domain string) (*model.CustomDomain, error) {
	var d model.CustomDomain
	err := r.replica.WithContext(ctx).
		Where("domain = ? AND status = ?", domain, model.DomainStatusVerified).
		First(&d).Error
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *customDomainRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.master.WithContext(ctx).
		Model(&model.CustomDomain{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *customDomainRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result := r.master.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.CustomDomain{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *customDomainRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.replica.WithContext(ctx).
		Model(&model.CustomDomain{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
