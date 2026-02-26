package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	bioSvc "github.com/shafikshaon/shortlink/service/bio"
	"github.com/shafikshaon/shortlink/transport"
)

type BioHandler struct {
	bioService bioSvc.BioServiceI
}

func NewBioHandler(bioService bioSvc.BioServiceI) *BioHandler {
	return &BioHandler{bioService: bioService}
}

func (h *BioHandler) userID(c *gin.Context) (uuid.UUID, bool) {
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	id, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return uuid.Nil, false
	}
	return id, true
}

// GetBioPage godoc
//
//	@Summary		Get current user's bio page
//	@Description	Returns the authenticated user's bio page (creates one if it doesn't exist).
//	@Tags			bio
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		500	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio [get]
func (h *BioHandler) GetBioPage(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	result, svcErr := h.bioService.GetOrCreate(c.Request.Context(), userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio page retrieved", result)
}

// UpdateBioPage godoc
//
//	@Summary		Update bio page settings
//	@Description	Updates the authenticated user's bio page (username, title, description, avatar, theme, published status).
//	@Tags			bio
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.UpdateBioPageRequest	true	"Bio page settings"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		409	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio [put]
func (h *BioHandler) UpdateBioPage(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	var req request.UpdateBioPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	result, svcErr := h.bioService.Update(c.Request.Context(), userID, bioSvc.UpdateBioRequest{
		Username:    req.Username,
		Title:       req.Title,
		Description: req.Description,
		AvatarURL:   req.AvatarURL,
		Theme:       req.Theme,
		IsPublished: req.IsPublished,
	})
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio page updated", result)
}

// GetPublicBioPage godoc
//
//	@Summary		Get public bio page by username
//	@Description	Returns a published bio page by username (no authentication required).
//	@Tags			bio
//	@Produce		json
//	@Param			username	path		string	true	"Bio page username"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/public/{username} [get]
func (h *BioHandler) GetPublicBioPage(c *gin.Context) {
	username := c.Param("username")
	result, svcErr := h.bioService.GetPublic(c.Request.Context(), username)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio page retrieved", result)
}

// ListBioLinks godoc
//
//	@Summary		List bio page links
//	@Description	Returns all links on the authenticated user's bio page.
//	@Tags			bio
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/links [get]
func (h *BioHandler) ListBioLinks(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	result, svcErr := h.bioService.ListLinks(c.Request.Context(), userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio links retrieved", result)
}

// CreateBioLink godoc
//
//	@Summary		Add a link to bio page
//	@Description	Adds a new link to the authenticated user's bio page.
//	@Tags			bio
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.CreateBioLinkRequest	true	"Link data"
//	@Success		201	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/links [post]
func (h *BioHandler) CreateBioLink(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	var req request.CreateBioLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	result, svcErr := h.bioService.CreateLink(c.Request.Context(), userID, req.Title, req.URL)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "Bio link created", result)
}

// UpdateBioLink godoc
//
//	@Summary		Update a bio page link
//	@Description	Updates a link on the authenticated user's bio page.
//	@Tags			bio
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Bio link UUID"
//	@Param			body	body		request.UpdateBioLinkRequest	true	"Link data"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/links/{id} [put]
func (h *BioHandler) UpdateBioLink(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid link ID")
		return
	}
	var req request.UpdateBioLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	result, svcErr := h.bioService.UpdateLink(c.Request.Context(), userID, linkID, req.Title, req.URL, isActive)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio link updated", result)
}

// DeleteBioLink godoc
//
//	@Summary		Delete a bio page link
//	@Description	Removes a link from the authenticated user's bio page.
//	@Tags			bio
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Bio link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/links/{id} [delete]
func (h *BioHandler) DeleteBioLink(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	linkID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid link ID")
		return
	}
	if svcErr := h.bioService.DeleteLink(c.Request.Context(), userID, linkID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio link deleted", nil)
}

// ReorderBioLinks godoc
//
//	@Summary		Reorder bio page links
//	@Description	Updates the display order of links on the authenticated user's bio page.
//	@Tags			bio
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.ReorderBioLinksRequest	true	"Ordered list of link IDs"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/links/reorder [patch]
func (h *BioHandler) ReorderBioLinks(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	var req request.ReorderBioLinksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid link ID: "+s)
			return
		}
		ids = append(ids, id)
	}
	if svcErr := h.bioService.ReorderLinks(c.Request.Context(), userID, ids); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Bio links reordered", nil)
}
