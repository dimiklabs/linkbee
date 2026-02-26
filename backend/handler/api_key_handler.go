package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	apiKeySvc "github.com/shafikshaon/shortlink/service/apikey"
	"github.com/shafikshaon/shortlink/transport"
)

type APIKeyHandler struct {
	svc apiKeySvc.APIKeyServiceI
}

func NewAPIKeyHandler(svc apiKeySvc.APIKeyServiceI) *APIKeyHandler {
	return &APIKeyHandler{svc: svc}
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

// ListAPIKeys GET /api-keys
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

// CreateAPIKey POST /api-keys
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

	resp, svcErr := h.svc.Create(c.Request.Context(), userID, req.Name, expiresAt)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": resp})
}

// RevokeAPIKey DELETE /api-keys/:id
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
	c.Status(http.StatusNoContent)
}
