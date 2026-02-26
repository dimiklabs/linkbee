package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	splitSvc "github.com/shafikshaon/shortlink/service/split"
	"github.com/shafikshaon/shortlink/transport"
	"github.com/shafikshaon/shortlink/util"
)

type SplitHandler struct {
	splitService splitSvc.SplitServiceI
}

func NewSplitHandler(splitService splitSvc.SplitServiceI) *SplitHandler {
	return &SplitHandler{splitService: splitService}
}

func (h *SplitHandler) userID(c *gin.Context) (uuid.UUID, bool) {
	raw, _ := c.Get(middlewares.ContextKeyUserID)
	id, err := uuid.Parse(raw.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return uuid.Nil, false
	}
	return id, true
}

func (h *SplitHandler) linkID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid link ID")
		return uuid.Nil, false
	}
	return id, true
}

// ListVariants godoc
//
//	@Summary		List split test variants
//	@Description	Returns all A/B test variants for a link.
//	@Tags			split-testing
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/variants [get]
func (h *SplitHandler) ListVariants(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	result, svcErr := h.splitService.ListVariants(ctx, linkID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Variants retrieved successfully", result)
}

// CreateVariant godoc
//
//	@Summary		Create a variant
//	@Description	Adds a weighted destination URL variant for A/B testing.
//	@Tags			split-testing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string					true	"Link UUID"
//	@Param			body	body		request.CreateVariantRequest	true	"Variant details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/variants [post]
func (h *SplitHandler) CreateVariant(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	var req request.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}
	result, svcErr := h.splitService.CreateVariant(ctx, linkID, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "Variant created successfully", result)
}

// UpdateVariant godoc
//
//	@Summary		Update a variant
//	@Description	Updates a variant's destination URL or weight.
//	@Tags			split-testing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id			path		string					true	"Link UUID"
//	@Param			variantId	path		string					true	"Variant UUID"
//	@Param			body		body		request.UpdateVariantRequest	true	"Updated fields"
//	@Success		200			{object}	transport.StandardResponse
//	@Failure		400			{object}	transport.ErrorResponse
//	@Failure		401			{object}	transport.ErrorResponse
//	@Failure		404			{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/variants/{variantId} [put]
func (h *SplitHandler) UpdateVariant(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	variantID, err := uuid.Parse(c.Param("variantId"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid variant ID")
		return
	}
	var req request.UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}
	result, svcErr := h.splitService.UpdateVariant(ctx, variantID, linkID, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Variant updated successfully", result)
}

// DeleteVariant godoc
//
//	@Summary		Delete a variant
//	@Description	Removes an A/B test variant.
//	@Tags			split-testing
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id			path		string	true	"Link UUID"
//	@Param			variantId	path		string	true	"Variant UUID"
//	@Success		200			{object}	transport.StandardResponse
//	@Failure		401			{object}	transport.ErrorResponse
//	@Failure		404			{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/variants/{variantId} [delete]
func (h *SplitHandler) DeleteVariant(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	variantID, err := uuid.Parse(c.Param("variantId"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid variant ID")
		return
	}
	if svcErr := h.splitService.DeleteVariant(ctx, variantID, linkID, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Variant deleted successfully", nil)
}

// ToggleSplitTest godoc
//
//	@Summary		Toggle A/B split test
//	@Description	Enables or disables A/B split testing for a link.
//	@Tags			split-testing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Link UUID"
//	@Param			body	body		request.ToggleSplitTestRequest	true	"Toggle payload"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/split-test [patch]
func (h *SplitHandler) ToggleSplitTest(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	var req request.ToggleSplitTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, err.Error())
		return
	}
	result, svcErr := h.splitService.ToggleSplitTest(ctx, linkID, userID, req.Enabled)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Split test toggled successfully", result)
}
