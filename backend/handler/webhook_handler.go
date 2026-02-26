package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	webhookSvc "github.com/shafikshaon/shortlink/service/webhook"
	"github.com/shafikshaon/shortlink/transport"
)

type WebhookHandler struct {
	webhookService webhookSvc.WebhookServiceI
}

func NewWebhookHandler(webhookService webhookSvc.WebhookServiceI) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService}
}

func (h *WebhookHandler) userID(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	id, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return uuid.Nil, false
	}
	return id, true
}

// ListWebhooks GET /api/v1/webhooks
func (h *WebhookHandler) ListWebhooks(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	webhooks, svcErr := h.webhookService.List(c.Request.Context(), userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Webhooks retrieved successfully", webhooks)
}

// CreateWebhook POST /api/v1/webhooks
func (h *WebhookHandler) CreateWebhook(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	var req request.CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	webhook, svcErr := h.webhookService.Create(c.Request.Context(), userID, req.URL, req.Events)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "Webhook created successfully", webhook)
}

// UpdateWebhook PUT /api/v1/webhooks/:id
func (h *WebhookHandler) UpdateWebhook(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid webhook ID")
		return
	}
	var req request.UpdateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	webhook, svcErr := h.webhookService.Update(c.Request.Context(), id, userID, req.URL, req.Events, isActive)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Webhook updated successfully", webhook)
}

// DeleteWebhook DELETE /api/v1/webhooks/:id
func (h *WebhookHandler) DeleteWebhook(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid webhook ID")
		return
	}
	if svcErr := h.webhookService.Delete(c.Request.Context(), id, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Webhook deleted successfully", nil)
}
