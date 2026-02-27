package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	billingSvc "github.com/shafikshaon/shortlink/service/billing"
	webhookSvc "github.com/shafikshaon/shortlink/service/webhook"
	"github.com/shafikshaon/shortlink/transport"
)

type WebhookHandler struct {
	webhookService webhookSvc.WebhookServiceI
	planEnforcer   billingSvc.PlanEnforcerI
}

func NewWebhookHandler(webhookService webhookSvc.WebhookServiceI, planEnforcer billingSvc.PlanEnforcerI) *WebhookHandler {
	return &WebhookHandler{webhookService: webhookService, planEnforcer: planEnforcer}
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

// ListWebhooks godoc
//
//	@Summary		List webhooks
//	@Description	Returns all webhooks for the authenticated user.
//	@Tags			webhooks
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/webhooks [get]
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

// CreateWebhook godoc
//
//	@Summary		Create a webhook
//	@Description	Registers a new webhook endpoint. The signing secret is generated server-side.
//	@Tags			webhooks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.CreateWebhookRequest	true	"Webhook details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/webhooks [post]
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

	if svcErr := h.planEnforcer.CheckWebhookLimit(c.Request.Context(), userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	webhook, svcErr := h.webhookService.Create(c.Request.Context(), userID, req.URL, req.Events)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "Webhook created successfully", webhook)
}

// UpdateWebhook godoc
//
//	@Summary		Update a webhook
//	@Description	Updates the URL, subscribed events, or active status of a webhook.
//	@Tags			webhooks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Webhook UUID"
//	@Param			body	body		request.UpdateWebhookRequest	true	"Fields to update"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/webhooks/{id} [put]
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

// DeleteWebhook godoc
//
//	@Summary		Delete a webhook
//	@Description	Permanently deletes a webhook.
//	@Tags			webhooks
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Webhook UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/webhooks/{id} [delete]
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
