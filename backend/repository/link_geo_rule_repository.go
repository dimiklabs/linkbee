package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/model"
)

type LinkGeoRuleRepositoryI interface {
	GetByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.LinkGeoRule, error)
	GetByLinkIDAndCountry(ctx context.Context, linkID uuid.UUID, countryCode string) (*model.LinkGeoRule, error)
	Create(ctx context.Context, rule *model.LinkGeoRule) error
	Update(ctx context.Context, rule *model.LinkGeoRule) error
	Delete(ctx context.Context, id, linkID uuid.UUID) error
	DeleteByLinkID(ctx context.Context, linkID uuid.UUID) error
}

type linkGeoRuleRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewLinkGeoRuleRepository(master, replica *gorm.DB) LinkGeoRuleRepositoryI {
	return &linkGeoRuleRepository{master: master, replica: replica}
}

func (r *linkGeoRuleRepository) GetByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.LinkGeoRule, error) {
	var rules []model.LinkGeoRule
	err := r.replica.WithContext(ctx).
		Where("link_id = ?", linkID).
		Order("priority ASC, created_at ASC").
		Find(&rules).Error
	return rules, err
}

func (r *linkGeoRuleRepository) GetByLinkIDAndCountry(ctx context.Context, linkID uuid.UUID, countryCode string) (*model.LinkGeoRule, error) {
	var rule model.LinkGeoRule
	err := r.replica.WithContext(ctx).
		Where("link_id = ? AND country_code = ?", linkID, countryCode).
		Order("priority ASC").
		First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *linkGeoRuleRepository) Create(ctx context.Context, rule *model.LinkGeoRule) error {
	return r.master.WithContext(ctx).Create(rule).Error
}

func (r *linkGeoRuleRepository) Update(ctx context.Context, rule *model.LinkGeoRule) error {
	return r.master.WithContext(ctx).Save(rule).Error
}

func (r *linkGeoRuleRepository) Delete(ctx context.Context, id, linkID uuid.UUID) error {
	result := r.master.WithContext(ctx).
		Where("id = ? AND link_id = ?", id, linkID).
		Delete(&model.LinkGeoRule{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *linkGeoRuleRepository) DeleteByLinkID(ctx context.Context, linkID uuid.UUID) error {
	return r.master.WithContext(ctx).
		Where("link_id = ?", linkID).
		Delete(&model.LinkGeoRule{}).Error
}
