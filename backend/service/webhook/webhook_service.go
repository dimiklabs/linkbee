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
	"io"
	"net/http"
	"time"
	"unicode/utf8"

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

type DeliveriesResponse struct {
	Deliveries []*model.WebhookDelivery `json:"deliveries"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
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
	// Delivery history & operations
	GetDeliveries(ctx context.Context, webhookID, userID uuid.UUID, page, limit int) (*DeliveriesResponse, *dto.ServiceError)
	ResendDelivery(ctx context.Context, deliveryID, userID uuid.UUID) (*model.WebhookDelivery, *dto.ServiceError)
	GetWebhookSecret(ctx context.Context, webhookID, userID uuid.UUID) (string, *dto.ServiceError)
	TestWebhook(ctx context.Context, webhookID, userID uuid.UUID) (*model.WebhookDelivery, *dto.ServiceError)
}

var deliveryClient = &http.Client{Timeout: 10 * time.Second}

const (
	maxRequestBodyLog  = 4 * 1024 // 4 KB stored in DB
	maxResponseBodyLog = 1024      // 1 KB
)

type webhookService struct {
	webhookRepo  repository.WebhookRepositoryI
	deliveryRepo repository.WebhookDeliveryRepositoryI
}

func NewWebhookService(webhookRepo repository.WebhookRepositoryI, deliveryRepo repository.WebhookDeliveryRepositoryI) WebhookServiceI {
	return &webhookService{webhookRepo: webhookRepo, deliveryRepo: deliveryRepo}
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
			wCopy := w
			go s.deliverAndLog(context.Background(), wCopy.ID, userID, wCopy.URL, wCopy.Secret, body, event)
		}
	}()
}

// deliverAndLog sends one webhook POST and persists the delivery result.
func (s *webhookService) deliverAndLog(ctx context.Context, webhookID, userID uuid.UUID, url, secret string, body []byte, event string) *model.WebhookDelivery {
	sig := signPayload(secret, body)
	start := time.Now()

	delivery := &model.WebhookDelivery{
		WebhookID:   webhookID,
		UserID:      userID,
		Event:       event,
		RequestBody: truncateString(string(body), maxRequestBodyLog),
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		delivery.ErrorMessage = err.Error()
		delivery.Success = false
		_ = s.deliveryRepo.Create(ctx, delivery)
		logger.Error("webhook deliver: failed to build request",
			zap.String("url", url), zap.Error(err))
		return delivery
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Event", event)
	req.Header.Set("X-Webhook-Signature", "sha256="+sig)

	resp, err := deliveryClient.Do(req)
	delivery.DurationMs = time.Since(start).Milliseconds()

	if err != nil {
		delivery.ErrorMessage = err.Error()
		delivery.Success = false
		_ = s.deliveryRepo.Create(ctx, delivery)
		logger.Warn("webhook deliver: request failed",
			zap.String("url", url), zap.Error(err))
		return delivery
	}
	defer resp.Body.Close()

	delivery.ResponseCode = resp.StatusCode
	delivery.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	if respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodyLog)); err == nil {
		delivery.ResponseBody = string(respBody)
	}

	if delivery.Success {
		logger.Info("webhook delivered",
			zap.String("url", url),
			zap.String("event", event),
			zap.Int("status", resp.StatusCode),
			zap.Int64("duration_ms", delivery.DurationMs))
	} else {
		delivery.ErrorMessage = fmt.Sprintf("HTTP %d", resp.StatusCode)
		logger.Warn("webhook deliver: non-2xx response",
			zap.String("url", url),
			zap.String("event", event),
			zap.Int("status", resp.StatusCode))
	}

	_ = s.deliveryRepo.Create(ctx, delivery)
	return delivery
}

func (s *webhookService) GetDeliveries(ctx context.Context, webhookID, userID uuid.UUID, page, limit int) (*DeliveriesResponse, *dto.ServiceError) {
	// Verify webhook ownership
	if _, err := s.webhookRepo.GetByID(ctx, webhookID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Webhook not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	deliveries, total, err := s.deliveryRepo.ListByWebhookID(ctx, webhookID, userID, page, limit)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return &DeliveriesResponse{
		Deliveries: deliveries,
		Total:      total,
		Page:       page,
		Limit:      limit,
	}, nil
}

func (s *webhookService) ResendDelivery(ctx context.Context, deliveryID, userID uuid.UUID) (*model.WebhookDelivery, *dto.ServiceError) {
	original, err := s.deliveryRepo.GetByID(ctx, deliveryID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Delivery not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	webhook, err := s.webhookRepo.GetByID(ctx, original.WebhookID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Webhook not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	delivery := s.deliverAndLog(ctx, webhook.ID, userID, webhook.URL, webhook.Secret,
		[]byte(original.RequestBody), original.Event)
	return delivery, nil
}

func (s *webhookService) GetWebhookSecret(ctx context.Context, webhookID, userID uuid.UUID) (string, *dto.ServiceError) {
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", dto.NewNotFoundError(constant.ErrCodeNotFound, "Webhook not found")
		}
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return webhook.Secret, nil
}

func (s *webhookService) TestWebhook(ctx context.Context, webhookID, userID uuid.UUID) (*model.WebhookDelivery, *dto.ServiceError) {
	webhook, err := s.webhookRepo.GetByID(ctx, webhookID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Webhook not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	payload := map[string]any{
		"event":     "test",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data": map[string]any{
			"message": "This is a test delivery from Shortlink.",
		},
	}
	body, _ := json.Marshal(payload)

	delivery := s.deliverAndLog(ctx, webhook.ID, userID, webhook.URL, webhook.Secret, body, "test")
	return delivery, nil
}

// truncateString trims s to at most maxBytes (in bytes, not runes) preserving UTF-8 validity.
func truncateString(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	for !utf8.ValidString(s[:maxBytes]) {
		maxBytes--
	}
	return s[:maxBytes]
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
