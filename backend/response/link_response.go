package response

import (
	"time"

	"github.com/google/uuid"
)

// HealthStatus values for a link.
const (
	HealthStatusUnknown   = "unknown"
	HealthStatusHealthy   = "healthy"
	HealthStatusUnhealthy = "unhealthy"
	HealthStatusTimeout   = "timeout"
	HealthStatusError     = "error"
)

// LinkResponse is the standard link representation returned to clients.
type LinkResponse struct {
	ID             uuid.UUID  `json:"id"`
	FolderID       *uuid.UUID `json:"folder_id,omitempty"`
	Slug           string     `json:"slug"`
	ShortURL       string     `json:"short_url"`
	DestinationURL string     `json:"destination_url"`
	Title          string     `json:"title,omitempty"`
	ClickCount     int64      `json:"click_count"`
	RedirectType   int16      `json:"redirect_type"`
	IsActive           bool       `json:"is_active"`
	IsStarred          bool       `json:"is_starred"`
	IsSplitTest        bool       `json:"is_split_test"`
	IsGeoRouting       bool       `json:"is_geo_routing"`
	IsPixelTracking    bool       `json:"is_pixel_tracking"`
	HealthStatus       string     `json:"health_status"`
	HealthStatusCode   int        `json:"health_status_code,omitempty"`
	HealthCheckedAt    *time.Time `json:"health_checked_at,omitempty"`
	Tags               []string   `json:"tags,omitempty"`
	HasPassword    bool       `json:"has_password"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	MaxClicks      *int64     `json:"max_clicks,omitempty"`
	UTMSource      string     `json:"utm_source,omitempty"`
	UTMMedium      string     `json:"utm_medium,omitempty"`
	UTMCampaign    string     `json:"utm_campaign,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// FolderResponse represents a link folder.
type FolderResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Color      string    `json:"color"`
	ClickCount int64     `json:"click_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BulkActionResponse reports how many links were affected by a bulk operation.
type BulkActionResponse struct {
	Affected int64  `json:"affected"`
	Action   string `json:"action"`
}

// LinkListResponse wraps a paginated list of links.
type LinkListResponse struct {
	Links      []*LinkResponse `json:"links"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// AnalyticsResponse is the analytics summary for a single link.
type AnalyticsResponse struct {
	LinkID             uuid.UUID        `json:"link_id"`
	TotalClicks        int64            `json:"total_clicks"`
	UniqueClicks       int64            `json:"unique_clicks"`
	FirstTimeVisitors  int64            `json:"first_time_visitors"`
	ReturningVisitors  int64            `json:"returning_visitors"`
	TimeSeries   []TimeSeriesData `json:"time_series"`
	Referrers    []ReferrerData   `json:"referrers"`
	Devices      []DeviceData     `json:"devices"`
	Countries    []CountryData    `json:"countries"`
	Browsers     []BrowserData    `json:"browsers"`
	OSBreakdown  []OSData         `json:"os_breakdown"`
	Sources      []SourceData     `json:"sources"`
	Heatmap      []HeatmapData    `json:"heatmap"`
	UTMSources         []UTMData             `json:"utm_sources"`
	UTMMediums         []UTMData             `json:"utm_mediums"`
	UTMCampaigns       []UTMData             `json:"utm_campaigns"`
	UTMContents        []UTMData             `json:"utm_contents"`
	UTMTerms           []UTMData             `json:"utm_terms"`
	Cities             []CityData            `json:"cities"`
	ReferrerCategories []ReferrerCategoryData `json:"referrer_categories"`
}

// SourceData holds click count broken down by source (web, qr, api).
type SourceData struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

// UTMData holds click count for a single UTM parameter value.
type UTMData struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type TimeSeriesData struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

type ReferrerData struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type DeviceData struct {
	DeviceType string `json:"device_type"`
	Count      int64  `json:"count"`
}

type CountryData struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type CityData struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

// ReferrerCategoryData holds click count for a single traffic channel category.
type ReferrerCategoryData struct {
	Category string `json:"category"` // direct | search | social | email | referral
	Count    int64  `json:"count"`
}

type BrowserData struct {
	Browser string `json:"browser"`
	Count   int64  `json:"count"`
}

type OSData struct {
	OS    string `json:"os"`
	Count int64  `json:"count"`
}

type HeatmapData struct {
	DayOfWeek int   `json:"day_of_week"` // 0 = Sunday … 6 = Saturday
	Hour      int   `json:"hour"`        // 0–23 UTC
	Count     int64 `json:"count"`
}

// PeriodMetrics holds aggregated stats for a single time window.
type PeriodMetrics struct {
	From         string `json:"from"`
	To           string `json:"to"`
	TotalClicks  int64  `json:"total_clicks"`
	UniqueClicks int64  `json:"unique_clicks"`
}

// PeriodChange describes the delta between current and previous periods.
type PeriodChange struct {
	CountChange   int64   `json:"count_change"`
	PercentChange float64 `json:"percent_change"` // positive = increase; 0 when no previous data
	Trend         string  `json:"trend"`           // "up" | "down" | "stable"
}

// PeriodComparisonResponse is returned by GET /links/:id/analytics/comparison.
type PeriodComparisonResponse struct {
	LinkID       uuid.UUID     `json:"link_id"`
	Current      PeriodMetrics `json:"current"`
	Previous     PeriodMetrics `json:"previous"`
	Clicks       PeriodChange  `json:"clicks"`
	UniqueClicks PeriodChange  `json:"unique_clicks"`
}

// LinkVariantResponse represents a single A/B split-test variant.
type LinkVariantResponse struct {
	ID             uuid.UUID `json:"id"`
	LinkID         uuid.UUID `json:"link_id"`
	Name           string    `json:"name"`
	DestinationURL string    `json:"destination_url"`
	Weight         int       `json:"weight"`
	ClickCount     int64     `json:"click_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ImportLinkError captures a per-row failure during CSV import.
type ImportLinkError struct {
	Row   int    `json:"row"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

// ImportLinksResponse is returned from the bulk CSV import endpoint.
type ImportLinksResponse struct {
	Total   int               `json:"total"`
	Created int               `json:"created"`
	Failed  int               `json:"failed"`
	Errors  []ImportLinkError `json:"errors,omitempty"`
}

// LinkComparisonMetric holds per-link analytics used in the multi-link comparison.
type LinkComparisonMetric struct {
	LinkID       uuid.UUID `json:"link_id"`
	Slug         string    `json:"slug"`
	ShortURL     string    `json:"short_url"`
	Title        string    `json:"title,omitempty"`
	TotalClicks  int64     `json:"total_clicks"`
	UniqueClicks int64     `json:"unique_clicks"`
	ClicksPerDay float64   `json:"clicks_per_day"`
	TopReferrer  string    `json:"top_referrer,omitempty"`
	TopCountry   string    `json:"top_country,omitempty"`
	TopBrowser   string    `json:"top_browser,omitempty"`
	TopDevice    string    `json:"top_device,omitempty"`
}

// MultiLinkComparisonResponse is returned by GET /api/v1/links/comparison.
type MultiLinkComparisonResponse struct {
	Links    []LinkComparisonMetric `json:"links"`
	From     string                 `json:"from"`
	To       string                 `json:"to"`
	SpanDays int64                  `json:"span_days"`
}

// DemoShortenResponse is returned from the demo shorten endpoint.
type DemoShortenResponse struct {
	ShortURL       string `json:"short_url"`
	Slug           string `json:"slug"`
	DestinationURL string `json:"destination_url"`
}
