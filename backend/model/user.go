package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email                string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password             string         `gorm:"type:varchar(255)" json:"-"`
	FirstName            string         `gorm:"type:varchar(100)" json:"first_name,omitempty"`
	LastName             string         `gorm:"type:varchar(100)" json:"last_name,omitempty"`
	Phone                string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	ProfilePicture       string         `gorm:"type:text" json:"profile_picture,omitempty"`
	ProfilePictureSource string         `gorm:"type:varchar(20)" json:"profile_picture_source,omitempty"`
	Status               string         `gorm:"type:varchar(20);default:'active';not null" json:"status"`
	Role                 string         `gorm:"type:varchar(50);default:'user';not null" json:"role"`
	AuthProvider         string         `gorm:"type:varchar(20);default:'local';not null" json:"auth_provider"`
	GoogleID             *string        `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	GitHubID             *string        `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	EmailVerified        bool           `gorm:"default:false;not null" json:"email_verified"`
	EmailVerifiedAt      *time.Time     `gorm:"type:timestamptz" json:"email_verified_at,omitempty"`
	LastLogin            *time.Time     `gorm:"type:timestamptz" json:"last_login,omitempty"`
	TotpSecret           string         `gorm:"type:varchar(255)" json:"-"`
	TotpEnabled          bool           `gorm:"default:false;not null" json:"totp_enabled"`
	ScheduledDeletionAt  *time.Time     `gorm:"type:timestamptz" json:"-"`
	CreatedAt            time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

const (
	AuthProviderLocal  = "local"
	AuthProviderGoogle = "google"
	AuthProviderGitHub = "github"
)

func (User) TableName() string {
	return "users"
}
