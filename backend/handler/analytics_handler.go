package handler

import (
	"io"
	"net/http"
	"strings"
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

// GetPeriodComparison godoc
//
//	@Summary		Period-over-period analytics comparison
//	@Description	Returns click metrics for the current period and the immediately preceding period of the same length, along with % change and trend direction.
//	@Tags			analytics
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path	string	true	"Link UUID"
//	@Param			from	query	string	false	"Period start (RFC3339, default 30 days ago)"
//	@Param			to		query	string	false	"Period end (RFC3339, default now)"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/analytics/comparison [get]
func (h *AnalyticsHandler) GetPeriodComparison(c *gin.Context) {
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

	result, svcErr := h.analyticsService.GetPeriodComparison(ctx, id, userID, from, to)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Period comparison retrieved successfully", result)
}

// GetMultiLinkComparison godoc
//
//	@Summary		Compare analytics for multiple links
//	@Description	Returns side-by-side analytics for 2–5 links owned by the authenticated user.
//	@Tags			analytics
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			ids		query	string	true	"Comma-separated link UUIDs (2–5)"
//	@Param			from	query	string	false	"Period start (RFC3339, default 30 days ago)"
//	@Param			to		query	string	false	"Period end (RFC3339, default now)"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/comparison [get]
func (h *AnalyticsHandler) GetMultiLinkComparison(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	idsParam := strings.TrimSpace(c.Query("ids"))
	if idsParam == "" {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "ids query parameter is required (comma-separated link UUIDs)")
		return
	}

	rawIDs := strings.Split(idsParam, ",")
	if len(rawIDs) < 2 || len(rawIDs) > 5 {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "provide between 2 and 5 link IDs")
		return
	}

	linkIDs := make([]uuid.UUID, 0, len(rawIDs))
	for _, raw := range rawIDs {
		id, parseErr := uuid.Parse(strings.TrimSpace(raw))
		if parseErr != nil {
			transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid link ID: "+raw)
			return
		}
		linkIDs = append(linkIDs, id)
	}

	var req request.AnalyticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

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

	result, svcErr := h.analyticsService.GetMultiLinkComparison(ctx, userID, linkIDs, from, to)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Comparison retrieved successfully", result)
}
