package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	"github.com/shafikshaon/linkbee/request"
	billingSvc "github.com/shafikshaon/linkbee/service/billing"
	bioSvc "github.com/shafikshaon/linkbee/service/bio"
	clickSvc "github.com/shafikshaon/linkbee/service/click"
	"github.com/shafikshaon/linkbee/transport"
)

type BioHandler struct {
	bioService      bioSvc.BioServiceI
	bioClickService clickSvc.BioClickServiceI
	planEnforcer    billingSvc.PlanEnforcerI
	uploadsDir      string
	baseURL         string
}

func NewBioHandler(bioService bioSvc.BioServiceI, bioClickService clickSvc.BioClickServiceI, planEnforcer billingSvc.PlanEnforcerI, uploadsDir, baseURL string) *BioHandler {
	return &BioHandler{
		bioService:      bioService,
		bioClickService: bioClickService,
		planEnforcer:    planEnforcer,
		uploadsDir:      uploadsDir,
		baseURL:         baseURL,
	}
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
	ctx := c.Request.Context()
	result, svcErr := h.bioService.GetOrCreate(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	// Click counts are a Pro feature; zero them out for Free-plan users.
	if h.planEnforcer.CheckAnalytics(ctx, userID) != nil {
		for i := range result.Links {
			result.Links[i].ClickCount = -1
		}
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
	ctx := c.Request.Context()
	result, svcErr := h.bioService.ListLinks(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	// Click counts are a Pro feature; zero them out for Free-plan users.
	if h.planEnforcer.CheckAnalytics(ctx, userID) != nil {
		for i := range result {
			result[i].ClickCount = -1
		}
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

// RecordBioLinkClick godoc
//
//	@Summary		Record a bio link click
//	@Description	Fire-and-forget endpoint called by the public bio page when a visitor clicks a link.
//	@Tags			bio
//	@Param			username	path	string	true	"Bio page username"
//	@Param			id			path	string	true	"Bio link UUID"
//	@Success		204	"No Content"
//	@Router			/api/v1/bio/public/{username}/links/{id}/click [post]
func (h *BioHandler) RecordBioLinkClick(c *gin.Context) {
	username := c.Param("username")
	linkIDStr := c.Param("id")

	linkID, err := uuid.Parse(linkIDStr)
	if err != nil {
		// Return 204 regardless — don't leak internal errors to the public.
		c.Status(http.StatusNoContent)
		return
	}

	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")

	// Fire-and-forget: enqueue asynchronously so the response is instant.
	go h.bioClickService.EnqueueBioClick(context.Background(), username, linkID, ip, userAgent, referrer)

	c.Status(http.StatusNoContent)
}

// UploadBioAvatar godoc
//
//	@Summary		Upload bio page avatar
//	@Description	Accepts a multipart image upload (max 2 MB, JPEG/PNG/GIF/WebP) and stores it as the bio page avatar.
//	@Tags			bio
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			avatar	formData	file	true	"Avatar image"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/bio/avatar [post]
func (h *BioHandler) UploadBioAvatar(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	// Validate content type.
	ct := header.Header.Get("Content-Type")
	allowed := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}
	ext, ok := allowed[ct]
	if !ok {
		// Fallback: derive from filename.
		fnExt := strings.ToLower(filepath.Ext(header.Filename))
		extMap := map[string]string{".jpg": ".jpg", ".jpeg": ".jpg", ".png": ".png", ".gif": ".gif", ".webp": ".webp"}
		ext, ok = extMap[fnExt]
		if !ok {
			transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "unsupported image type; use JPEG, PNG, GIF, or WebP")
			return
		}
	}

	// Validate size (max 2 MB).
	const maxSize = 2 << 20
	if header.Size > maxSize {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "avatar must be 2 MB or smaller")
		return
	}

	// Ensure upload directory exists.
	avatarDir := filepath.Join(h.uploadsDir, "avatars")
	if err := os.MkdirAll(avatarDir, 0o755); err != nil {
		transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, "could not create upload directory")
		return
	}

	// Save file as <userID><ext>, overwriting any previous avatar.
	filename := userID.String() + ext
	dst := filepath.Join(avatarDir, filename)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, "failed to save avatar")
		return
	}

	avatarURL := fmt.Sprintf("%s/uploads/avatars/%s", strings.TrimRight(h.baseURL, "/"), filename)

	// Persist the URL on the bio page (avatar-only update).
	result, svcErr := h.bioService.Update(c.Request.Context(), userID, bioSvc.UpdateBioRequest{
		AvatarURL:  avatarURL,
		AvatarOnly: true,
	})
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Avatar uploaded", result)
}
