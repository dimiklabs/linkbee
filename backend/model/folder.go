package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Folder struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Color     string         `gorm:"type:varchar(20);not null;default:'#635bff'" json:"color"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Folder) TableName() string {
	return "folders"
}
