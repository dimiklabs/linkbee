package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/model"
)

type FolderRepositoryI interface {
	Create(ctx context.Context, folder *model.Folder) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Folder, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Folder, error)
	Update(ctx context.Context, folder *model.Folder) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type folderRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewFolderRepository(masterDB, replicaDB *gorm.DB) FolderRepositoryI {
	return &folderRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

func (r *folderRepository) Create(ctx context.Context, folder *model.Folder) error {
	return r.masterDB.WithContext(ctx).Create(folder).Error
}

func (r *folderRepository) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Folder, error) {
	var folder model.Folder
	if err := r.replicaDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&folder).Error; err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *folderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Folder, error) {
	var folders []model.Folder
	if err := r.replicaDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *folderRepository) Update(ctx context.Context, folder *model.Folder) error {
	return r.masterDB.WithContext(ctx).Save(folder).Error
}

// Delete unlinks all links in the folder (sets folder_id = NULL) then soft-deletes the folder.
func (r *folderRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := r.masterDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("folder_id = ?", id).
		Update("folder_id", nil).Error; err != nil {
		return err
	}

	result := r.masterDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Folder{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
