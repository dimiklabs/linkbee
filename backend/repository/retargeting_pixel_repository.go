package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type RetargetingPixelRepositoryI interface {
	Create(ctx context.Context, pixel *model.RetargetingPixel) error
	GetByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.RetargetingPixel, error)
	GetActiveByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.RetargetingPixel, error)
	Delete(ctx context.Context, id, linkID uuid.UUID) error
}

type retargetingPixelRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewRetargetingPixelRepository(master, replica *gorm.DB) RetargetingPixelRepositoryI {
	return &retargetingPixelRepository{master: master, replica: replica}
}

func (r *retargetingPixelRepository) Create(ctx context.Context, pixel *model.RetargetingPixel) error {
	return r.master.WithContext(ctx).Create(pixel).Error
}

func (r *retargetingPixelRepository) GetByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.RetargetingPixel, error) {
	var pixels []model.RetargetingPixel
	err := r.replica.WithContext(ctx).
		Where("link_id = ?", linkID).
		Order("created_at ASC").
		Find(&pixels).Error
	return pixels, err
}

func (r *retargetingPixelRepository) GetActiveByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.RetargetingPixel, error) {
	var pixels []model.RetargetingPixel
	err := r.replica.WithContext(ctx).
		Where("link_id = ? AND is_active = true", linkID).
		Order("created_at ASC").
		Find(&pixels).Error
	return pixels, err
}

func (r *retargetingPixelRepository) Delete(ctx context.Context, id, linkID uuid.UUID) error {
	result := r.master.WithContext(ctx).
		Where("id = ? AND link_id = ?", id, linkID).
		Delete(&model.RetargetingPixel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
