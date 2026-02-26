package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BioPage struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Username    string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Title       string         `gorm:"type:varchar(100)" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	AvatarURL   string         `gorm:"type:text" json:"avatar_url"`
	Theme       string         `gorm:"type:varchar(20);not null;default:'light'" json:"theme"`
	IsPublished bool           `gorm:"default:false;not null" json:"is_published"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Links       []BioLink      `gorm:"foreignKey:BioPageID;references:ID" json:"links,omitempty"`
}

func (BioPage) TableName() string { return "bio_pages" }
