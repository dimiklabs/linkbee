package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	folderSvc "github.com/shafikshaon/shortlink/service/folder"
	"github.com/shafikshaon/shortlink/transport"
)

type FolderHandler struct {
	folderService folderSvc.FolderServiceI
}

func NewFolderHandler(folderService folderSvc.FolderServiceI) *FolderHandler {
	return &FolderHandler{folderService: folderService}
}

// ListFolders godoc
//
//	@Summary		List folders
//	@Description	Returns all folders for the authenticated user.
//	@Tags			folders
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/folders [get]
func (h *FolderHandler) ListFolders(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	folders, svcErr := h.folderService.ListFolders(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Folders retrieved successfully", folders)
}

// CreateFolder godoc
//
//	@Summary		Create a folder
//	@Description	Creates a new folder for organising links.
//	@Tags			folders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.CreateFolderRequest	true	"Folder details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/folders [post]
func (h *FolderHandler) CreateFolder(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req request.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	folder, svcErr := h.folderService.CreateFolder(ctx, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusCreated, "Folder created successfully", folder)
}

// UpdateFolder godoc
//
//	@Summary		Update a folder
//	@Description	Updates the name or icon of a folder.
//	@Tags			folders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string					true	"Folder UUID"
//	@Param			body	body		request.UpdateFolderRequest	true	"Updated fields"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/folders/{id} [put]
func (h *FolderHandler) UpdateFolder(c *gin.Context) {
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
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid folder ID")
		return
	}

	var req request.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	folder, svcErr := h.folderService.UpdateFolder(ctx, id, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Folder updated successfully", folder)
}

// DeleteFolder godoc
//
//	@Summary		Delete a folder
//	@Description	Deletes a folder (links inside are not deleted).
//	@Tags			folders
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Folder UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/folders/{id} [delete]
func (h *FolderHandler) DeleteFolder(c *gin.Context) {
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
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid folder ID")
		return
	}

	if svcErr := h.folderService.DeleteFolder(ctx, id, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Folder deleted successfully", nil)
}
