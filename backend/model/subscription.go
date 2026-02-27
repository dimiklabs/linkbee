package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	SubStatusActive    = "active"
	SubStatusCancelled = "cancelled"
	SubStatusExpired   = "expired"
	SubStatusPastDue   = "past_due"
	SubStatusPaused    = "paused"
	SubStatusTrialing  = "on_trial"
)

type Subscription struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID                 uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	PlanID                 string     `gorm:"type:varchar(20);not null;default:'free'" json:"plan_id"`
	Status                 string     `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	LemonSqueezySubID      string     `gorm:"type:varchar(100)" json:"lemon_squeezy_sub_id,omitempty"`
	LemonSqueezyOrderID    string     `gorm:"type:varchar(100)" json:"lemon_squeezy_order_id,omitempty"`
	LemonSqueezyCustomerID string     `gorm:"type:varchar(100)" json:"lemon_squeezy_customer_id,omitempty"`
	LemonSqueezyVariantID  string     `gorm:"type:varchar(100)" json:"lemon_squeezy_variant_id,omitempty"`
	CurrentPeriodEnd       *time.Time `gorm:"type:timestamptz" json:"current_period_end,omitempty"`
	CancelledAt            *time.Time `gorm:"type:timestamptz" json:"cancelled_at,omitempty"`
	CreatedAt              time.Time  `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt              time.Time  `gorm:"type:timestamptz;not null" json:"updated_at"`
}

func (Subscription) TableName() string { return "subscriptions" }
