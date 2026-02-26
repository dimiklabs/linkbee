package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/request"
	demoSvc "github.com/shafikshaon/shortlink/service/demo"
	"github.com/shafikshaon/shortlink/transport"
	"github.com/shafikshaon/shortlink/util"
)

type DemoHandler struct {
	demoService demoSvc.DemoServiceI
}

func NewDemoHandler(demoService demoSvc.DemoServiceI) *DemoHandler {
	return &DemoHandler{demoService: demoService}
}

// ShortenURL handles POST /api/v1/demo/shorten
func (h *DemoHandler) ShortenURL(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.DemoShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		code, msg := util.TranslateValidationError(err)
		transport.RespondWithError(c, http.StatusBadRequest, code, msg)
		return
	}

	ip := c.ClientIP()
	result, svcErr := h.demoService.ShortenURL(ctx, req.DestinationURL, ip)
	if svcErr != nil {
		statusCode := svcErr.StatusCode
		if svcErr.ErrorCode == constant.ErrCodeRateLimited {
			statusCode = http.StatusTooManyRequests
		}
		transport.RespondWithError(c, statusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusCreated, "URL shortened successfully", result)
}
