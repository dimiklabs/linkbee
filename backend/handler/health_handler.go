package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/service/health"
	"github.com/shafikshaon/linkbee/transport"
)

type HealthHandler struct {
	healthService health.HealthServiceI
}

func NewHealthHandler(healthService health.HealthServiceI) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// Check godoc
//
//	@Summary		Health check
//	@Description	Returns the health status of the API server and its dependencies.
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	transport.StandardResponse
//	@Failure		503	{object}	transport.StandardResponse
//	@Router			/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	ctx := c.Request.Context()

	logger.DebugCtx(ctx, "Health check requested",
		zap.String("client_ip", c.ClientIP()))

	result := h.healthService.Check(ctx)

	if result.Status == health.StatusUp {
		transport.RespondWithSuccess(c, http.StatusOK, "Service is healthy", result)
	} else {
		logger.WarnCtx(ctx, "Health check failed: service is degraded",
			zap.String("status", result.Status))
		transport.RespondWithSuccess(c, http.StatusServiceUnavailable, "Service is degraded", result)
	}
}
