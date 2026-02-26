package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey represents a user-issued API key used to authenticate programmatic requests.
// The full key is shown to the user exactly once at creation and never stored in plaintext.
type APIKey struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	KeyPrefix  string         `gorm:"type:varchar(16);not null;uniqueIndex" json:"key_prefix"` // first N chars for fast lookup
	KeyHash    string         `gorm:"type:varchar(64);not null" json:"-"`                       // SHA-256 hex of full key
	LastUsedAt *time.Time     `gorm:"type:timestamptz" json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time     `gorm:"type:timestamptz" json:"expires_at,omitempty"`
	CreatedAt  time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (APIKey) TableName() string {
	return "api_keys"
}
