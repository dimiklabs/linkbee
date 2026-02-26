package model

import (
	"time"

	"github.com/google/uuid"
)

// LinkVariant represents one destination URL in a split test.
// Weights are relative integers; the redirect service picks a variant
// proportionally (e.g., weight 60 vs 40 gives a 60/40 split).
type LinkVariant struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LinkID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"type:varchar(100);not null"`
	DestinationURL string    `gorm:"type:text;not null"`
	Weight         int       `gorm:"not null;default:50"`
	ClickCount     int64     `gorm:"default:0;not null"`
	CreatedAt      time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt      time.Time `gorm:"type:timestamptz;not null"`
}

func (LinkVariant) TableName() string {
	return "link_variants"
}
