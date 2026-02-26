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

type LinkRepositoryI interface {
	// Create
	Create(ctx context.Context, link *model.Link) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*model.Link, error)
	GetBySlug(ctx context.Context, slug string) (*model.Link, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int, search string, folderID *uuid.UUID, starred *bool, healthStatus string) ([]model.Link, int64, error)
	SlugExists(ctx context.Context, slug string) (bool, error)

	// Duplicate check
	FindByDestinationURL(ctx context.Context, userID uuid.UUID, destURL string) (*model.Link, error)

	// Update
	Update(ctx context.Context, link *model.Link) error
	IncrementClickCount(ctx context.Context, id uuid.UUID) error
	ToggleStar(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error)
	UpdateHealthStatus(ctx context.Context, id uuid.UUID, status string, statusCode int) error

	// Health check
	GetLinksForHealthCheck(ctx context.Context, staleBefore time.Time, limit int) ([]model.Link, error)

	// Count (all links, admin use)
	Count(ctx context.Context) (int64, error)

	// Delete
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type LinkRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewLinkRepository(masterDB, replicaDB *gorm.DB) LinkRepositoryI {
	return &LinkRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

// --------------- Create ---------------

func (r *LinkRepository) Create(ctx context.Context, link *model.Link) error {
	logger.DebugCtx(ctx, "Creating link",
		zap.String("slug", link.Slug))

	if err := r.masterDB.WithContext(ctx).Create(link).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create link",
			zap.String("slug", link.Slug),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Link created",
		zap.String("link_id", link.ID.String()),
		zap.String("slug", link.Slug))
	return nil
}

// --------------- Read ---------------

func (r *LinkRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Link, error) {
	logger.DebugCtx(ctx, "Fetching link by ID",
		zap.String("link_id", id.String()))

	var link model.Link
	if err := r.replicaDB.WithContext(ctx).
		Where("id = ?", id).
		First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) GetBySlug(ctx context.Context, slug string) (*model.Link, error) {
	logger.DebugCtx(ctx, "Fetching link by slug",
		zap.String("slug", slug))

	var link model.Link
	if err := r.replicaDB.WithContext(ctx).
		Where("slug = ? AND is_active = ?", slug, true).
		First(&link).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) FindByDestinationURL(ctx context.Context, userID uuid.UUID, destURL string) (*model.Link, error) {
	var link model.Link
	err := r.replicaDB.WithContext(ctx).
		Where("user_id = ? AND destination_url = ? AND deleted_at IS NULL", userID, destURL).
		Order("created_at ASC").
		First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int, search string, folderID *uuid.UUID, starred *bool, healthStatus string) ([]model.Link, int64, error) {
	logger.DebugCtx(ctx, "Fetching links for user",
		zap.String("user_id", userID.String()),
		zap.Int("page", page),
		zap.Int("limit", limit))

	offset := (page - 1) * limit

	query := r.replicaDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("user_id = ?", userID)

	if folderID != nil {
		query = query.Where("folder_id = ?", *folderID)
	}

	if starred != nil {
		query = query.Where("is_starred = ?", *starred)
	}

	if healthStatus != "" {
		query = query.Where("health_status = ?", healthStatus)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where(
			"(slug ILIKE ? OR title ILIKE ? OR destination_url ILIKE ?)",
			searchTerm, searchTerm, searchTerm,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count links for user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, 0, err
	}

	var links []model.Link
	if err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&links).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch links for user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, 0, err
	}

	return links, total, nil
}

func (r *LinkRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("slug = ?", slug).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to check slug existence",
			zap.String("slug", slug),
			zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

// --------------- Update ---------------

func (r *LinkRepository) Update(ctx context.Context, link *model.Link) error {
	logger.DebugCtx(ctx, "Updating link",
		zap.String("link_id", link.ID.String()))

	link.UpdatedAt = time.Now()
	if err := r.masterDB.WithContext(ctx).Save(link).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update link",
			zap.String("link_id", link.ID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Link updated",
		zap.String("link_id", link.ID.String()))
	return nil
}

func (r *LinkRepository) IncrementClickCount(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Incrementing click count",
		zap.String("link_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("id = ?", id).
		Update("click_count", gorm.Expr("click_count + 1")).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to increment click count",
			zap.String("link_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *LinkRepository) UpdateHealthStatus(ctx context.Context, id uuid.UUID, status string, statusCode int) error {
	now := time.Now()
	if err := r.masterDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"health_status":      status,
			"health_status_code": statusCode,
			"health_checked_at":  now,
			"updated_at":         now,
		}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update health status",
			zap.String("link_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *LinkRepository) GetLinksForHealthCheck(ctx context.Context, staleBefore time.Time, limit int) ([]model.Link, error) {
	var links []model.Link
	err := r.replicaDB.WithContext(ctx).
		Where("is_active = ? AND (health_checked_at IS NULL OR health_checked_at < ?)", true, staleBefore).
		Order("health_checked_at ASC NULLS FIRST").
		Limit(limit).
		Find(&links).Error
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch links for health check", zap.Error(err))
		return nil, err
	}
	return links, nil
}

func (r *LinkRepository) ToggleStar(ctx context.Context, id uuid.UUID, userID uuid.UUID) (bool, error) {
	logger.DebugCtx(ctx, "Toggling star for link",
		zap.String("link_id", id.String()),
		zap.String("user_id", userID.String()))

	var link model.Link
	if err := r.masterDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&link).Error; err != nil {
		return false, err
	}

	link.IsStarred = !link.IsStarred
	link.UpdatedAt = time.Now()
	if err := r.masterDB.WithContext(ctx).
		Model(&link).
		Update("is_starred", link.IsStarred).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to toggle star",
			zap.String("link_id", id.String()),
			zap.Error(err))
		return false, err
	}

	logger.InfoCtx(ctx, "Link star toggled",
		zap.String("link_id", id.String()),
		zap.Bool("is_starred", link.IsStarred))
	return link.IsStarred, nil
}

// --------------- Delete ---------------

func (r *LinkRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting link",
		zap.String("link_id", id.String()),
		zap.String("user_id", userID.String()))

	result := r.masterDB.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Link{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete link",
			zap.String("link_id", id.String()),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.InfoCtx(ctx, "Link deleted",
		zap.String("link_id", id.String()))
	return nil
}

func (r *LinkRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).Model(&model.Link{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *LinkRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting all links for user",
		zap.String("user_id", userID.String()))

	result := r.masterDB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Link{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete links for user",
			zap.String("user_id", userID.String()),
			zap.Error(result.Error))
		return result.Error
	}

	logger.InfoCtx(ctx, "Links deleted for user",
		zap.String("user_id", userID.String()),
		zap.Int64("count", result.RowsAffected))
	return nil
}
