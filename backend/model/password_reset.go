package model

import (
	"time"

	"github.com/google/uuid"
)

type PasswordReset struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	Used      bool      `gorm:"default:false;not null" json:"used"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
