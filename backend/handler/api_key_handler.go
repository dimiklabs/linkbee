package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	"github.com/shafikshaon/linkbee/model"
	apiKeySvc "github.com/shafikshaon/linkbee/service/apikey"
	auditSvc "github.com/shafikshaon/linkbee/service/audit"
	billingSvc "github.com/shafikshaon/linkbee/service/billing"
	"github.com/shafikshaon/linkbee/transport"
)

type APIKeyHandler struct {
	svc          apiKeySvc.APIKeyServiceI
	planEnforcer billingSvc.PlanEnforcerI
	auditService auditSvc.AuditServiceI
}

func NewAPIKeyHandler(svc apiKeySvc.APIKeyServiceI, planEnforcer billingSvc.PlanEnforcerI, auditService auditSvc.AuditServiceI) *APIKeyHandler {
	return &APIKeyHandler{svc: svc, planEnforcer: planEnforcer, auditService: auditService}
}

func (h *APIKeyHandler) userID(c *gin.Context) (uuid.UUID, bool) {
	v, _ := c.Get(middlewares.ContextKeyUserID)
	id, err := uuid.Parse(v.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return uuid.Nil, false
	}
	return id, true
}

// ListAPIKeys godoc
//
//	@Summary		List API keys
//	@Description	Returns all API keys for the authenticated user.
//	@Tags			api-keys
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/api-keys [get]
func (h *APIKeyHandler) ListAPIKeys(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	keys, svcErr := h.svc.List(c.Request.Context(), userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": keys})
}

// CreateAPIKey godoc
//
//	@Summary		Create an API key
//	@Description	Creates a new API key. The full key is returned once and cannot be retrieved again.
//	@Tags			api-keys
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body	object{name=string,expires_at=string}	true	"Key details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/api-keys [post]
func (h *APIKeyHandler) CreateAPIKey(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}

	var req struct {
		Name      string  `json:"name" binding:"required,min=1,max=100"`
		ExpiresAt *string `json:"expires_at"` // optional ISO-8601 date
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			transport.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", "expires_at must be ISO-8601 (RFC3339)")
			return
		}
		expiresAt = &t
	}

	if svcErr := h.planEnforcer.CheckAPIKeyLimit(c.Request.Context(), userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	resp, svcErr := h.svc.Create(c.Request.Context(), userID, req.Name, expiresAt)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionAPIKeyCreated,
		ResourceType: model.AuditResourceAPIKey, ResourceName: req.Name,
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

// RevokeAPIKey godoc
//
//	@Summary		Revoke an API key
//	@Description	Permanently revokes an API key.
//	@Tags			api-keys
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	string	true	"API key UUID"
//	@Success		204
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/api-keys/{id} [delete]
func (h *APIKeyHandler) RevokeAPIKey(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Invalid API key ID")
		return
	}
	if svcErr := h.svc.Revoke(c.Request.Context(), id, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionAPIKeyRevoked,
		ResourceType: model.AuditResourceAPIKey, ResourceID: id.String(),
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})
	c.Status(http.StatusNoContent)
}
