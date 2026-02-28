package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	adminSvc "github.com/shafikshaon/linkbee/service/admin"
	"github.com/shafikshaon/linkbee/transport"
)

type AdminHandler struct {
	adminService adminSvc.AdminServiceI
	cfg          *config.AppConfig
}

func NewAdminHandler(adminService adminSvc.AdminServiceI, cfg *config.AppConfig) *AdminHandler {
	return &AdminHandler{adminService: adminService, cfg: cfg}
}

// GetStats returns platform-wide statistics.
func (h *AdminHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()
	stats, svcErr := h.adminService.GetStats(ctx)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "stats retrieved", stats)
}

// ListUsers returns a paginated, searchable list of all users.
func (h *AdminHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	search := strings.TrimSpace(c.Query("search"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, svcErr := h.adminService.ListUsers(ctx, search, page, limit)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "users retrieved", result)
}

// GetGrowthTimeSeries returns daily registration and link creation counts for the last 30 days.
func (h *AdminHandler) GetGrowthTimeSeries(c *gin.Context) {
	result, err := h.adminService.GetGrowthTimeSeries(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateUserStatus changes the status of a user (active / inactive / banned).
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid user id")
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "status is required")
		return
	}

	if svcErr := h.adminService.UpdateUserStatus(ctx, userID, body.Status); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "user status updated", nil)
}

// UpdateUserRole changes the role of a user (admin / user).
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	ctx := c.Request.Context()
	adminIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	adminID, _ := uuid.Parse(adminIDStr.(string))

	targetIDStr := c.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid user id")
		return
	}

	var body struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "role is required")
		return
	}

	if svcErr := h.adminService.UpdateUserRole(ctx, adminID, targetID, body.Role); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "user role updated", nil)
}

// ImpersonateUser generates a short-lived access token for the target user.
func (h *AdminHandler) ImpersonateUser(c *gin.Context) {
	ctx := c.Request.Context()
	adminIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	adminID, _ := uuid.Parse(adminIDStr.(string))

	targetIDStr := c.Param("id")
	targetID, err := uuid.Parse(targetIDStr)
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid user id")
		return
	}

	result, svcErr := h.adminService.ImpersonateUser(ctx, adminID, targetID, h.cfg)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "impersonation token generated", result)
}
