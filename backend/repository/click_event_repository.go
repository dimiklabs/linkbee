package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
)

// TimeSeriesPoint is a single data point in a click time series.
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

// ReferrerPoint holds referrer stats for a link.
type ReferrerPoint struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

// DevicePoint holds device type stats for a link.
type DevicePoint struct {
	DeviceType string `json:"device_type"`
	Count      int64  `json:"count"`
}

// CountryPoint holds per-country click stats for a link.
type CountryPoint struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

// BrowserPoint holds per-browser click stats for a link.
type BrowserPoint struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}

// OSPoint holds per-OS click stats for a link.
type OSPoint struct {
	OS    string `json:"os"`
	Count int64  `json:"count"`
}

type ClickEventRepositoryI interface {
	// Write
	BulkCreate(ctx context.Context, events []*model.ClickEvent) error

	// Aggregates
	GetClickCountByLinkID(ctx context.Context, linkID uuid.UUID) (int64, error)
	GetUniqueClickCountByLinkID(ctx context.Context, linkID uuid.UUID) (int64, error)
	GetTimeSeriesData(ctx context.Context, linkID uuid.UUID, from, to time.Time, granularity string) ([]TimeSeriesPoint, error)
	GetTopReferrers(ctx context.Context, linkID uuid.UUID, limit int) ([]ReferrerPoint, error)
	GetDeviceBreakdown(ctx context.Context, linkID uuid.UUID) ([]DevicePoint, error)
	GetCountryBreakdown(ctx context.Context, linkID uuid.UUID) ([]CountryPoint, error)
	GetBrowserBreakdown(ctx context.Context, linkID uuid.UUID) ([]BrowserPoint, error)
	GetOSBreakdown(ctx context.Context, linkID uuid.UUID) ([]OSPoint, error)

	// Delete
	DeleteByLinkID(ctx context.Context, linkID uuid.UUID) error
}

type ClickEventRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewClickEventRepository(masterDB, replicaDB *gorm.DB) ClickEventRepositoryI {
	return &ClickEventRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

// --------------- Write ---------------

func (r *ClickEventRepository) BulkCreate(ctx context.Context, events []*model.ClickEvent) error {
	if len(events) == 0 {
		return nil
	}

	logger.DebugCtx(ctx, "Bulk inserting click events",
		zap.Int("count", len(events)))

	if err := r.masterDB.WithContext(ctx).CreateInBatches(events, 100).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to bulk insert click events",
			zap.Error(err))
		return err
	}

	logger.DebugCtx(ctx, "Click events inserted",
		zap.Int("count", len(events)))
	return nil
}

// --------------- Aggregates ---------------

func (r *ClickEventRepository) GetClickCountByLinkID(ctx context.Context, linkID uuid.UUID) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Where("link_id = ?", linkID).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count clicks for link",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetUniqueClickCountByLinkID(ctx context.Context, linkID uuid.UUID) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Where("link_id = ? AND ip_hash != ''", linkID).
		Distinct("ip_hash").
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count unique clicks for link",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetTimeSeriesData(ctx context.Context, linkID uuid.UUID, from, to time.Time, granularity string) ([]TimeSeriesPoint, error) {
	// Validate granularity — default to day
	allowed := map[string]bool{"hour": true, "day": true, "week": true, "month": true}
	if !allowed[granularity] {
		granularity = "day"
	}

	type result struct {
		Timestamp time.Time
		Count     int64
	}

	// Granularity is already validated above, safe to use in SQL
	truncExpr := fmt.Sprintf("date_trunc('%s', clicked_at)", granularity)
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select(truncExpr+" AS timestamp, COUNT(*) AS count").
		Where("link_id = ? AND clicked_at BETWEEN ? AND ?", linkID, from, to).
		Group(truncExpr).
		Order("timestamp ASC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get time series data",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]TimeSeriesPoint, len(rows))
	for i, r := range rows {
		points[i] = TimeSeriesPoint{Timestamp: r.Timestamp, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetTopReferrers(ctx context.Context, linkID uuid.UUID, limit int) ([]ReferrerPoint, error) {
	type result struct {
		Referrer string
		Count    int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(referrer, ''), 'Direct') AS referrer, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(referrer, ''), 'Direct')").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get top referrers",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]ReferrerPoint, len(rows))
	for i, r := range rows {
		points[i] = ReferrerPoint{Referrer: r.Referrer, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetDeviceBreakdown(ctx context.Context, linkID uuid.UUID) ([]DevicePoint, error) {
	type result struct {
		DeviceType string
		Count      int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(device_type, ''), 'unknown') AS device_type, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(device_type, ''), 'unknown')").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get device breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]DevicePoint, len(rows))
	for i, r := range rows {
		points[i] = DevicePoint{DeviceType: r.DeviceType, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetCountryBreakdown(ctx context.Context, linkID uuid.UUID) ([]CountryPoint, error) {
	type result struct {
		Country string
		Count   int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(country, ''), 'Unknown') AS country, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(country, ''), 'Unknown')").
		Order("count DESC").
		Limit(15).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get country breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]CountryPoint, len(rows))
	for i, r := range rows {
		points[i] = CountryPoint{Country: r.Country, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetBrowserBreakdown(ctx context.Context, linkID uuid.UUID) ([]BrowserPoint, error) {
	type result struct {
		Browser string
		Count   int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(browser, ''), 'Unknown') AS browser, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(browser, ''), 'Unknown')").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get browser breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]BrowserPoint, len(rows))
	for i, r := range rows {
		points[i] = BrowserPoint{Browser: r.Browser, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetOSBreakdown(ctx context.Context, linkID uuid.UUID) ([]OSPoint, error) {
	type result struct {
		OS    string
		Count int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(os, ''), 'Unknown') AS os, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(os, ''), 'Unknown')").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get OS breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]OSPoint, len(rows))
	for i, r := range rows {
		points[i] = OSPoint{OS: r.OS, Count: r.Count}
	}
	return points, nil
}

// --------------- Delete ---------------

func (r *ClickEventRepository) DeleteByLinkID(ctx context.Context, linkID uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting click events for link",
		zap.String("link_id", linkID.String()))

	result := r.masterDB.WithContext(ctx).
		Where("link_id = ?", linkID).
		Delete(&model.ClickEvent{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete click events for link",
			zap.String("link_id", linkID.String()),
			zap.Error(result.Error))
		return result.Error
	}

	logger.InfoCtx(ctx, "Click events deleted for link",
		zap.String("link_id", linkID.String()),
		zap.Int64("count", result.RowsAffected))
	return nil
}
