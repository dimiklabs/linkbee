package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/model"
)

type BioRepositoryI interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.BioPage, error)
	GetByUsername(ctx context.Context, username string) (*model.BioPage, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	Create(ctx context.Context, page *model.BioPage) error
	Update(ctx context.Context, page *model.BioPage) error

	ListLinks(ctx context.Context, bioPageID uuid.UUID) ([]model.BioLink, error)
	GetLink(ctx context.Context, id, bioPageID uuid.UUID) (*model.BioLink, error)
	CreateLink(ctx context.Context, link *model.BioLink) error
	UpdateLink(ctx context.Context, link *model.BioLink) error
	DeleteLink(ctx context.Context, id, bioPageID uuid.UUID) error
	ReorderLinks(ctx context.Context, bioPageID uuid.UUID, ids []uuid.UUID) error
}

type BioRepository struct {
	master  *gorm.DB
	replica *gorm.DB
}

func NewBioRepository(master, replica *gorm.DB) BioRepositoryI {
	return &BioRepository{master: master, replica: replica}
}

func (r *BioRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.BioPage, error) {
	var page model.BioPage
	err := r.replica.WithContext(ctx).
		Preload("Links", func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *BioRepository) GetByUsername(ctx context.Context, username string) (*model.BioPage, error) {
	var page model.BioPage
	err := r.replica.WithContext(ctx).
		Preload("Links", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = true").Order("position ASC")
		}).
		Where("username = ? AND is_published = true AND deleted_at IS NULL", username).
		First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

func (r *BioRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.replica.WithContext(ctx).Model(&model.BioPage{}).
		Where("username = ? AND deleted_at IS NULL", username).
		Count(&count).Error
	return count > 0, err
}

func (r *BioRepository) Create(ctx context.Context, page *model.BioPage) error {
	return r.master.WithContext(ctx).Create(page).Error
}

func (r *BioRepository) Update(ctx context.Context, page *model.BioPage) error {
	return r.master.WithContext(ctx).Save(page).Error
}

func (r *BioRepository) ListLinks(ctx context.Context, bioPageID uuid.UUID) ([]model.BioLink, error) {
	var links []model.BioLink
	err := r.replica.WithContext(ctx).
		Where("bio_page_id = ?", bioPageID).
		Order("position ASC").
		Find(&links).Error
	return links, err
}

func (r *BioRepository) GetLink(ctx context.Context, id, bioPageID uuid.UUID) (*model.BioLink, error) {
	var link model.BioLink
	err := r.replica.WithContext(ctx).
		Where("id = ? AND bio_page_id = ?", id, bioPageID).
		First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *BioRepository) CreateLink(ctx context.Context, link *model.BioLink) error {
	return r.master.WithContext(ctx).Create(link).Error
}

func (r *BioRepository) UpdateLink(ctx context.Context, link *model.BioLink) error {
	return r.master.WithContext(ctx).Save(link).Error
}

func (r *BioRepository) DeleteLink(ctx context.Context, id, bioPageID uuid.UUID) error {
	return r.master.WithContext(ctx).
		Where("id = ? AND bio_page_id = ?", id, bioPageID).
		Delete(&model.BioLink{}).Error
}

func (r *BioRepository) ReorderLinks(ctx context.Context, bioPageID uuid.UUID, ids []uuid.UUID) error {
	return r.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.BioLink{}).
				Where("id = ? AND bio_page_id = ?", id, bioPageID).
				Update("position", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
