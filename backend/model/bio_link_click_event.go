package model

import (
	"time"

	"github.com/google/uuid"
)

// BioLinkClickEvent records a single click on a bio page link.
// The composite primary key (id, clicked_at) satisfies TimescaleDB's
// requirement that all UNIQUE/PRIMARY KEY constraints include the
// partition column (clicked_at).
type BioLinkClickEvent struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BioLinkID  uuid.UUID `gorm:"type:uuid;not null;index"                        json:"bio_link_id"`
	BioPageID  uuid.UUID `gorm:"type:uuid;not null;index"                        json:"bio_page_id"`
	ClickedAt  time.Time `gorm:"type:timestamptz;primaryKey;not null"            json:"clicked_at"`
	IPHash     string    `gorm:"type:varchar(64)"                                json:"ip_hash"`
	Country    string    `gorm:"type:char(2)"                                    json:"country,omitempty"`
	City       string    `gorm:"type:varchar(255)"                               json:"city,omitempty"`
	DeviceType string    `gorm:"type:varchar(50)"                                json:"device_type,omitempty"`
	OS         string    `gorm:"type:varchar(100)"                               json:"os,omitempty"`
	Browser    string    `gorm:"type:varchar(100)"                               json:"browser,omitempty"`
	Referrer   string    `gorm:"type:text"                                       json:"referrer,omitempty"`
}

func (BioLinkClickEvent) TableName() string { return "bio_link_click_events" }
