package model

import (
	"time"

	"github.com/google/uuid"
)

type BioLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BioPageID uuid.UUID `gorm:"type:uuid;not null;index" json:"bio_page_id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
	Position  int       `gorm:"not null;default:0" json:"position"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (BioLink) TableName() string { return "bio_links" }
