package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	linkSvc "github.com/shafikshaon/shortlink/service/link"
	previewSvc "github.com/shafikshaon/shortlink/service/preview"
	"github.com/shafikshaon/shortlink/transport"
)

type PreviewHandler struct {
	previewService previewSvc.PreviewServiceI
	linkService    linkSvc.LinkServiceI
}

func NewPreviewHandler(previewService previewSvc.PreviewServiceI, linkService linkSvc.LinkServiceI) *PreviewHandler {
	return &PreviewHandler{
		previewService: previewService,
		linkService:    linkService,
	}
}

// GetLinkPreview godoc
//
//	@Summary		Get OG/link preview metadata
//	@Description	Fetches Open Graph metadata (title, description, image, site name, favicon) for a link's destination URL. Results are cached for 24 hours.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/preview [get]
func (h *PreviewHandler) GetLinkPreview(c *gin.Context) {
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid link ID")
		return
	}

	// Verify ownership and get destination URL
	link, svcErr := h.linkService.GetLink(c.Request.Context(), linkID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	preview, svcErr := h.previewService.GetPreview(c.Request.Context(), linkID, link.DestinationURL)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Link preview retrieved", preview)
}
