package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	analyticsSvc "github.com/shafikshaon/shortlink/service/analytics"
	"github.com/shafikshaon/shortlink/transport"
)

type AnalyticsHandler struct {
	analyticsService analyticsSvc.AnalyticsServiceI
}

func NewAnalyticsHandler(analyticsService analyticsSvc.AnalyticsServiceI) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

// GetLinkAnalytics handles GET /api/v1/links/:id/analytics
func (h *AnalyticsHandler) GetLinkAnalytics(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid link ID")
		return
	}

	var req request.AnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	// Default date range: last 30 days
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)

	if req.From != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, req.From); parseErr == nil {
			from = parsed
		}
	}
	if req.To != "" {
		if parsed, parseErr := time.Parse(time.RFC3339, req.To); parseErr == nil {
			to = parsed
		}
	}

	granularity := req.Granularity
	if granularity == "" {
		granularity = "day"
	}

	result, svcErr := h.analyticsService.GetLinkAnalytics(ctx, id, userID, from, to, granularity)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Analytics retrieved successfully", result)
}
