package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	AvatarURL   string         `gorm:"type:text" json:"avatar_url,omitempty"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Team) TableName() string {
	return "teams"
}

// TeamRole values: owner, admin, member
type TeamMember struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TeamID      uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_team_user" json:"team_id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index;uniqueIndex:idx_team_user" json:"user_id"`
	Role        string         `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	InvitedBy   uuid.UUID      `gorm:"type:uuid" json:"invited_by"`
	JoinedAt    *time.Time     `gorm:"type:timestamptz" json:"joined_at,omitempty"`
	InviteEmail string         `gorm:"type:varchar(255)" json:"invite_email,omitempty"`
	InviteToken string         `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	Status      string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, active, declined
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TeamMember) TableName() string {
	return "team_members"
}
