package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/request"
	reportSvc "github.com/shafikshaon/shortlink/service/reporting"
	"github.com/shafikshaon/shortlink/transport"
)

type ReportHandler struct {
	reportingService reportSvc.ReportingServiceI
}

func NewReportHandler(svc reportSvc.ReportingServiceI) *ReportHandler {
	return &ReportHandler{reportingService: svc}
}

func (h *ReportHandler) userID(c *gin.Context) (uuid.UUID, bool) {
	raw, _ := c.Get(middlewares.ContextKeyUserID)
	id, err := uuid.Parse(raw.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return uuid.Nil, false
	}
	return id, true
}

// ListReports godoc
//
//	@Summary		List scheduled reports
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Success		200	{object}	transport.StandardResponse
//	@Router			/api/v1/reports [get]
func (h *ReportHandler) ListReports(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	reports, svcErr := h.reportingService.ListReports(c.Request.Context(), userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Reports retrieved successfully", reports)
}

// CreateReport godoc
//
//	@Summary		Create a scheduled report
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			body	body		request.CreateReportRequest	true	"Report config"
//	@Success		201		{object}	transport.StandardResponse
//	@Router			/api/v1/reports [post]
func (h *ReportHandler) CreateReport(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	var req request.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	report, svcErr := h.reportingService.CreateReport(c.Request.Context(), userID, req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusCreated, "Report created successfully", report)
}

// GetReport godoc
//
//	@Summary		Get a single report
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Report UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Router			/api/v1/reports/{id} [get]
func (h *ReportHandler) GetReport(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid report ID")
		return
	}
	report, svcErr := h.reportingService.GetReport(c.Request.Context(), userID, id)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Report retrieved successfully", report)
}

// UpdateReport godoc
//
//	@Summary		Update a scheduled report
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path		string						true	"Report UUID"
//	@Param			body	body		request.UpdateReportRequest	true	"Fields to update"
//	@Success		200		{object}	transport.StandardResponse
//	@Router			/api/v1/reports/{id} [put]
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid report ID")
		return
	}
	var req request.UpdateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}
	report, svcErr := h.reportingService.UpdateReport(c.Request.Context(), userID, id, req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Report updated successfully", report)
}

// DeleteReport godoc
//
//	@Summary		Delete a scheduled report
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Report UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Router			/api/v1/reports/{id} [delete]
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid report ID")
		return
	}
	if svcErr := h.reportingService.DeleteReport(c.Request.Context(), userID, id); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Report deleted successfully", nil)
}

// SendReportNow godoc
//
//	@Summary		Trigger an immediate report delivery
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Report UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Router			/api/v1/reports/{id}/send [post]
func (h *ReportHandler) SendReportNow(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid report ID")
		return
	}
	if svcErr := h.reportingService.SendReportNow(c.Request.Context(), userID, id); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Report sent successfully", nil)
}

// GetReportDeliveries godoc
//
//	@Summary		List delivery history for a report
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id	path		string	true	"Report UUID"
//	@Success		200	{object}	transport.StandardResponse
//	@Router			/api/v1/reports/{id}/deliveries [get]
func (h *ReportHandler) GetReportDeliveries(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid report ID")
		return
	}
	deliveries, svcErr := h.reportingService.GetDeliveries(c.Request.Context(), userID, id)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}
	transport.RespondWithSuccess(c, http.StatusOK, "Deliveries retrieved successfully", deliveries)
}
