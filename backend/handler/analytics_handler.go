package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/request"
	analyticsSvc "github.com/shafikshaon/shortlink/service/analytics"
	"github.com/shafikshaon/shortlink/transport"
	"github.com/shafikshaon/shortlink/util"
)

type AnalyticsHandler struct {
	analyticsService analyticsSvc.AnalyticsServiceI
	linkRepo         repository.LinkRepositoryI
	clickRepo        repository.ClickEventRepositoryI
	appCfg           *config.AppConfig
}

func NewAnalyticsHandler(
	analyticsService analyticsSvc.AnalyticsServiceI,
	linkRepo repository.LinkRepositoryI,
	clickRepo repository.ClickEventRepositoryI,
	appCfg *config.AppConfig,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		linkRepo:         linkRepo,
		clickRepo:        clickRepo,
		appCfg:           appCfg,
	}
}

// GetLinkAnalytics godoc
//
//	@Summary		Get link analytics
//	@Description	Returns aggregated analytics for a link including click count, time-series data, top referrers, device breakdown, and country breakdown.
//	@Tags			analytics
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id			path	string	true	"Link UUID"
//	@Param			start		query	string	false	"Start date (RFC3339)"
//	@Param			end			query	string	false	"End date (RFC3339)"
//	@Param			granularity	query	string	false	"Time bucket: hour, day, week, month (default day)"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/analytics [get]
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

// StreamLiveCount handles GET /api/v1/links/:id/analytics/live (SSE)
// Auth is done inline via ?token= query param since EventSource cannot set custom headers.
func (h *AnalyticsHandler) StreamLiveCount(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, authErr := util.ValidateAccessToken(tokenStr, h.appCfg.JWTSecret, h.appCfg.JWTIssuer)
	if authErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	idStr := c.Param("id")
	linkID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid link id"})
		return
	}

	ctx := c.Request.Context()

	// Verify link ownership
	link, linkErr := h.linkRepo.GetByID(ctx, linkID)
	if linkErr != nil || link.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	first := true
	c.Stream(func(w io.Writer) bool {
		if first {
			first = false
			total, _ := h.clickRepo.GetClickCountByLinkID(ctx, linkID)
			unique, _ := h.clickRepo.GetUniqueClickCountByLinkID(ctx, linkID)
			c.SSEvent("count", gin.H{"total_clicks": total, "unique_clicks": unique})
			return true
		}
		select {
		case <-ticker.C:
			total, fetchErr := h.clickRepo.GetClickCountByLinkID(ctx, linkID)
			if fetchErr != nil {
				return false
			}
			unique, _ := h.clickRepo.GetUniqueClickCountByLinkID(ctx, linkID)
			c.SSEvent("count", gin.H{"total_clicks": total, "unique_clicks": unique})
			return true
		case <-ctx.Done():
			return false
		}
	})
}
