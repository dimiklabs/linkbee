package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	DomainStatusPending  = "pending"
	DomainStatusVerified = "verified"
	DomainStatusFailed   = "failed"
)

type CustomDomain struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Domain      string    `gorm:"type:varchar(253);not null;uniqueIndex"`
	Status      string    `gorm:"type:varchar(20);not null;default:'pending'"`
	VerifyToken string    `gorm:"type:varchar(64);not null"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null"`
}

func (CustomDomain) TableName() string {
	return "custom_domains"
}
