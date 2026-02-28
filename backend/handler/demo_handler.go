package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/request"
	demoSvc "github.com/shafikshaon/linkbee/service/demo"
	"github.com/shafikshaon/linkbee/transport"
	"github.com/shafikshaon/linkbee/util"
)

type DemoHandler struct {
	demoService demoSvc.DemoServiceI
}

func NewDemoHandler(demoService demoSvc.DemoServiceI) *DemoHandler {
	return &DemoHandler{demoService: demoService}
}

// ShortenURL godoc
//
//	@Summary		Demo URL shortener
//	@Description	Shorten a URL without an account (rate-limited per IP). Returns a short URL valid for 1 hour.
//	@Tags			demo
//	@Accept			json
//	@Produce		json
//	@Param			body	body		request.DemoShortenRequest	true	"URL to shorten"
//	@Success		201		{object}	transport.StandardResponse
//	@Failure		400		{object}	transport.ErrorResponse
//	@Failure		429		{object}	transport.ErrorResponse
//	@Router			/api/v1/demo/shorten [post]
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
