package model

import (
	"time"

	"github.com/google/uuid"
)

// LinkGeoRule routes visitors from a specific country to a custom destination URL.
type LinkGeoRule struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LinkID         uuid.UUID `gorm:"type:uuid;not null;index" json:"link_id"`
	CountryCode    string    `gorm:"type:char(2);not null" json:"country_code"` // ISO 3166-1 alpha-2
	DestinationURL string    `gorm:"type:text;not null" json:"destination_url"`
	Priority       int       `gorm:"not null;default:0" json:"priority"` // lower = higher priority
	CreatedAt      time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (LinkGeoRule) TableName() string {
	return "link_geo_rules"
}
