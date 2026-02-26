package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	adminSvc "github.com/shafikshaon/shortlink/service/admin"
	"github.com/shafikshaon/shortlink/transport"
)

type AdminHandler struct {
	adminService adminSvc.AdminServiceI
}

func NewAdminHandler(adminService adminSvc.AdminServiceI) *AdminHandler {
	return &AdminHandler{adminService: adminService}
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
