package model

import (
	"time"

	"github.com/google/uuid"
)

type TotpBackupCode struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CodeHash  string    `gorm:"type:varchar(255);not null" json:"-"`
	Used      bool      `gorm:"default:false;not null" json:"used"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
}

func (TotpBackupCode) TableName() string {
	return "totp_backup_codes"
}
