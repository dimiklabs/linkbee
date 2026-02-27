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

// CityPoint holds per-city click stats for a link.
type CityPoint struct {
	City    string `json:"city"`
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

// SourcePoint holds click count broken down by source (web, qr, api).
type SourcePoint struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

// HeatmapPoint holds click count for a specific hour-of-day × day-of-week cell.
type HeatmapPoint struct {
	DayOfWeek int   `json:"day_of_week"` // 0 = Sunday … 6 = Saturday
	Hour      int   `json:"hour"`        // 0–23 UTC
	Count     int64 `json:"count"`
}

// UTMPoint holds click count for a single UTM parameter value.
type UTMPoint struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// ReferrerCategoryPoint holds click count broken down by traffic channel category.
type ReferrerCategoryPoint struct {
	Category string `json:"category"` // direct | search | social | email | referral
	Count    int64  `json:"count"`
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
	GetCityBreakdown(ctx context.Context, linkID uuid.UUID, limit int) ([]CityPoint, error)
	GetBrowserBreakdown(ctx context.Context, linkID uuid.UUID) ([]BrowserPoint, error)
	GetOSBreakdown(ctx context.Context, linkID uuid.UUID) ([]OSPoint, error)
	GetSourceBreakdown(ctx context.Context, linkID uuid.UUID) ([]SourcePoint, error)
	GetHeatmapData(ctx context.Context, linkID uuid.UUID, from, to time.Time) ([]HeatmapPoint, error)
	GetUTMBreakdown(ctx context.Context, linkID uuid.UUID, field string, limit int) ([]UTMPoint, error)
	GetReferrerCategoryBreakdown(ctx context.Context, linkID uuid.UUID) ([]ReferrerCategoryPoint, error)
	GetReturnVisitorStats(ctx context.Context, linkID uuid.UUID) (firstTime int64, returning int64, err error)
	GetTotalClicksByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
	GetClicksInPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
	GetClicksInPeriodByLinkID(ctx context.Context, linkID uuid.UUID, from, to time.Time) (int64, error)
	GetUniqueClicksInPeriodByLinkID(ctx context.Context, linkID uuid.UUID, from, to time.Time) (int64, error)
	GetTimeSeriesByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, granularity string) ([]TimeSeriesPoint, error)
	GetUniqueClicksInPeriodByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error)
	GetCountryBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]CountryPoint, error)
	GetDeviceBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]DevicePoint, error)
	GetBrowserBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]BrowserPoint, error)
	GetTopReferrersByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]ReferrerPoint, error)

	GetOSBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]OSPoint, error)
	GetCityBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]CityPoint, error)
	GetSourceBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]SourcePoint, error)
	GetReferrerCategoryBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]ReferrerCategoryPoint, error)
	GetUTMBreakdownByUserID(ctx context.Context, userID uuid.UUID, field string, from, to time.Time, limit int) ([]UTMPoint, error)

	GetClickCountByFolderID(ctx context.Context, folderID uuid.UUID) (int64, error)
	GetTotalClicks(ctx context.Context) (int64, error)

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

