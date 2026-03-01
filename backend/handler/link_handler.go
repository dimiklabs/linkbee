package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/request"
	auditSvc "github.com/shafikshaon/linkbee/service/audit"
	billingSvc "github.com/shafikshaon/linkbee/service/billing"
	linkSvc "github.com/shafikshaon/linkbee/service/link"
	webhookSvc "github.com/shafikshaon/linkbee/service/webhook"
	"github.com/shafikshaon/linkbee/transport"
	"github.com/shafikshaon/linkbee/util"
)

type LinkHandler struct {
	linkService    linkSvc.LinkServiceI
	webhookService webhookSvc.WebhookServiceI
	planEnforcer   billingSvc.PlanEnforcerI
	auditService   auditSvc.AuditServiceI
}

func NewLinkHandler(linkService linkSvc.LinkServiceI, webhookService webhookSvc.WebhookServiceI, planEnforcer billingSvc.PlanEnforcerI, auditService auditSvc.AuditServiceI) *LinkHandler {
	return &LinkHandler{linkService: linkService, webhookService: webhookService, planEnforcer: planEnforcer, auditService: auditService}
}

// ListLinks godoc
//
//	@Summary		List links
//	@Description	Returns a paginated list of links for the authenticated user.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			page			query		int		false	"Page number (default 1)"
//	@Param			limit			query		int		false	"Items per page (default 20)"
//	@Param			search			query		string	false	"Search by slug, title or URL"
//	@Param			folder_id		query		string	false	"Filter by folder UUID"
//	@Param			starred			query		bool	false	"Filter starred links"
//	@Param			health_status	query		string	false	"Filter by health status"
//	@Success		200				{object}	transport.StandardResponse
//	@Failure		401				{object}	transport.ErrorResponse
//	@Router			/api/v1/links [get]
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

	var folderFilter *uuid.UUID
	if req.FolderID != "" {
		if parsed, err := uuid.Parse(req.FolderID); err == nil {
			folderFilter = &parsed
		}
	}

	result, svcErr := h.linkService.ListLinks(ctx, userID, req.Page, req.Limit, req.Search, folderFilter, req.Starred, req.HealthStatus, req.Tags, req.ExpiringSoon)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Links retrieved successfully", result)
}

// GetUserTags godoc
//
//	@Summary		List all tags
//	@Description	Returns all unique tags used across the authenticated user's links.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/tags [get]
func (h *LinkHandler) GetUserTags(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	tags, svcErr := h.linkService.GetUserTags(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Tags retrieved successfully", tags)
}

// CreateLink godoc
//
//	@Summary		Create a link
//	@Description	Creates a new shortened link.
//	@Tags			links
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.CreateLinkRequest	true	"Link details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		409		{object}	transport.ErrorResponse
//	@Router			/api/v1/links [post]
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

	if svcErr := h.planEnforcer.CheckLinkLimit(ctx, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	if req.Slug != "" {
		if svcErr := h.planEnforcer.CheckCustomSlug(ctx, userID); svcErr != nil {
			transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
			return
		}
	}

	result, svcErr := h.linkService.CreateLink(ctx, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	h.webhookService.Trigger(userID, webhookSvc.EventLinkCreated, result)
	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionLinkCreated,
		ResourceType: model.AuditResourceLink, ResourceID: result.ID.String(), ResourceName: result.Slug,
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusCreated, "Link created successfully", result)
}

// GetLink godoc
//
//	@Summary		Get a link
//	@Description	Returns a single link by ID.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id} [get]
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
//
//	@Summary		Update a link
//	@Description	Updates fields of an existing link.
//	@Tags			links
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string					true	"Link UUID"
//	@Param			body	body		request.UpdateLinkRequest	true	"Fields to update"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id} [put]
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

	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionLinkUpdated,
		ResourceType: model.AuditResourceLink, ResourceID: result.ID.String(), ResourceName: result.Slug,
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusOK, "Link updated successfully", result)
}

