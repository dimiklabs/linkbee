package request

// CreateLinkRequest is the request body for creating a new short link.
type CreateLinkRequest struct {
	DestinationURL string   `json:"destination_url" binding:"required,url"`
	Slug           string   `json:"slug" binding:"omitempty,min=3,max=20,alphanum"`
	Title          string   `json:"title" binding:"omitempty,max=500"`
	Password       string   `json:"password" binding:"omitempty,min=4,max=128"`
	ExpiresAt      *string  `json:"expires_at"` // RFC3339 string, nullable
	MaxClicks      *int64   `json:"max_clicks" binding:"omitempty,min=1"`
	RedirectType   *int16   `json:"redirect_type" binding:"omitempty,oneof=301 302"`
	Tags           []string `json:"tags" binding:"omitempty,max=10"`
	UTMSource      string   `json:"utm_source" binding:"omitempty,max=255"`
	UTMMedium      string   `json:"utm_medium" binding:"omitempty,max=255"`
	UTMCampaign    string   `json:"utm_campaign" binding:"omitempty,max=255"`
	FolderID       *string  `json:"folder_id"` // UUID string; null = no folder
}

// UpdateLinkRequest is the request body for updating an existing short link.
type UpdateLinkRequest struct {
	DestinationURL string   `json:"destination_url" binding:"omitempty,url"`
	Title          string   `json:"title" binding:"omitempty,max=500"`
	Password       string   `json:"password" binding:"omitempty,min=4,max=128"`
	ExpiresAt      *string  `json:"expires_at"` // RFC3339 string, nullable; null clears expiry
	MaxClicks      *int64   `json:"max_clicks" binding:"omitempty,min=1"`
	RedirectType   *int16   `json:"redirect_type" binding:"omitempty,oneof=301 302"`
	Tags           []string `json:"tags" binding:"omitempty,max=10"`
	IsActive       *bool    `json:"is_active"`
	UTMSource      string   `json:"utm_source" binding:"omitempty,max=255"`
	UTMMedium      string   `json:"utm_medium" binding:"omitempty,max=255"`
	UTMCampaign    string   `json:"utm_campaign" binding:"omitempty,max=255"`
	FolderID       *string  `json:"folder_id"` // UUID string; null = remove from folder
}

// ListLinksRequest holds query params for paginating the link list.
type ListLinksRequest struct {
	Page         int    `form:"page,default=1" binding:"min=1"`
	Limit        int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search       string `form:"search"`
	FolderID     string `form:"folder_id"`     // optional UUID filter
	Starred      *bool  `form:"starred"`        // optional; true = starred only
	HealthStatus string `form:"health_status"` // optional; healthy|unhealthy|timeout|error|unknown
}

// DemoShortenRequest is for the unauthenticated demo shorten endpoint.
type DemoShortenRequest struct {
	DestinationURL string `json:"destination_url" binding:"required,url"`
}

// AnalyticsRequest holds query params for analytics.
type AnalyticsRequest struct {
	From        string `form:"from"`        // RFC3339
	To          string `form:"to"`          // RFC3339
	Granularity string `form:"granularity"` // hour, day, week, month
}
