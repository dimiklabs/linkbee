package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/repository"
	billingSvc "github.com/shafikshaon/shortlink/service/billing"
	"github.com/shafikshaon/shortlink/transport"
)

type BillingHandler struct {
	billingService billingSvc.BillingServiceI
	linkRepo       repository.LinkRepositoryI
	apiKeyRepo     repository.APIKeyRepositoryI
	webhookRepo    repository.WebhookRepositoryI
}

func NewBillingHandler(
	billingService billingSvc.BillingServiceI,
	linkRepo repository.LinkRepositoryI,
	apiKeyRepo repository.APIKeyRepositoryI,
	webhookRepo repository.WebhookRepositoryI,
) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
		linkRepo:       linkRepo,
		apiKeyRepo:     apiKeyRepo,
		webhookRepo:    webhookRepo,
	}
}

// GetSubscription returns the current user's subscription and plan info.
func (h *BillingHandler) GetSubscription(c *gin.Context) {
	ctx := c.Request.Context()
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	result, svcErr := h.billingService.GetSubscription(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "subscription retrieved", result)
}

// GetCheckoutURL returns a Lemon Squeezy checkout URL for the requested plan.
func (h *BillingHandler) GetCheckoutURL(c *gin.Context) {
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	planID := c.Param("plan")
	url, svcErr := h.billingService.GetCheckoutURL(planID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "checkout URL generated", gin.H{"checkout_url": url})
}

// GetUsage returns the current resource usage counts for the authenticated user.
func (h *BillingHandler) GetUsage(c *gin.Context) {
	ctx := c.Request.Context()
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	links, apiKeys, webhooks, err := h.fetchUsage(ctx, userID)
	if err != nil {
		transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, "failed to fetch usage")
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "usage retrieved", gin.H{
		"links":    links,
		"api_keys": apiKeys,
		"webhooks": webhooks,
	})
}

func (h *BillingHandler) fetchUsage(ctx context.Context, userID uuid.UUID) (links, apiKeys, webhooks int64, err error) {
	links, err = h.linkRepo.CountByUserID(ctx, userID)
	if err != nil {
		return
	}
	apiKeys, err = h.apiKeyRepo.CountByUserID(ctx, userID)
	if err != nil {
		return
	}
	webhooks, err = h.webhookRepo.CountByUserID(ctx, userID)
	return
}

// LemonSqueezyWebhook handles incoming Lemon Squeezy webhook events.
// The raw body is read before parsing so the HMAC signature can be verified.
func (h *BillingHandler) LemonSqueezyWebhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "failed to read request body")
		return
	}

	signature := c.GetHeader("X-Signature")
	if svcErr := h.billingService.HandleWebhook(c.Request.Context(), body, signature); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "webhook processed", nil)
}
