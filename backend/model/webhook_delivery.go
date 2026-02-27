package model

import (
	"time"

	"github.com/google/uuid"
)

type WebhookDelivery struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WebhookID    uuid.UUID `gorm:"type:uuid;not null;index" json:"webhook_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Event        string    `gorm:"type:varchar(100);not null" json:"event"`
	RequestBody  string    `gorm:"type:text" json:"request_body"`
	ResponseCode int       `gorm:"default:0" json:"response_code"`
	ResponseBody string    `gorm:"type:varchar(1024)" json:"response_body"`
	ErrorMessage string    `gorm:"type:text" json:"error_message,omitempty"`
	Success      bool      `gorm:"default:false;not null" json:"success"`
	DurationMs   int64     `gorm:"default:0" json:"duration_ms"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
}

func (WebhookDelivery) TableName() string {
	return "webhook_deliveries"
}
