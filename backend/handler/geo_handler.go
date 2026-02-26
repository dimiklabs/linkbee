package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/request"
	geoSvc "github.com/shafikshaon/shortlink/service/georouting"
	"github.com/shafikshaon/shortlink/transport"
)

type GeoHandler struct {
	geoRoutingService geoSvc.GeoRoutingServiceI
}

func NewGeoHandler(geoRoutingService geoSvc.GeoRoutingServiceI) *GeoHandler {
	return &GeoHandler{geoRoutingService: geoRoutingService}
}

func (h *GeoHandler) userID(c *gin.Context) uuid.UUID {
	id, _ := c.Get("userID")
	return id.(uuid.UUID)
}

func (h *GeoHandler) linkID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Invalid link ID")
		return uuid.Nil, false
	}
	return id, true
}

func (h *GeoHandler) ruleID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "INVALID_ID", "Invalid rule ID")
		return uuid.Nil, false
	}
	return id, true
}

// ListGeoRules godoc
//
//	@Summary		List geo-routing rules
//	@Description	Returns all geo-routing rules for a link.
//	@Tags			geo-routing
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Link UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/geo-rules [get]
func (h *GeoHandler) ListGeoRules(c *gin.Context) {
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	rules, svcErr := h.geoRoutingService.ListRules(c.Request.Context(), linkID, h.userID(c))
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rules})
}

// CreateGeoRule godoc
//
//	@Summary		Create a geo-routing rule
//	@Description	Adds a country-based redirect rule to a link.
//	@Tags			geo-routing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Link UUID"
//	@Param			body	body		request.CreateGeoRuleRequest	true	"Rule details"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/geo-rules [post]
func (h *GeoHandler) CreateGeoRule(c *gin.Context) {
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	var req request.CreateGeoRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	rule, svcErr := h.geoRoutingService.CreateRule(c.Request.Context(), linkID, h.userID(c), req.CountryCode, req.DestinationURL, req.Priority)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": rule})
}

// UpdateGeoRule godoc
//
//	@Summary		Update a geo-routing rule
//	@Description	Updates the country code, destination URL, or priority of a rule.
//	@Tags			geo-routing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Link UUID"
//	@Param			ruleId	path		string						true	"Rule UUID"
//	@Param			body	body		request.UpdateGeoRuleRequest	true	"Updated fields"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/geo-rules/{ruleId} [put]
func (h *GeoHandler) UpdateGeoRule(c *gin.Context) {
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	ruleID, ok := h.ruleID(c)
	if !ok {
		return
	}
	var req request.UpdateGeoRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	rule, svcErr := h.geoRoutingService.UpdateRule(c.Request.Context(), ruleID, linkID, h.userID(c), req.CountryCode, req.DestinationURL, req.Priority)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rule})
}

// DeleteGeoRule godoc
//
//	@Summary		Delete a geo-routing rule
//	@Description	Removes a geo-routing rule from a link.
//	@Tags			geo-routing
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path	string	true	"Link UUID"
//	@Param			ruleId	path	string	true	"Rule UUID"
//	@Success		204
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/geo-rules/{ruleId} [delete]
func (h *GeoHandler) DeleteGeoRule(c *gin.Context) {
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	ruleID, ok := h.ruleID(c)
	if !ok {
		return
	}
	if svcErr := h.geoRoutingService.DeleteRule(c.Request.Context(), ruleID, linkID, h.userID(c)); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.Status(http.StatusNoContent)
}

// ToggleGeoRouting godoc
//
//	@Summary		Toggle geo-routing
//	@Description	Enables or disables geo-routing for a link.
//	@Tags			geo-routing
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string							true	"Link UUID"
//	@Param			body	body		request.ToggleGeoRoutingRequest	true	"Toggle payload"
//	@Success		200		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		401		{object}	transport.ErrorResponse
//	@Failure		404		{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/geo-routing [patch]
func (h *GeoHandler) ToggleGeoRouting(c *gin.Context) {
	linkID, ok := h.linkID(c)
	if !ok {
		return
	}
	var req request.ToggleGeoRoutingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	link, svcErr := h.geoRoutingService.ToggleGeoRouting(c.Request.Context(), linkID, h.userID(c), req.Enabled)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": link})
}
