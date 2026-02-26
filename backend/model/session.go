package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	RefreshTokenJTI string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"-"`
	UserAgent       string    `gorm:"type:varchar(512)" json:"user_agent"`
	IPAddress       string    `gorm:"type:varchar(45)" json:"ip_address"`
	CreatedAt       time.Time `gorm:"type:timestamptz;not null" json:"created_at"`
	ExpiresAt       time.Time `gorm:"type:timestamptz;not null" json:"expires_at"`
	LastActivityAt  time.Time `gorm:"type:timestamptz;not null" json:"last_activity_at"`

	// Enhanced session fields
	RememberMe     bool       `gorm:"default:false;not null" json:"remember_me"`
	DeviceName     string     `gorm:"type:varchar(255)" json:"device_name"`
	DeviceType     string     `gorm:"type:varchar(50)" json:"device_type"`
	Browser        string     `gorm:"type:varchar(100)" json:"browser"`
	OS             string     `gorm:"type:varchar(100)" json:"os"`
	Location       string     `gorm:"type:varchar(255)" json:"location,omitempty"`
	LoginMethod    string     `gorm:"type:varchar(50)" json:"login_method"` // local, google, github, etc.
	LastActivityIP string     `gorm:"type:varchar(45)" json:"last_activity_ip"`
	ActivityCount  int64      `gorm:"default:0" json:"activity_count"`
	NotifiedAt     *time.Time `gorm:"type:timestamptz" json:"notified_at,omitempty"` // When user was notified of this new session
}

const (
	LoginMethodLocal    = "local"
	LoginMethodGoogle   = "google"
	LoginMethodGitHub   = "github"
	LoginMethodFacebook = "facebook"
)

func (Session) TableName() string {
	return "sessions"
}
