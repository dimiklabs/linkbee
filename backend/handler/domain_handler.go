package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	domainSvc "github.com/shafikshaon/shortlink/service/domain"
	"github.com/shafikshaon/shortlink/transport"
)

type DomainHandler struct {
	domainService domainSvc.DomainServiceI
}

func NewDomainHandler(domainService domainSvc.DomainServiceI) *DomainHandler {
	return &DomainHandler{domainService: domainService}
}

// ListDomains godoc
//
//	@Summary		List custom domains
//	@Description	Returns all custom domains registered by the authenticated user.
//	@Tags			domains
//	@Produce		json
//	@Success		200	{object}	transport.SuccessResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Router			/api/v1/domains [get]
func (h *DomainHandler) ListDomains(c *gin.Context) {
	ctx := c.Request.Context()
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	domains, svcErr := h.domainService.ListDomains(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "domains retrieved", domains)
}

// AddDomain godoc
//
//	@Summary		Add a custom domain
//	@Description	Registers a new custom domain for the authenticated user. The domain starts in pending status and must be verified via DNS TXT record.
//	@Tags			domains
//	@Accept			json
//	@Produce		json
//	@Param			body	body		object{domain=string}	true	"Domain to add"
//	@Success		201	{object}	transport.SuccessResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		409	{object}	transport.ErrorResponse
//	@Router			/api/v1/domains [post]
func (h *DomainHandler) AddDomain(c *gin.Context) {
	ctx := c.Request.Context()
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, err.Error())
		return
	}

	result, svcErr := h.domainService.AddDomain(ctx, userID, req.Domain)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "domain added", result)
}

// VerifyDomain godoc
//
//	@Summary		Verify a custom domain
//	@Description	Triggers a DNS TXT record check for the domain. The record _shortlink-verify.<domain> must contain the verify_token.
//	@Tags			domains
//	@Produce		json
//	@Param			id	path	string	true	"Domain ID"
//	@Success		200	{object}	transport.SuccessResponse
//	@Failure		400	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/domains/{id}/verify [post]
func (h *DomainHandler) VerifyDomain(c *gin.Context) {
	ctx := c.Request.Context()
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	domainID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid domain ID")
		return
	}

	result, svcErr := h.domainService.VerifyDomain(ctx, domainID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "domain verified", result)
}

// DeleteDomain godoc
//
//	@Summary		Delete a custom domain
//	@Description	Removes a custom domain registration for the authenticated user.
//	@Tags			domains
//	@Produce		json
//	@Param			id	path	string	true	"Domain ID"
//	@Success		200	{object}	transport.SuccessResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/domains/{id} [delete]
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	ctx := c.Request.Context()
	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "unauthorized")
		return
	}

	domainID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "invalid domain ID")
		return
	}

	if svcErr := h.domainService.DeleteDomain(ctx, domainID, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "domain deleted", nil)
}
