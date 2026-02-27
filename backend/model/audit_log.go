package model

import (
	"time"

	"github.com/google/uuid"
)

// Audit action constants — keep these in sync with the frontend filter list.
const (
	AuditActionUserSignup      = "user_signup"
	AuditActionUserLogin       = "user_login"
	AuditActionUserLogout      = "user_logout"
	AuditActionPasswordChanged = "password_changed"
	AuditActionAccountDeleted  = "account_deleted"

	AuditActionLinkCreated  = "link_created"
	AuditActionLinkUpdated  = "link_updated"
	AuditActionLinkDeleted  = "link_deleted"
	AuditActionLinksImported = "links_imported"

	AuditActionDomainAdded    = "domain_added"
	AuditActionDomainVerified = "domain_verified"
	AuditActionDomainDeleted  = "domain_deleted"

	AuditActionAPIKeyCreated = "api_key_created"
	AuditActionAPIKeyRevoked = "api_key_revoked"
)

// Audit resource type constants.
const (
	AuditResourceUser   = "user"
	AuditResourceLink   = "link"
	AuditResourceDomain = "domain"
	AuditResourceAPIKey = "api_key"
)

type AuditLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Action       string    `gorm:"type:varchar(64);not null;index"`
	ResourceType string    `gorm:"type:varchar(64);not null"`
	ResourceID   string    `gorm:"type:varchar(255)"`
	ResourceName string    `gorm:"type:varchar(500)"`
	IPAddress    string    `gorm:"type:varchar(64)"`
	UserAgent    string    `gorm:"type:varchar(512)"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;index"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