// DeleteLink godoc
//
//	@Summary		Delete a link
//	@Description	Permanently deletes a link.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id} [delete]
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

	h.webhookService.Trigger(userID, webhookSvc.EventLinkDeleted, map[string]string{"id": id.String()})
	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionLinkDeleted,
		ResourceType: model.AuditResourceLink, ResourceID: id.String(),
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusOK, "Link deleted successfully", nil)
}

// CheckLinkHealth godoc
//
//	@Summary		Check link health
//	@Description	Probes the destination URL and updates its health status.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/health-check [post]
func (h *LinkHandler) CheckLinkHealth(c *gin.Context) {
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

	result, svcErr := h.linkService.CheckLinkHealth(ctx, id, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Health check complete", result)
}

// ImportLinks godoc
//
//	@Summary		Bulk import links from CSV
//	@Description	Accepts a multipart CSV file with a 'destination_url' column (max 500 rows).
//	@Tags			links
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			file	formData	file	true	"CSV file"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/import [post]
func (h *LinkHandler) ImportLinks(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	if c.Request.ContentLength > 5*1024*1024 {
		transport.RespondWithError(c, http.StatusRequestEntityTooLarge, constant.ErrCodeBadRequest, "File too large. Maximum file size is 5MB")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "A CSV file is required (field name: file)")
		return
	}
	defer file.Close()

	result, svcErr := h.linkService.ImportLinks(ctx, userID, file)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionLinksImported,
		ResourceType: model.AuditResourceLink, ResourceName: fmt.Sprintf("%d created", result.Created),
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusOK, "Import complete", result)
}

// ToggleStar godoc
//
//	@Summary		Toggle star on a link
//	@Description	Toggles the starred flag on the given link.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/star [patch]
func (h *LinkHandler) ToggleStar(c *gin.Context) {
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

	result, svcErr := h.linkService.ToggleStar(ctx, id, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Link star toggled successfully", result)
}

// BulkAction godoc
//
//	@Summary		Bulk link action
//	@Description	Performs a bulk action (delete, activate, deactivate, move_folder, add_tags, remove_tags) on multiple links.
//	@Tags			links
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.BulkLinkActionRequest	true	"Bulk action payload"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/bulk [post]
func (h *LinkHandler) BulkAction(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req request.BulkLinkActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}

	result, svcErr := h.linkService.BulkAction(ctx, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, fmt.Sprintf("Bulk %s completed: %d link(s) affected", req.Action, result.Affected), result)
}

// CloneLink godoc
//
//	@Summary		Clone a link
//	@Description	Creates a copy of an existing link with a new slug, resetting click count and flags.
//	@Tags			links
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Source link UUID"
//	@Param			body	body		request.CloneLinkRequest	false	"Optional overrides for the clone"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Failure		409		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/clone [post]
func (h *LinkHandler) CloneLink(c *gin.Context) {
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

	var req request.CloneLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}

	if svcErr := h.planEnforcer.CheckLinkLimit(ctx, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	result, svcErr := h.linkService.CloneLink(ctx, id, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	h.webhookService.Trigger(userID, webhookSvc.EventLinkCreated, result)
	h.auditService.LogAsync(auditSvc.LogEntry{
		UserID: userID, Action: model.AuditActionLinkCreated,
		ResourceType: model.AuditResourceLink, ResourceID: result.ID.String(), ResourceName: result.Slug,
		IPAddress: c.ClientIP(), UserAgent: c.GetHeader("User-Agent"),
	})

	transport.RespondWithSuccess(c, http.StatusCreated, "Link cloned successfully", result)
}

// CheckDuplicate godoc
//
//	@Summary		Check if a destination URL already exists
//	@Description	Returns the existing shortened link if the given destination URL is already in the user's account.
//	@Tags			links
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			url	query		string	true	"Destination URL to check"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/duplicate [get]
func (h *LinkHandler) CheckDuplicate(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	destURL := c.Query("url")
	if destURL == "" {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "url query parameter is required")
		return
	}

	result, svcErr := h.linkService.CheckDuplicate(ctx, userID, destURL)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Duplicate link found", result)
}

