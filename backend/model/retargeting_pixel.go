package model

import (
	"time"

	"github.com/google/uuid"
)

// Supported pixel types.
const (
	PixelTypeFacebook  = "facebook"
	PixelTypeGoogleAds = "google_ads"
	PixelTypeTikTok    = "tiktok"
	PixelTypeLinkedIn  = "linkedin"
	PixelTypeCustom    = "custom"
)

// RetargetingPixel stores a tracking pixel configuration for a link.
// When the link has IsPixelTracking=true, the redirect handler serves an
// intermediate HTML page that fires all active pixels before redirecting.
type RetargetingPixel struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LinkID       uuid.UUID `gorm:"type:uuid;not null;index" json:"link_id"`
	PixelType    string    `gorm:"type:varchar(20);not null" json:"pixel_type"`
	PixelID      string    `gorm:"type:varchar(200)" json:"pixel_id,omitempty"`
	CustomScript string    `gorm:"type:text" json:"custom_script,omitempty"`
	IsActive     bool      `gorm:"default:true;not null" json:"is_active"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (RetargetingPixel) TableName() string { return "retargeting_pixels" }
