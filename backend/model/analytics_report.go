package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// AnalyticsReport defines a scheduled report for one or more links.
type AnalyticsReport struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	LinkIDs   pq.StringArray `gorm:"type:text[];not null" json:"link_ids"`
	Frequency string         `gorm:"type:varchar(50);not null" json:"frequency"` // daily | weekly | monthly
	NextRunAt *time.Time     `gorm:"type:timestamptz" json:"next_run_at,omitempty"`
	IsActive  bool           `gorm:"default:true;not null" json:"is_active"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AnalyticsReport) TableName() string { return "analytics_reports" }

// ReportDelivery records each attempt to send a scheduled report.
type ReportDelivery struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ReportID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"report_id"`
	Status        string     `gorm:"type:varchar(50);not null" json:"status"` // sent | failed
	FailureReason string     `gorm:"type:text" json:"failure_reason,omitempty"`
	DeliveredAt   *time.Time `gorm:"type:timestamptz" json:"delivered_at,omitempty"`
	CreatedAt     time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
}

func (ReportDelivery) TableName() string { return "report_deliveries" }
