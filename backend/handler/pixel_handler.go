package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	"github.com/shafikshaon/linkbee/request"
	pixelSvc "github.com/shafikshaon/linkbee/service/pixel"
	"github.com/shafikshaon/linkbee/transport"
)

type PixelHandler struct {
	pixelService pixelSvc.PixelServiceI
}

func NewPixelHandler(pixelService pixelSvc.PixelServiceI) *PixelHandler {
	return &PixelHandler{pixelService: pixelService}
}

func (h *PixelHandler) userID(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	id, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return uuid.Nil, false
	}
	return id, true
}

func (h *PixelHandler) linkID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid link ID")
		return uuid.Nil, false
	}
	return id, true
}

// ListPixels godoc
//
//	@Summary		List retargeting pixels
//	@Description	Returns all retargeting pixels configured for a link.
//	@Tags			pixels
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/pixels [get]
func (h *PixelHandler) ListPixels(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	pixels, svcErr := h.pixelService.List(c.Request.Context(), linkID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Pixels retrieved successfully", pixels)
}

// CreatePixel godoc
//
//	@Summary		Add a retargeting pixel
//	@Description	Adds a tracking pixel (Facebook, Google Ads, TikTok, LinkedIn, or custom) to a link.
//	@Tags			pixels
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string					true	"Link UUID"
//	@Param			body	body		request.CreatePixelRequest	true	"Pixel details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/pixels [post]
func (h *PixelHandler) CreatePixel(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	var req request.CreatePixelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	pixel, svcErr := h.pixelService.Create(c.Request.Context(), linkID, userID, req.PixelType, req.PixelID, req.CustomScript)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "Pixel created successfully", pixel)
}

// DeletePixel godoc
//
//	@Summary		Delete a retargeting pixel
//	@Description	Removes a tracking pixel from a link.
//	@Tags			pixels
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path	string	true	"Link UUID"
//	@Param			pixelId	path	string	true	"Pixel UUID"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/pixels/{pixelId} [delete]
func (h *PixelHandler) DeletePixel(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	pixelID, err := uuid.Parse(c.Param("pixelId"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid pixel ID")
		return
	}
	if svcErr := h.pixelService.Delete(c.Request.Context(), pixelID, linkID, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Pixel deleted successfully", nil)
}

// TogglePixelTracking godoc
//
//	@Summary		Toggle pixel tracking
//	@Description	Enables or disables pixel tracking for a link. When enabled, clicks serve an intermediate HTML page that fires all configured pixels before redirecting.
//	@Tags			pixels
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string							true	"Link UUID"
//	@Param			body	body		request.TogglePixelTrackingRequest	true	"Toggle payload"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/pixel-tracking [patch]
func (h *PixelHandler) TogglePixelTracking(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	var req request.TogglePixelTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	if svcErr := h.pixelService.TogglePixelTracking(c.Request.Context(), linkID, userID, req.Enabled); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Pixel tracking updated", nil)
}
