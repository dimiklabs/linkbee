package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/transport"
)

type AuditLogHandler struct {
	repo repository.AuditLogRepositoryI
}

func NewAuditLogHandler(repo repository.AuditLogRepositoryI) *AuditLogHandler {
	return &AuditLogHandler{repo: repo}
}

// ListAuditLogs godoc
//
//	@Summary		List audit logs
//	@Description	Returns a paginated list of audit log entries for the authenticated user.
//	@Tags			audit
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page			query	int		false	"Page (default 1)"
//	@Param			limit			query	int		false	"Items per page (default 20, max 100)"
//	@Param			action			query	string	false	"Filter by action"
//	@Param			resource_type	query	string	false	"Filter by resource type"
//	@Param			from			query	string	false	"Start date (RFC3339)"
//	@Param			to				query	string	false	"End date (RFC3339)"
//	@Success		200	{object}	transport.SuccessResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/audit-logs [get]
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	ctx := c.Request.Context()

	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		// fall back to string parse (some handlers store as string)
		if s, ok2 := rawID.(string); ok2 {
			parsed, err := uuid.Parse(s)
			if err != nil {
				transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
				return
			}
			userID = parsed
		} else {
			transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
			return
		}
	}

	page := 1
	limit := 20
	if v := c.Query("page"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	action := c.Query("action")
	resourceType := c.Query("resource_type")

	var from, to *time.Time
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from2 := t
			to = &from2
		}
	}

	logs, total, err := h.repo.ListByUserID(ctx, userID, page, limit, action, resourceType, from, to)
	if err != nil {
		transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, "failed to fetch audit logs")
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "audit logs retrieved", gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscan(s, &n)
	return n, err
}
