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

type LinkVariantRepositoryI interface {
	Create(ctx context.Context, variant *model.LinkVariant) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.LinkVariant, error)
	GetByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.LinkVariant, error)
	Update(ctx context.Context, variant *model.LinkVariant) error
	Delete(ctx context.Context, id uuid.UUID, linkID uuid.UUID) error
	IncrementClickCount(ctx context.Context, id uuid.UUID) error
	DeleteByLinkID(ctx context.Context, linkID uuid.UUID) error
}

type linkVariantRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewLinkVariantRepository(masterDB, replicaDB *gorm.DB) LinkVariantRepositoryI {
	return &linkVariantRepository{masterDB: masterDB, replicaDB: replicaDB}
}

func (r *linkVariantRepository) Create(ctx context.Context, variant *model.LinkVariant) error {
	if err := r.masterDB.WithContext(ctx).Create(variant).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create link variant",
			zap.String("link_id", variant.LinkID.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *linkVariantRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.LinkVariant, error) {
	var v model.LinkVariant
	if err := r.replicaDB.WithContext(ctx).Where("id = ?", id).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *linkVariantRepository) GetByLinkID(ctx context.Context, linkID uuid.UUID) ([]model.LinkVariant, error) {
	var variants []model.LinkVariant
	if err := r.replicaDB.WithContext(ctx).
		Where("link_id = ?", linkID).
		Order("created_at ASC").
		Find(&variants).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get variants for link",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}
	return variants, nil
}

func (r *linkVariantRepository) Update(ctx context.Context, variant *model.LinkVariant) error {
	variant.UpdatedAt = time.Now()
	if err := r.masterDB.WithContext(ctx).Save(variant).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update link variant",
			zap.String("variant_id", variant.ID.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *linkVariantRepository) Delete(ctx context.Context, id uuid.UUID, linkID uuid.UUID) error {
	result := r.masterDB.WithContext(ctx).
		Where("id = ? AND link_id = ?", id, linkID).
		Delete(&model.LinkVariant{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *linkVariantRepository) IncrementClickCount(ctx context.Context, id uuid.UUID) error {
	return r.masterDB.WithContext(ctx).
		Model(&model.LinkVariant{}).
		Where("id = ?", id).
		Update("click_count", gorm.Expr("click_count + 1")).Error
}

func (r *linkVariantRepository) DeleteByLinkID(ctx context.Context, linkID uuid.UUID) error {
	return r.masterDB.WithContext(ctx).
		Where("link_id = ?", linkID).
		Delete(&model.LinkVariant{}).Error
}
