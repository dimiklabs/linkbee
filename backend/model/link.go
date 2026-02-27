package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Link struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	FolderID       *uuid.UUID     `gorm:"type:uuid;index" json:"folder_id,omitempty"`
	Slug           string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"slug"`
	DestinationURL string         `gorm:"type:text;not null" json:"destination_url"`
	Title          string         `gorm:"type:varchar(500)" json:"title,omitempty"`
	PasswordHash   string         `gorm:"type:varchar(255)" json:"-"`
	ExpiresAt      *time.Time     `gorm:"type:timestamptz" json:"expires_at,omitempty"`
	MaxClicks      *int64         `gorm:"type:bigint" json:"max_clicks,omitempty"`
	ClickCount     int64          `gorm:"default:0;not null" json:"click_count"`
	RedirectType   int16          `gorm:"default:302;not null" json:"redirect_type"`
	IsActive           bool           `gorm:"default:true;not null" json:"is_active"`
	IsStarred          bool           `gorm:"default:false;not null" json:"is_starred"`
	IsSplitTest        bool           `gorm:"default:false;not null" json:"is_split_test"`
	IsGeoRouting       bool           `gorm:"default:false;not null" json:"is_geo_routing"`
	IsPixelTracking    bool           `gorm:"default:false;not null" json:"is_pixel_tracking"`
	HealthStatus       string         `gorm:"type:varchar(20);not null;default:'unknown'" json:"health_status"`
	HealthStatusCode   int            `gorm:"default:0;not null" json:"health_status_code"`
	HealthCheckedAt    *time.Time     `gorm:"type:timestamptz" json:"health_checked_at,omitempty"`
	ExpiryNotifiedAt   *time.Time     `gorm:"type:timestamptz" json:"-"`
	Tags           pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	UTMSource      string         `gorm:"type:varchar(255)" json:"utm_source,omitempty"`
	UTMMedium      string         `gorm:"type:varchar(255)" json:"utm_medium,omitempty"`
	UTMCampaign    string         `gorm:"type:varchar(255)" json:"utm_campaign,omitempty"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Link) TableName() string {
	return "links"
}