func (r *ClickEventRepository) GetCityBreakdown(ctx context.Context, linkID uuid.UUID, limit int) ([]CityPoint, error) {
	type result struct {
		City    string
		Country string
		Count   int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("city, COALESCE(NULLIF(country, ''), 'Unknown') AS country, COUNT(*) AS count").
		Where("link_id = ? AND city != ''", linkID).
		Group("city, COALESCE(NULLIF(country, ''), 'Unknown')").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get city breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]CityPoint, len(rows))
	for i, r := range rows {
		points[i] = CityPoint{City: r.City, Country: r.Country, Count: r.Count}
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

func (r *ClickEventRepository) GetSourceBreakdown(ctx context.Context, linkID uuid.UUID) ([]SourcePoint, error) {
	type result struct {
		Source string
		Count  int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(source, ''), 'web') AS source, COUNT(*) AS count").
		Where("link_id = ?", linkID).
		Group("COALESCE(NULLIF(source, ''), 'web')").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get source breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]SourcePoint, len(rows))
	for i, r := range rows {
		points[i] = SourcePoint{Source: r.Source, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetHeatmapData(ctx context.Context, linkID uuid.UUID, from, to time.Time) ([]HeatmapPoint, error) {
	type result struct {
		DayOfWeek int
		Hour      int
		Count     int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("EXTRACT(DOW FROM clicked_at)::int AS day_of_week, EXTRACT(HOUR FROM clicked_at)::int AS hour, COUNT(*) AS count").
		Where("link_id = ? AND clicked_at BETWEEN ? AND ?", linkID, from, to).
		Group("EXTRACT(DOW FROM clicked_at), EXTRACT(HOUR FROM clicked_at)").
		Order("day_of_week, hour").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get heatmap data",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]HeatmapPoint, len(rows))
	for i, r := range rows {
		points[i] = HeatmapPoint{DayOfWeek: r.DayOfWeek, Hour: r.Hour, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetUTMBreakdown(ctx context.Context, linkID uuid.UUID, field string, limit int) ([]UTMPoint, error) {
	// Validate field against allowed UTM columns to prevent SQL injection
	allowed := map[string]bool{
		"utm_source": true, "utm_medium": true,
		"utm_campaign": true, "utm_content": true, "utm_term": true,
	}
	if !allowed[field] {
		return nil, nil
	}

	type result struct {
		Value string
		Count int64
	}

	expr := fmt.Sprintf("NULLIF(%s, '')", field)
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select(expr+" AS value, COUNT(*) AS count").
		Where("link_id = ? AND "+field+" != ''", linkID).
		Group(expr).
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM breakdown",
			zap.String("link_id", linkID.String()),
			zap.String("field", field),
			zap.Error(err))
		return nil, err
	}

	points := make([]UTMPoint, len(rows))
	for i, r := range rows {
		points[i] = UTMPoint{Value: r.Value, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetReferrerCategoryBreakdown(ctx context.Context, linkID uuid.UUID) ([]ReferrerCategoryPoint, error) {
	type result struct {
		Category string
		Count    int64
	}

	// Use PostgreSQL case-insensitive regex (~*) to classify referrers into
	// five standard marketing channels.  Order of WHEN clauses matters: direct
	// is checked first (empty string), then the more-specific categories before
	// the catch-all referral bucket.
	q := `
		SELECT category, COUNT(*) AS count
		FROM (
			SELECT CASE
				WHEN COALESCE(NULLIF(referrer, ''), '') = ''
					THEN 'direct'
				WHEN referrer ~* 'google\.|bing\.com|yahoo\.com|duckduckgo\.com|baidu\.com|yandex\.|ask\.com|ecosia\.org|qwant\.com|startpage\.com|search\.yahoo\.com'
					THEN 'search'
				WHEN referrer ~* 'facebook\.com|fb\.com/|twitter\.com|t\.co/|instagram\.com|linkedin\.com|youtube\.com|youtu\.be|tiktok\.com|pinterest\.com|reddit\.com|snapchat\.com|discord\.com|discord\.gg|tumblr\.com|mastodon\.|threads\.net|bsky\.app'
					THEN 'social'
				WHEN referrer ~* 'mail\.google\.com|outlook\.|mail\.yahoo\.com|mail\.proton\.me|protonmail\.com|webmail\.|email\.'
					THEN 'email'
				ELSE 'referral'
			END AS category
			FROM click_events
			WHERE link_id = $1
		) sub
		GROUP BY category
		ORDER BY count DESC
	`

	var rows []result
	if err := r.replicaDB.WithContext(ctx).Raw(q, linkID).Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get referrer category breakdown",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]ReferrerCategoryPoint, len(rows))
	for i, row := range rows {
		points[i] = ReferrerCategoryPoint{Category: row.Category, Count: row.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetReturnVisitorStats(ctx context.Context, linkID uuid.UUID) (firstTime int64, returning int64, err error) {
	type result struct {
		FirstTime int64
		Returning int64
	}

	// A visitor is "first-time" if their ip_hash appears exactly once for this
	// link; "returning" if it appears more than once.  Anonymous clicks
	// (empty ip_hash) are excluded from both counts.
	q := `
		SELECT
			SUM(CASE WHEN visit_count = 1 THEN 1 ELSE 0 END) AS first_time,
			SUM(CASE WHEN visit_count > 1 THEN 1 ELSE 0 END) AS returning
		FROM (
			SELECT ip_hash, COUNT(*) AS visit_count
			FROM click_events
			WHERE link_id = $1 AND ip_hash != ''
			GROUP BY ip_hash
		) sub
	`

	var row result
	if dbErr := r.replicaDB.WithContext(ctx).Raw(q, linkID).Scan(&row).Error; dbErr != nil {
		logger.ErrorCtx(ctx, "Failed to get return visitor stats",
			zap.String("link_id", linkID.String()),
			zap.Error(dbErr))
		return 0, 0, dbErr
	}
	return row.FirstTime, row.Returning, nil
}

func (r *ClickEventRepository) GetClicksInPeriodByLinkID(ctx context.Context, linkID uuid.UUID, from, to time.Time) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Where("link_id = ? AND clicked_at BETWEEN ? AND ?", linkID, from, to).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count clicks in period for link",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetUniqueClicksInPeriodByLinkID(ctx context.Context, linkID uuid.UUID, from, to time.Time) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Where("link_id = ? AND ip_hash != '' AND clicked_at BETWEEN ? AND ?", linkID, from, to).
		Distinct("ip_hash").
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count unique clicks in period for link",
			zap.String("link_id", linkID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetTotalClicksByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL", userID).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count total clicks for user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetClicksInPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ?", userID, from, to).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count clicks in period for user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetUniqueClicksInPeriodByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ? AND click_events.ip_hash != ''", userID, from, to).
		Distinct("click_events.ip_hash").
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count unique clicks in period by user", zap.String("user_id", userID.String()), zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetCountryBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]CountryPoint, error) {
	type result struct {
		Country string
		Count   int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("click_events.country, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ? AND click_events.country != ''", userID, from, to).
		Group("click_events.country").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get country breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]CountryPoint, len(rows))
	for i, r := range rows {
		points[i] = CountryPoint{Country: r.Country, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetDeviceBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]DevicePoint, error) {
	type result struct {
		DeviceType string
		Count      int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("click_events.device_type, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ?", userID, from, to).
		Group("click_events.device_type").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get device breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]DevicePoint, len(rows))
	for i, r := range rows {
		points[i] = DevicePoint{DeviceType: r.DeviceType, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetBrowserBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]BrowserPoint, error) {
	type result struct {
		Browser string
		Count   int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("click_events.browser, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ? AND click_events.browser != ''", userID, from, to).
		Group("click_events.browser").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get browser breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]BrowserPoint, len(rows))
	for i, r := range rows {
		points[i] = BrowserPoint{Browser: r.Browser, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetTopReferrersByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]ReferrerPoint, error) {
	type result struct {
		Referrer string
		Count    int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(click_events.referrer, ''), 'Direct') AS referrer, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ?", userID, from, to).
		Group("COALESCE(NULLIF(click_events.referrer, ''), 'Direct')").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get top referrers by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]ReferrerPoint, len(rows))
	for i, r := range rows {
		points[i] = ReferrerPoint{Referrer: r.Referrer, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetTimeSeriesByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, granularity string) ([]TimeSeriesPoint, error) {
	allowed := map[string]bool{"hour": true, "day": true, "week": true, "month": true}
	if !allowed[granularity] {
		granularity = "day"
	}

	type result struct {
		Timestamp time.Time
		Count     int64
	}

	truncExpr := fmt.Sprintf("date_trunc('%s', click_events.clicked_at)", granularity)
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select(truncExpr+" AS timestamp, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ?", userID, from, to).
		Group(truncExpr).
		Order("timestamp ASC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get time series by user",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return nil, err
	}

	points := make([]TimeSeriesPoint, len(rows))
	for i, r := range rows {
		points[i] = TimeSeriesPoint{Timestamp: r.Timestamp, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetOSBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]OSPoint, error) {
	type result struct {
		OS    string
		Count int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(click_events.os, ''), 'Unknown') AS os, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ?", userID, from, to).
		Group("COALESCE(NULLIF(click_events.os, ''), 'Unknown')").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get OS breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]OSPoint, len(rows))
	for i, r := range rows {
		points[i] = OSPoint{OS: r.OS, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetCityBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time, limit int) ([]CityPoint, error) {
	type result struct {
		City    string
		Country string
		Count   int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("click_events.city, COALESCE(NULLIF(click_events.country, ''), 'Unknown') AS country, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ? AND click_events.city != ''", userID, from, to).
		Group("click_events.city, COALESCE(NULLIF(click_events.country, ''), 'Unknown')").
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get city breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]CityPoint, len(rows))
	for i, r := range rows {
		points[i] = CityPoint{City: r.City, Country: r.Country, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetSourceBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]SourcePoint, error) {
	type result struct {
		Source string
		Count  int64
	}
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select("COALESCE(NULLIF(click_events.source, ''), 'web') AS source, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ?", userID, from, to).
		Group("COALESCE(NULLIF(click_events.source, ''), 'web')").
		Order("count DESC").
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get source breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]SourcePoint, len(rows))
	for i, r := range rows {
		points[i] = SourcePoint{Source: r.Source, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetReferrerCategoryBreakdownByUserID(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]ReferrerCategoryPoint, error) {
	type result struct {
		Category string
		Count    int64
	}

	q := `
		SELECT category, COUNT(*) AS count
		FROM (
			SELECT CASE
				WHEN COALESCE(NULLIF(ce.referrer, ''), '') = ''
					THEN 'direct'
				WHEN ce.referrer ~* 'google\.|bing\.com|yahoo\.com|duckduckgo\.com|baidu\.com|yandex\.|ask\.com|ecosia\.org|qwant\.com|startpage\.com|search\.yahoo\.com'
					THEN 'search'
				WHEN ce.referrer ~* 'facebook\.com|fb\.com/|twitter\.com|t\.co/|instagram\.com|linkedin\.com|youtube\.com|youtu\.be|tiktok\.com|pinterest\.com|reddit\.com|snapchat\.com|discord\.com|discord\.gg|tumblr\.com|mastodon\.|threads\.net|bsky\.app'
					THEN 'social'
				WHEN ce.referrer ~* 'mail\.google\.com|outlook\.|mail\.yahoo\.com|mail\.proton\.me|protonmail\.com|webmail\.|email\.'
					THEN 'email'
				ELSE 'referral'
			END AS category
			FROM click_events ce
			JOIN links l ON l.id = ce.link_id
			WHERE l.user_id = $1 AND l.deleted_at IS NULL AND ce.clicked_at BETWEEN $2 AND $3
		) sub
		GROUP BY category
		ORDER BY count DESC
	`

	var rows []result
	if err := r.replicaDB.WithContext(ctx).Raw(q, userID, from, to).Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get referrer category breakdown by user", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}
	points := make([]ReferrerCategoryPoint, len(rows))
	for i, row := range rows {
		points[i] = ReferrerCategoryPoint{Category: row.Category, Count: row.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetUTMBreakdownByUserID(ctx context.Context, userID uuid.UUID, field string, from, to time.Time, limit int) ([]UTMPoint, error) {
	// Validate field against allowed UTM columns to prevent SQL injection
	allowed := map[string]bool{
		"utm_source": true, "utm_medium": true,
		"utm_campaign": true, "utm_content": true, "utm_term": true,
	}
	if !allowed[field] {
		return nil, nil
	}

	type result struct {
		Value string
		Count int64
	}

	expr := fmt.Sprintf("NULLIF(click_events.%s, '')", field)
	whereField := fmt.Sprintf("click_events.%s != ''", field)
	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Select(expr+" AS value, COUNT(*) AS count").
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.user_id = ? AND links.deleted_at IS NULL AND click_events.clicked_at BETWEEN ? AND ? AND "+whereField, userID, from, to).
		Group(expr).
		Order("count DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM breakdown by user",
			zap.String("user_id", userID.String()),
			zap.String("field", field),
			zap.Error(err))
		return nil, err
	}

	points := make([]UTMPoint, len(rows))
	for i, r := range rows {
		points[i] = UTMPoint{Value: r.Value, Count: r.Count}
	}
	return points, nil
}

func (r *ClickEventRepository) GetClickCountByFolderID(ctx context.Context, folderID uuid.UUID) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.ClickEvent{}).
		Joins("JOIN links ON links.id = click_events.link_id").
		Where("links.folder_id = ? AND links.deleted_at IS NULL", folderID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ClickEventRepository) GetTotalClicks(ctx context.Context) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).Model(&model.ClickEvent{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
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
