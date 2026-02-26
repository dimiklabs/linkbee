package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	linkSvc "github.com/shafikshaon/shortlink/service/link"
	"github.com/shafikshaon/shortlink/transport"
	"github.com/shafikshaon/shortlink/util"
)

type LinkHandler struct {
	linkService linkSvc.LinkServiceI
}

func NewLinkHandler(linkService linkSvc.LinkServiceI) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

// ListLinks godoc
// GET /api/v1/links
func (h *LinkHandler) ListLinks(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req request.ListLinksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	result, svcErr := h.linkService.ListLinks(ctx, userID, req.Page, req.Limit, req.Search)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Links retrieved successfully", result)
}

// CreateLink godoc
// POST /api/v1/links
func (h *LinkHandler) CreateLink(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req request.CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}

	result, svcErr := h.linkService.CreateLink(ctx, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusCreated, "Link created successfully", result)
}

// GetLink godoc
// GET /api/v1/links/:id
func (h *LinkHandler) GetLink(c *gin.Context) {
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

	result, svcErr := h.linkService.GetLink(ctx, id, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Link retrieved successfully", result)
}

// UpdateLink godoc
// PUT /api/v1/links/:id
func (h *LinkHandler) UpdateLink(c *gin.Context) {
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

	var req request.UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}

	result, svcErr := h.linkService.UpdateLink(ctx, id, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Link updated successfully", result)
}

// DeleteLink godoc
// DELETE /api/v1/links/:id
func (h *LinkHandler) DeleteLink(c *gin.Context) {
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

	if svcErr := h.linkService.DeleteLink(ctx, id, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Link deleted successfully", nil)
}

