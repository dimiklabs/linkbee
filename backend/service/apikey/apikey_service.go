package apikey

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
)

const (
	keyPrefix    = "sl_"
	keyBodyBytes = 32 // 32 random bytes → 64 hex chars → total key ≈ 67 chars
	prefixLen    = 12 // first 12 chars stored for fast DB lookup
)

// APIKeyResponse is the safe public representation of an API key.
// KeyHash and full key are never included.
type APIKeyResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// CreateAPIKeyResponse wraps the response and includes the full key one time only.
type CreateAPIKeyResponse struct {
	APIKeyResponse
	// Key is the full plaintext key — shown to the user exactly once.
	Key string `json:"key"`
}

// ValidatedKey carries the user context resolved from a valid API key.
type ValidatedKey struct {
	KeyID  uuid.UUID
	UserID uuid.UUID
}

type APIKeyServiceI interface {
	Create(ctx context.Context, userID uuid.UUID, name string, expiresAt *time.Time) (*CreateAPIKeyResponse, *dto.ServiceError)
	List(ctx context.Context, userID uuid.UUID) ([]APIKeyResponse, *dto.ServiceError)
	Revoke(ctx context.Context, id, userID uuid.UUID) *dto.ServiceError
	Validate(ctx context.Context, rawKey string) (*ValidatedKey, *dto.ServiceError)
}

type apiKeyService struct {
	repo repository.APIKeyRepositoryI
}

func NewAPIKeyService(repo repository.APIKeyRepositoryI) APIKeyServiceI {
	return &apiKeyService{repo: repo}
}

func (s *apiKeyService) Create(ctx context.Context, userID uuid.UUID, name string, expiresAt *time.Time) (*CreateAPIKeyResponse, *dto.ServiceError) {
	rawKey, err := generateKey()
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "Failed to generate API key")
	}

	keyHash := hashKey(rawKey)
	prefix := rawKey[:prefixLen]

	record := &model.APIKey{
		UserID:    userID,
		Name:      name,
		KeyPrefix: prefix,
		KeyHash:   keyHash,
		ExpiresAt: expiresAt,
	}
	if err := s.repo.Create(ctx, record); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return &CreateAPIKeyResponse{
		APIKeyResponse: toResponse(record),
		Key:            rawKey,
	}, nil
}

func (s *apiKeyService) List(ctx context.Context, userID uuid.UUID) ([]APIKeyResponse, *dto.ServiceError) {
	records, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	out := make([]APIKeyResponse, 0, len(records))
	for _, r := range records {
		out = append(out, toResponse(&r))
	}
	return out, nil
}

func (s *apiKeyService) Revoke(ctx context.Context, id, userID uuid.UUID) *dto.ServiceError {
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeLinkNotFound, "API key not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

// Validate looks up the key by its prefix, compares the SHA-256 hash,
// checks expiry, and records the last-used timestamp asynchronously.
func (s *apiKeyService) Validate(ctx context.Context, rawKey string) (*ValidatedKey, *dto.ServiceError) {
	if len(rawKey) < prefixLen {
		return nil, dto.NewUnauthorizedError(constant.ErrCodeUnauthorized, "Invalid API key")
	}

	prefix := rawKey[:prefixLen]
	record, err := s.repo.GetByPrefix(ctx, prefix)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewUnauthorizedError(constant.ErrCodeUnauthorized, "Invalid API key")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Constant-time comparison via hash
	if hashKey(rawKey) != record.KeyHash {
		return nil, dto.NewUnauthorizedError(constant.ErrCodeUnauthorized, "Invalid API key")
	}

	// Check expiry
	if record.ExpiresAt != nil && time.Now().After(*record.ExpiresAt) {
		return nil, dto.NewUnauthorizedError(constant.ErrCodeUnauthorized, "API key has expired")
	}

	// Update last-used asynchronously so it doesn't add latency
	keyID := record.ID
	go func() {
		_ = s.repo.UpdateLastUsed(context.Background(), keyID)
	}()

	return &ValidatedKey{
		KeyID:  record.ID,
		UserID: record.UserID,
	}, nil
}

// generateKey creates a new key: "sl_" + hex(32 random bytes).
func generateKey() (string, error) {
	b := make([]byte, keyBodyBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return keyPrefix + hex.EncodeToString(b), nil
}

// hashKey returns the SHA-256 hex digest of the given key.
func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func toResponse(r *model.APIKey) APIKeyResponse {
	return APIKeyResponse{
		ID:         r.ID,
		Name:       r.Name,
		KeyPrefix:  r.KeyPrefix,
		LastUsedAt: r.LastUsedAt,
		ExpiresAt:  r.ExpiresAt,
		CreatedAt:  r.CreatedAt,
	}
}
