package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

// Supported webhook event types.
const (
	EventLinkCreated = "link.created"
	EventLinkDeleted = "link.deleted"
	EventLinkClicked = "link.clicked"
)

// WebhookResponse is the outward-facing webhook representation (no secret).
type WebhookResponse struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WebhookServiceI defines operations for managing and firing webhooks.
type WebhookServiceI interface {
	List(ctx context.Context, userID uuid.UUID) ([]WebhookResponse, *dto.ServiceError)
	Create(ctx context.Context, userID uuid.UUID, url string, events []string) (*WebhookResponse, *dto.ServiceError)
	Update(ctx context.Context, id, userID uuid.UUID, url string, events []string, isActive bool) (*WebhookResponse, *dto.ServiceError)
	Delete(ctx context.Context, id, userID uuid.UUID) *dto.ServiceError
	// Trigger fires all active webhooks for userID that subscribe to event.
	// It is intentionally fire-and-forget — errors are logged, not returned.
	Trigger(userID uuid.UUID, event string, data any)
}

var deliveryClient = &http.Client{Timeout: 10 * time.Second}

type webhookService struct {
	webhookRepo repository.WebhookRepositoryI
}

func NewWebhookService(webhookRepo repository.WebhookRepositoryI) WebhookServiceI {
	return &webhookService{webhookRepo: webhookRepo}
}

func (s *webhookService) List(ctx context.Context, userID uuid.UUID) ([]WebhookResponse, *dto.ServiceError) {
	webhooks, err := s.webhookRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	out := make([]WebhookResponse, len(webhooks))
	for i, w := range webhooks {
		out[i] = toResponse(&w)
	}
	return out, nil
}

func (s *webhookService) Create(ctx context.Context, userID uuid.UUID, url string, events []string) (*WebhookResponse, *dto.ServiceError) {
	secret, err := generateSecret()
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	w := &model.Webhook{
		UserID:   userID,
		URL:      url,
		Secret:   secret,
		Events:   events,
		IsActive: true,
	}
	if err := s.webhookRepo.Create(ctx, w); err != nil {
		logger.ErrorCtx(ctx, "Failed to create webhook", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	resp := toResponse(w)
	resp.ID = w.ID // ensure ID is populated after creation
	return &resp, nil
}

func (s *webhookService) Update(ctx context.Context, id, userID uuid.UUID, url string, events []string, isActive bool) (*WebhookResponse, *dto.ServiceError) {
	w, err := s.webhookRepo.GetByID(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Webhook not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if url != "" {
		w.URL = url
	}
	if len(events) > 0 {
		w.Events = events
	}
	w.IsActive = isActive
	if err := s.webhookRepo.Update(ctx, w); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	resp := toResponse(w)
	return &resp, nil
}

func (s *webhookService) Delete(ctx context.Context, id, userID uuid.UUID) *dto.ServiceError {
	if err := s.webhookRepo.Delete(ctx, id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Webhook not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

// Trigger fetches matching webhooks and delivers the payload in background goroutines.
func (s *webhookService) Trigger(userID uuid.UUID, event string, data any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		webhooks, err := s.webhookRepo.GetActiveByEvent(ctx, userID, event)
		if err != nil {
			logger.Error("webhook trigger: failed to fetch webhooks",
				zap.String("user_id", userID.String()),
				zap.String("event", event),
				zap.Error(err))
			return
		}

		payload := map[string]any{
			"event":     event,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"data":      data,
		}
		body, err := json.Marshal(payload)
		if err != nil {
			logger.Error("webhook trigger: failed to marshal payload", zap.Error(err))
			return
		}

		for _, w := range webhooks {
			go deliver(w.URL, w.Secret, body, event)
		}
	}()
}

// deliver sends one webhook POST with HMAC-SHA256 signature.
func deliver(url, secret string, body []byte, event string) {
	sig := signPayload(secret, body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		logger.Error("webhook deliver: failed to build request",
			zap.String("url", url), zap.Error(err))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", event)
	req.Header.Set("X-Webhook-Signature", "sha256="+sig)

	resp, err := deliveryClient.Do(req)
	if err != nil {
		logger.Warn("webhook deliver: request failed",
			zap.String("url", url), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Info("webhook delivered",
			zap.String("url", url),
			zap.String("event", event),
			zap.Int("status", resp.StatusCode))
	} else {
		logger.Warn("webhook deliver: non-2xx response",
			zap.String("url", url),
			zap.String("event", event),
			zap.Int("status", resp.StatusCode))
	}
}

// signPayload computes HMAC-SHA256(secret, body) and returns a hex string.
func signPayload(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// generateSecret creates a random 32-byte hex secret (64 chars).
func generateSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generateSecret: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func toResponse(w *model.Webhook) WebhookResponse {
	events := make([]string, len(w.Events))
	copy(events, w.Events)
	return WebhookResponse{
		ID:        w.ID,
		URL:       w.URL,
		Events:    events,
		IsActive:  w.IsActive,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
	}
}
