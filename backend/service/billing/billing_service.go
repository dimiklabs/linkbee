package billing

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
)

// ─── Response types ───────────────────────────────────────────────────────────

type SubscriptionResponse struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	PlanID           string     `json:"plan_id"`
	Status           string     `json:"status"`
	CurrentPeriodEnd *time.Time `json:"current_period_end,omitempty"`
	CancelledAt      *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type SubscriptionWithPlanResponse struct {
	Subscription SubscriptionResponse `json:"subscription"`
	Plan         PlanInfo             `json:"plan"`
}

// ─── Paddle API types ─────────────────────────────────────────────────────────

type paddleTransactionRequest struct {
	Items      []paddleItem      `json:"items"`
	CustomData map[string]string `json:"custom_data"`
}

type paddleItem struct {
	PriceID  string `json:"price_id"`
	Quantity int    `json:"quantity"`
}

type paddleTransactionResponse struct {
	Data struct {
		Checkout struct {
			URL string `json:"url"`
		} `json:"checkout"`
	} `json:"data"`
	Error *struct {
		Type   string `json:"type"`
		Code   string `json:"code"`
		Detail string `json:"detail"`
	} `json:"error,omitempty"`
}

// ─── Paddle webhook types ─────────────────────────────────────────────────────

type paddleWebhookPayload struct {
	EventType string        `json:"event_type"`
	Data      paddleSubData `json:"data"`
}

type paddleSubData struct {
	ID                   string            `json:"id"`
	Status               string            `json:"status"`
	CustomerID           string            `json:"customer_id"`
	Items                []paddleSubItem   `json:"items"`
	CustomData           map[string]string `json:"custom_data"`
	CurrentBillingPeriod *struct {
		EndsAt string `json:"ends_at"`
	} `json:"current_billing_period"`
	NextBilledAt *string `json:"next_billed_at"`
	CanceledAt   *string `json:"canceled_at"`
	PausedAt     *string `json:"paused_at"`
}

