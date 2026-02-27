package model

import (
	"time"

	"github.com/google/uuid"
)

type ClickEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LinkID      uuid.UUID `gorm:"type:uuid;not null;index" json:"link_id"`
	ClickedAt   time.Time `gorm:"type:timestamptz;not null;index" json:"clicked_at"`
	IPHash      string    `gorm:"type:varchar(64)" json:"ip_hash"`
	Country     string    `gorm:"type:char(2)" json:"country,omitempty"`
	City        string    `gorm:"type:varchar(255)" json:"city,omitempty"`
	DeviceType  string    `gorm:"type:varchar(50)" json:"device_type,omitempty"`
	OS          string    `gorm:"type:varchar(100)" json:"os,omitempty"`
	Browser     string    `gorm:"type:varchar(100)" json:"browser,omitempty"`
	Referrer    string    `gorm:"type:text" json:"referrer,omitempty"`
	Source      string    `gorm:"type:varchar(20)" json:"source,omitempty"`
	UTMSource   string    `gorm:"type:varchar(255)" json:"utm_source,omitempty"`
	UTMMedium   string    `gorm:"type:varchar(255)" json:"utm_medium,omitempty"`
	UTMCampaign string    `gorm:"type:varchar(255)" json:"utm_campaign,omitempty"`
	UTMContent  string    `gorm:"type:varchar(255)" json:"utm_content,omitempty"`
	UTMTerm     string    `gorm:"type:varchar(255)" json:"utm_term,omitempty"`
}

const (
	ClickSourceWeb = "web"
	ClickSourceQR  = "qr"
	ClickSourceAPI = "api"
)

func (ClickEvent) TableName() string {
	return "click_events"
}
