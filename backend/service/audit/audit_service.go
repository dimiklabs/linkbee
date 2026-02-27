package audit

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

// LogEntry holds all data for a single audit event.
type LogEntry struct {
	UserID       uuid.UUID
	Action       string
	ResourceType string
	ResourceID   string // UUID string or human-readable identifier
	ResourceName string // slug, domain name, key name, etc.
	IPAddress    string
	UserAgent    string
}

type AuditServiceI interface {
	// LogAsync records an audit event in a background goroutine (fire-and-forget).
	LogAsync(entry LogEntry)
}

type auditService struct {
	repo repository.AuditLogRepositoryI
}

func NewAuditService(repo repository.AuditLogRepositoryI) AuditServiceI {
	return &auditService{repo: repo}
}

func (s *auditService) LogAsync(entry LogEntry) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log := &model.AuditLog{
			UserID:       entry.UserID,
			Action:       entry.Action,
			ResourceType: entry.ResourceType,
			ResourceID:   entry.ResourceID,
			ResourceName: entry.ResourceName,
			IPAddress:    entry.IPAddress,
			UserAgent:    entry.UserAgent,
			CreatedAt:    time.Now(),
		}
		_ = s.repo.Create(ctx, log)
	}()
}
