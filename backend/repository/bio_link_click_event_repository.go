package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/model"
)

type BioLinkClickEventRepositoryI interface {
	BulkCreate(ctx context.Context, events []*model.BioLinkClickEvent) error
	GetClicksByBioLinkID(ctx context.Context, bioLinkID uuid.UUID) (int64, error)
	GetClicksByBioPageID(ctx context.Context, bioPageID uuid.UUID) (int64, error)
}

type BioLinkClickEventRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewBioLinkClickEventRepository(master, replica *gorm.DB) BioLinkClickEventRepositoryI {
	return &BioLinkClickEventRepository{master: master, replica: replica}
}

func (r *BioLinkClickEventRepository) BulkCreate(ctx context.Context, events []*model.BioLinkClickEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.master.WithContext(ctx).CreateInBatches(events, 100).Error
}

func (r *BioLinkClickEventRepository) GetClicksByBioLinkID(ctx context.Context, bioLinkID uuid.UUID) (int64, error) {
	var count int64
	err := r.replica.WithContext(ctx).
		Model(&model.BioLinkClickEvent{}).
		Where("bio_link_id = ?", bioLinkID).
		Count(&count).Error
	return count, err
}

func (r *BioLinkClickEventRepository) GetClicksByBioPageID(ctx context.Context, bioPageID uuid.UUID) (int64, error) {
	var count int64
	err := r.replica.WithContext(ctx).
		Model(&model.BioLinkClickEvent{}).
		Where("bio_page_id = ?", bioPageID).
		Count(&count).Error
	return count, err
}
