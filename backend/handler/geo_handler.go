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

// ListGeoRules GET /links/:id/geo-rules
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

// CreateGeoRule POST /links/:id/geo-rules
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

// UpdateGeoRule PUT /links/:id/geo-rules/:ruleId
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

// DeleteGeoRule DELETE /links/:id/geo-rules/:ruleId
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

// ToggleGeoRouting PATCH /links/:id/geo-routing
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
