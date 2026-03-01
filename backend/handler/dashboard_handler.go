package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	billingSvc "github.com/shafikshaon/linkbee/service/billing"
	dashboardSvc "github.com/shafikshaon/linkbee/service/dashboard"
	"github.com/shafikshaon/linkbee/transport"
)

type DashboardHandler struct {
	dashboardService dashboardSvc.DashboardServiceI
	planEnforcer     billingSvc.PlanEnforcerI
}

func NewDashboardHandler(dashboardService dashboardSvc.DashboardServiceI, planEnforcer billingSvc.PlanEnforcerI) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService, planEnforcer: planEnforcer}
}

// GetOverview godoc
//
//	@Summary		Dashboard overview
//	@Description	Returns aggregated stats (total links, total clicks, clicks last 7/30 days, top links, recent links).
//	@Tags			dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/dashboard/overview [get]
func (h *DashboardHandler) GetOverview(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	result, svcErr := h.dashboardService.GetOverview(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Dashboard overview retrieved successfully", result)
}

// GetGlobalAnalytics godoc
//
//	@Summary		Global analytics
//	@Description	Returns account-wide click analytics (time series, top countries, devices, browsers, referrers) for the given date range.
//	@Tags			dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			from	query	string	false	"Start date (RFC3339)"
//	@Param			to	query	string	false	"End date (RFC3339)"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/dashboard/analytics [get]
func (h *DashboardHandler) GetGlobalAnalytics(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	if svcErr := h.planEnforcer.CheckAnalytics(ctx, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -30)
	to := now

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t.UTC()
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t.UTC()
		}
	}

	result, svcErr := h.dashboardService.GetGlobalAnalytics(ctx, userID, from, to)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Global analytics retrieved successfully", result)
}

// GetGlobalAnalyticsComparison godoc
//
//	@Summary		Global analytics period comparison
//	@Description	Returns period-over-period comparison for account-wide click analytics.
//	@Tags			dashboard
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			from	query	string	false	"Start date (RFC3339)"
//	@Param			to	query	string	false	"End date (RFC3339)"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/dashboard/analytics/comparison [get]
func (h *DashboardHandler) GetGlobalAnalyticsComparison(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	if svcErr := h.planEnforcer.CheckAnalytics(ctx, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	now := time.Now().UTC()
	from := now.AddDate(0, 0, -30)
	to := now

	if fromStr := c.Query("from"); fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t.UTC()
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t.UTC()
		}
	}

	result, svcErr := h.dashboardService.GetGlobalAnalyticsComparison(ctx, userID, from, to)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Global analytics comparison retrieved successfully", result)
}