type paddleSubItem struct {
	Price struct {
		ID string `json:"id"`
	} `json:"price"`
	Quantity int `json:"quantity"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type BillingServiceI interface {
	GetSubscription(ctx context.Context, userID uuid.UUID) (*SubscriptionWithPlanResponse, *dto.ServiceError)
	GetCheckoutURL(planID string, userID uuid.UUID) (string, *dto.ServiceError)
	HandleWebhook(ctx context.Context, body []byte, signature string) *dto.ServiceError
	VerifyWebhookSignature(body []byte, signature string) bool
}

type billingService struct {
	subRepo    repository.SubscriptionRepositoryI
	billingCfg *config.BillingConfig
}

func NewBillingService(subRepo repository.SubscriptionRepositoryI, billingCfg *config.BillingConfig) BillingServiceI {
	return &billingService{subRepo: subRepo, billingCfg: billingCfg}
}

// GetSubscription returns the user's subscription and plan details.
// If no subscription record exists, a synthetic free-plan subscription is returned.
func (s *billingService) GetSubscription(ctx context.Context, userID uuid.UUID) (*SubscriptionWithPlanResponse, *dto.ServiceError) {
	sub, err := s.subRepo.GetByUserID(ctx, userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to fetch subscription")
	}

	if err == gorm.ErrRecordNotFound || sub == nil {
		now := time.Now()
		sub = &model.Subscription{
			ID:        uuid.Nil,
			UserID:    userID,
			PlanID:    PlanFree,
			Status:    model.SubStatusActive,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	plan := GetPlan(sub.PlanID)
	return &SubscriptionWithPlanResponse{
		Subscription: toResponse(sub),
		Plan:         plan,
	}, nil
}

// GetCheckoutURL creates a Paddle transaction and returns the hosted checkout URL.
func (s *billingService) GetCheckoutURL(planID string, userID uuid.UUID) (string, *dto.ServiceError) {
	var priceID string
	switch planID {
	case PlanPro:
		priceID = s.billingCfg.PaddleProPriceID
	case PlanGrowth:
		priceID = s.billingCfg.PaddleGrowthPriceID
	default:
		return "", dto.NewBadRequestError(constant.ErrCodeBadRequest, "invalid plan: must be pro or growth")
	}

	if priceID == "" {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "checkout not configured for this plan")
	}

	if s.billingCfg.PaddleAPIKey == "" {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "payment provider not configured")
	}

	reqBody := paddleTransactionRequest{
		Items:      []paddleItem{{PriceID: priceID, Quantity: 1}},
		CustomData: map[string]string{"user_id": userID.String()},
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "failed to build checkout request")
	}

	apiBase := "https://api.paddle.com"
	if s.billingCfg.PaddleEnvironment == "sandbox" {
		apiBase = "https://sandbox-api.paddle.com"
	}

	req, err := http.NewRequest(http.MethodPost, apiBase+"/transactions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "failed to create checkout request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.billingCfg.PaddleAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "failed to reach payment provider")
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var paddleResp paddleTransactionResponse
	if err := json.Unmarshal(respBytes, &paddleResp); err != nil {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "unexpected response from payment provider")
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		detail := "checkout creation failed"
		if paddleResp.Error != nil {
			detail = paddleResp.Error.Detail
		}
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, detail)
	}

	if paddleResp.Data.Checkout.URL == "" {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "no checkout URL returned by payment provider")
	}
	return paddleResp.Data.Checkout.URL, nil
}

// VerifyWebhookSignature validates a Paddle webhook signature.
// Paddle-Signature header format: ts=TIMESTAMP;h1=HMAC_HEX
// Verification: HMAC-SHA256(secret, "ts:body")
func (s *billingService) VerifyWebhookSignature(body []byte, signature string) bool {
	if s.billingCfg.PaddleWebhookSecret == "" {
		return true // skip verification in dev mode
	}

	var ts, h1 string
	for _, part := range strings.Split(signature, ";") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "ts":
			ts = kv[1]
		case "h1":
			h1 = kv[1]
		}
	}
	if ts == "" || h1 == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(s.billingCfg.PaddleWebhookSecret))
	mac.Write([]byte(ts + ":"))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(h1))
}

// HandleWebhook processes a verified Paddle webhook event.
func (s *billingService) HandleWebhook(ctx context.Context, body []byte, signature string) *dto.ServiceError {
	if !s.VerifyWebhookSignature(body, signature) {
		return dto.NewUnauthorizedError(constant.ErrCodeUnauthorized, "invalid webhook signature")
	}

	var payload paddleWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "invalid webhook payload")
	}

	if !isPaddleSubscriptionEvent(payload.EventType) {
		return nil // not a subscription event — ignore
	}

	userIDStr := payload.Data.CustomData["user_id"]
	if userIDStr == "" {
		return nil
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil
	}

	d := payload.Data
	planID := paddlePriceToPlan(d.Items, s.billingCfg)
	status := normalisePaddleStatus(d.Status)

	// Revert to free on cancellation
	if status == model.SubStatusCancelled || status == model.SubStatusExpired {
		planID = PlanFree
	}

	var periodEnd *time.Time
	if d.CurrentBillingPeriod != nil && d.CurrentBillingPeriod.EndsAt != "" {
		if t, err := time.Parse(time.RFC3339, d.CurrentBillingPeriod.EndsAt); err == nil {
			periodEnd = &t
		}
	}

	var cancelledAt *time.Time
	if d.CanceledAt != nil && *d.CanceledAt != "" {
		if t, err := time.Parse(time.RFC3339, *d.CanceledAt); err == nil {
			cancelledAt = &t
		}
	}

	priceID := ""
	if len(d.Items) > 0 {
		priceID = d.Items[0].Price.ID
	}

	now := time.Now()
	sub := &model.Subscription{
		UserID:           userID,
		PlanID:           planID,
		Status:           status,
		PaddleSubID:      d.ID,
		PaddleCustomerID: d.CustomerID,
		PaddlePriceID:    priceID,
		CurrentPeriodEnd: periodEnd,
		CancelledAt:      cancelledAt,
		UpdatedAt:        now,
		CreatedAt:        now,
	}

	if svcErr := s.subRepo.Upsert(ctx, sub); svcErr != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to persist subscription")
	}
	return nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func toResponse(s *model.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:               s.ID.String(),
		UserID:           s.UserID.String(),
		PlanID:           s.PlanID,
		Status:           s.Status,
		CurrentPeriodEnd: s.CurrentPeriodEnd,
		CancelledAt:      s.CancelledAt,
		CreatedAt:        s.CreatedAt,
	}
}

func isPaddleSubscriptionEvent(name string) bool {
	switch name {
	case "subscription.created", "subscription.updated", "subscription.canceled":
		return true
	}
	return false
}

func paddlePriceToPlan(items []paddleSubItem, cfg *config.BillingConfig) string {
	if len(items) == 0 {
		return PlanFree
	}
	switch items[0].Price.ID {
	case cfg.PaddleProPriceID:
		return PlanPro
	case cfg.PaddleGrowthPriceID:
		return PlanGrowth
	}
	return PlanFree
}

func normalisePaddleStatus(status string) string {
	switch status {
	case "active":
		return model.SubStatusActive
	case "canceled":
		return model.SubStatusCancelled
	case "past_due":
		return model.SubStatusPastDue
	case "paused":
		return model.SubStatusPaused
	case "trialing":
		return model.SubStatusTrialing
	}
	return model.SubStatusActive
}
