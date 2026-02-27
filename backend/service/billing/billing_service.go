package billing

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
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

// ─── Lemon Squeezy webhook payload types ─────────────────────────────────────

type lsWebhookMeta struct {
	EventName  string            `json:"event_name"`
	CustomData map[string]string `json:"custom_data"`
}

type lsSubAttributes struct {
	Status     string     `json:"status"`
	VariantID  int64      `json:"variant_id"`
	OrderID    int64      `json:"order_id"`
	CustomerID int64      `json:"customer_id"`
	Cancelled  bool       `json:"cancelled"`
	EndsAt     *time.Time `json:"ends_at"`
	RenewsAt   *time.Time `json:"renews_at"`
}

type lsSubData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes lsSubAttributes `json:"attributes"`
}

type lsWebhookPayload struct {
	Meta lsWebhookMeta `json:"meta"`
	Data lsSubData     `json:"data"`
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

	// No record → synthesise a free subscription
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

// GetCheckoutURL builds a Lemon Squeezy checkout URL for the requested plan,
// embedding the user_id as custom checkout data.
func (s *billingService) GetCheckoutURL(planID string, userID uuid.UUID) (string, *dto.ServiceError) {
	var base string
	switch planID {
	case PlanPro:
		base = s.billingCfg.ProCheckoutURL
	case PlanBusiness:
		base = s.billingCfg.BusinessCheckoutURL
	default:
		return "", dto.NewBadRequestError(constant.ErrCodeBadRequest, "invalid plan: must be pro or business")
	}

	if base == "" {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, "checkout URL not configured for this plan")
	}

	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	url := fmt.Sprintf("%s%scheckout[custom][user_id]=%s", base, sep, userID.String())
	return url, nil
}

// VerifyWebhookSignature checks the HMAC-SHA256 signature from Lemon Squeezy.
func (s *billingService) VerifyWebhookSignature(body []byte, signature string) bool {
	if s.billingCfg.LemonSqueezyWebhookSecret == "" {
		return true // skip verification when secret not configured (dev mode)
	}
	mac := hmac.New(sha256.New, []byte(s.billingCfg.LemonSqueezyWebhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// HandleWebhook processes a verified Lemon Squeezy webhook event.
func (s *billingService) HandleWebhook(ctx context.Context, body []byte, signature string) *dto.ServiceError {
	if !s.VerifyWebhookSignature(body, signature) {
		return dto.NewUnauthorizedError(constant.ErrCodeUnauthorized, "invalid webhook signature")
	}

	var payload lsWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "invalid webhook payload")
	}

	eventName := payload.Meta.EventName
	if !isSubscriptionEvent(eventName) {
		return nil // not a subscription event — ignore
	}

	// Extract user_id from custom data
	userIDStr := payload.Meta.CustomData["user_id"]
	if userIDStr == "" {
		return nil // no user_id — cannot process
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil
	}

	attr := payload.Data.Attributes
	planID := variantToPlan(payload.Data.Attributes.VariantID, s.billingCfg)
	status := normaliseLSStatus(attr.Status, attr.Cancelled)

	var periodEnd *time.Time
	if attr.RenewsAt != nil {
		periodEnd = attr.RenewsAt
	} else if attr.EndsAt != nil {
		periodEnd = attr.EndsAt
	}

	var cancelledAt *time.Time
	if attr.Cancelled {
		now := time.Now()
		cancelledAt = &now
	}

	// If expired/cancelled with no paid plan → revert to free
	if (status == model.SubStatusExpired || status == model.SubStatusCancelled) && planID != PlanFree {
		planID = PlanFree
	}

	now := time.Now()
	sub := &model.Subscription{
		UserID:                 userID,
		PlanID:                 planID,
		Status:                 status,
		LemonSqueezySubID:      payload.Data.ID,
		LemonSqueezyOrderID:    fmt.Sprintf("%d", attr.OrderID),
		LemonSqueezyCustomerID: fmt.Sprintf("%d", attr.CustomerID),
		LemonSqueezyVariantID:  fmt.Sprintf("%d", attr.VariantID),
		CurrentPeriodEnd:       periodEnd,
		CancelledAt:            cancelledAt,
		UpdatedAt:              now,
		CreatedAt:              now,
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

func isSubscriptionEvent(name string) bool {
	switch name {
	case "subscription_created", "subscription_updated",
		"subscription_cancelled", "subscription_expired", "subscription_resumed":
		return true
	}
	return false
}

// variantToPlan maps a Lemon Squeezy variant ID to an internal plan ID.
func variantToPlan(variantID int64, cfg *config.BillingConfig) string {
	variantStr := fmt.Sprintf("%d", variantID)
	switch variantStr {
	case cfg.ProVariantID:
		return PlanPro
	case cfg.BusinessVariantID:
		return PlanBusiness
	}
	return PlanFree
}

// normaliseLSStatus maps a LS status string to an internal status.
func normaliseLSStatus(lsStatus string, cancelled bool) string {
	if cancelled {
		return model.SubStatusCancelled
	}
	switch lsStatus {
	case "active":
		return model.SubStatusActive
	case "past_due":
		return model.SubStatusPastDue
	case "paused":
		return model.SubStatusPaused
	case "expired":
		return model.SubStatusExpired
	case "cancelled":
		return model.SubStatusCancelled
	case "on_trial":
		return model.SubStatusTrialing
	}
	return model.SubStatusActive
}
