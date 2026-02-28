package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
)

// DomainResponse is the API representation of a custom domain.
type DomainResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Domain      string    `json:"domain"`
	Status      string    `json:"status"`
	VerifyToken string    `json:"verify_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DomainServiceI interface {
	AddDomain(ctx context.Context, userID uuid.UUID, domain string) (*DomainResponse, *dto.ServiceError)
	ListDomains(ctx context.Context, userID uuid.UUID) ([]DomainResponse, *dto.ServiceError)
	VerifyDomain(ctx context.Context, id, userID uuid.UUID) (*DomainResponse, *dto.ServiceError)
	DeleteDomain(ctx context.Context, id, userID uuid.UUID) *dto.ServiceError
}

type domainService struct {
	repo repository.CustomDomainRepositoryI
}

func NewDomainService(repo repository.CustomDomainRepositoryI) DomainServiceI {
	return &domainService{repo: repo}
}

func (s *domainService) AddDomain(ctx context.Context, userID uuid.UUID, domain string) (*DomainResponse, *dto.ServiceError) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	if domain == "" {
		return nil, dto.NewBadRequestError(constant.ErrCodeBadRequest, "domain is required")
	}

	// Check if already registered (any user)
	existing, err := s.repo.GetByDomain(ctx, domain)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to check domain")
	}
	if existing != nil {
		return nil, dto.NewConflictError(constant.ErrCodeDomainAlreadyExists, "domain is already registered")
	}

	// Also check if this user already has a pending/failed entry for the same domain
	// by checking across all statuses (GetByDomain only checks verified)
	token, genErr := generateToken()
	if genErr != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to generate verification token")
	}

	now := time.Now()
	d := &model.CustomDomain{
		UserID:      userID,
		Domain:      domain,
		Status:      model.DomainStatusPending,
		VerifyToken: token,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, d); err != nil {
		// Unique constraint violation means already registered
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			return nil, dto.NewConflictError(constant.ErrCodeDomainAlreadyExists, "domain is already registered")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to add domain")
	}

	return toResponse(d), nil
}

func (s *domainService) ListDomains(ctx context.Context, userID uuid.UUID) ([]DomainResponse, *dto.ServiceError) {
	domains, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to list domains")
	}
	result := make([]DomainResponse, len(domains))
	for i := range domains {
		result[i] = *toResponse(&domains[i])
	}
	return result, nil
}

// VerifyDomain checks the DNS TXT record _shortlink-verify.<domain> for the stored token.
func (s *domainService) VerifyDomain(ctx context.Context, id, userID uuid.UUID) (*DomainResponse, *dto.ServiceError) {
	d, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeDomainNotFound, "domain not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch domain")
	}

	txtHost := "_shortlink-verify." + d.Domain
	verified := checkTXTRecord(txtHost, d.VerifyToken)

	newStatus := model.DomainStatusFailed
	if verified {
		newStatus = model.DomainStatusVerified
	}

	if updateErr := s.repo.UpdateStatus(ctx, id, newStatus); updateErr != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to update domain status")
	}
	d.Status = newStatus

	if !verified {
		return toResponse(d), dto.NewBadRequestError(constant.ErrCodeDomainVerifyFailed,
			"TXT record not found — add _shortlink-verify."+d.Domain+" = "+d.VerifyToken)
	}

	return toResponse(d), nil
}

func (s *domainService) DeleteDomain(ctx context.Context, id, userID uuid.UUID) *dto.ServiceError {
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeDomainNotFound, "domain not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to delete domain")
	}
	return nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func generateToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// checkTXTRecord looks up TXT records on host and returns true if any record
// exactly matches the expected token.
func checkTXTRecord(host, token string) bool {
	records, err := net.LookupTXT(host)
	if err != nil {
		return false
	}
	for _, r := range records {
		if strings.TrimSpace(r) == token {
			return true
		}
	}
	return false
}

func toResponse(d *model.CustomDomain) *DomainResponse {
	return &DomainResponse{
		ID:          d.ID.String(),
		UserID:      d.UserID.String(),
		Domain:      d.Domain,
		Status:      d.Status,
		VerifyToken: d.VerifyToken,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}
