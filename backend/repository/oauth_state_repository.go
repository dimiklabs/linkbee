package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go/valkeycompat"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/logger"
)

const (
	oauthStatePrefix   = "oauth_state:"
	stateTTL           = 1 * time.Minute // OAuth state expires in 1 minute
	stateLength        = 32              // 256 bits of entropy
	codeVerifierLength = 43              // PKCE code verifier length (recommended minimum)
)

// OAuthStateData holds the state information for OAuth flow
type OAuthStateData struct {
	State        string    `json:"state"`
	CodeVerifier string    `json:"code_verifier"` // PKCE code verifier
	Nonce        string    `json:"nonce"`         // For ID token validation
	IPAddress    string    `json:"ip_address"`    // Client IP binding
	UserAgent    string    `json:"user_agent"`    // User agent binding
	RedirectURI  string    `json:"redirect_uri"`  // Where to redirect after auth
	CreatedAt    time.Time `json:"created_at"`
}

type OAuthStateRepositoryI interface {
	// Create generates and stores a new OAuth state
	Create(ctx context.Context, ipAddress, userAgent, redirectURI string) (*OAuthStateData, error)

	// Get retrieves and validates the OAuth state
	Get(ctx context.Context, state string) (*OAuthStateData, error)

	// Delete removes the OAuth state (should be called after use)
	Delete(ctx context.Context, state string) error

	// GenerateCodeChallenge generates the PKCE code challenge from the verifier
	GenerateCodeChallenge(codeVerifier string) string
}

type OAuthStateRepository struct {
	cache valkeycompat.Cmdable
}

func NewOAuthStateRepository(cache valkeycompat.Cmdable) OAuthStateRepositoryI {
	return &OAuthStateRepository{
		cache: cache,
	}
}

func (r *OAuthStateRepository) Create(ctx context.Context, ipAddress, userAgent, redirectURI string) (*OAuthStateData, error) {
	logger.DebugCtx(ctx, "Creating OAuth state")

	// Generate a cryptographically secure state token
	state, err := generateSecureToken(stateLength)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate state token", zap.Error(err))
		return nil, fmt.Errorf("failed to generate state token: %w", err)
	}

	// Generate PKCE code verifier
	codeVerifier, err := generateSecureToken(codeVerifierLength)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate code verifier", zap.Error(err))
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	// Generate nonce for ID token validation
	nonce, err := generateSecureToken(stateLength)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate nonce", zap.Error(err))
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	stateData := &OAuthStateData{
		State:        state,
		CodeVerifier: codeVerifier,
		Nonce:        nonce,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		RedirectURI:  redirectURI,
		CreatedAt:    time.Now(),
	}

	// Serialize state data
	data, err := json.Marshal(stateData)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to marshal state data", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal state data: %w", err)
	}

	// Store in Valkey with TTL
	key := fmt.Sprintf("%s%s", oauthStatePrefix, state)
	if err := r.cache.Set(ctx, key, data, stateTTL).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to store OAuth state", zap.Error(err))
		return nil, fmt.Errorf("failed to store OAuth state: %w", err)
	}

	logger.InfoCtx(ctx, "OAuth state created",
		zap.String("state", state[:8]+"...")) // Log only the first 8 chars for security

	return stateData, nil
}

func (r *OAuthStateRepository) Get(ctx context.Context, state string) (*OAuthStateData, error) {
	logger.DebugCtx(ctx, "Retrieving OAuth state")

	key := fmt.Sprintf("%s%s", oauthStatePrefix, state)
	data, err := r.cache.Get(ctx, key).Bytes()
	if err != nil {
		if err == valkeycompat.Nil {
			logger.WarnCtx(ctx, "OAuth state not found or expired")
			return nil, fmt.Errorf("OAuth state not found or expired")
		}
		logger.ErrorCtx(ctx, "Failed to retrieve OAuth state", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve OAuth state: %w", err)
	}

	var stateData OAuthStateData
	if err := json.Unmarshal(data, &stateData); err != nil {
		logger.ErrorCtx(ctx, "Failed to unmarshal state data", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal state data: %w", err)
	}

	return &stateData, nil
}

func (r *OAuthStateRepository) Delete(ctx context.Context, state string) error {
	logger.DebugCtx(ctx, "Deleting OAuth state")

	key := fmt.Sprintf("%s%s", oauthStatePrefix, state)
	if err := r.cache.Del(ctx, key).Err(); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete OAuth state", zap.Error(err))
		return fmt.Errorf("failed to delete OAuth state: %w", err)
	}

	logger.DebugCtx(ctx, "OAuth state deleted")
	return nil
}

// GenerateCodeChallenge generates a PKCE code challenge using SHA256
func (r *OAuthStateRepository) GenerateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// generateSecureToken generates a cryptographically secure URL-safe token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
