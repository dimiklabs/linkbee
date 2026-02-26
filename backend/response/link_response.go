package response

import (
	"time"

	"github.com/google/uuid"
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
	IsActive       bool       `json:"is_active"`
	Tags           []string   `json:"tags,omitempty"`
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
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	LinkID       uuid.UUID        `json:"link_id"`
	TotalClicks  int64            `json:"total_clicks"`
	UniqueClicks int64            `json:"unique_clicks"`
	TimeSeries   []TimeSeriesData `json:"time_series"`
	Referrers    []ReferrerData   `json:"referrers"`
	Devices      []DeviceData     `json:"devices"`
	Countries    []CountryData    `json:"countries"`
	Browsers     []BrowserData    `json:"browsers"`
	OSBreakdown  []OSData         `json:"os_breakdown"`
	Heatmap      []HeatmapData    `json:"heatmap"`
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

// DemoShortenResponse is returned from the demo shorten endpoint.
type DemoShortenResponse struct {
	ShortURL       string `json:"short_url"`
	Slug           string `json:"slug"`
	DestinationURL string `json:"destination_url"`
}
