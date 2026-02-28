package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	linkSvc "github.com/shafikshaon/linkbee/service/link"
	qrSvc "github.com/shafikshaon/linkbee/service/qr"
	"github.com/shafikshaon/linkbee/transport"
)

type QRHandler struct {
	qrService   qrSvc.QRServiceI
	linkService linkSvc.LinkServiceI
	appCfg      *config.AppConfig
}

func NewQRHandler(qrService qrSvc.QRServiceI, linkService linkSvc.LinkServiceI, appCfg *config.AppConfig) *QRHandler {
	return &QRHandler{
		qrService:   qrService,
		linkService: linkService,
		appCfg:      appCfg,
	}
}

// GetQRCode godoc
//
//	@Summary		Get QR code for a link
//	@Description	Returns a QR code image in PNG (default, inline) or SVG (download) format. Supports custom colours, size, and error correction level.
//	@Tags			links
//	@Produce		image/png
//	@Security		BearerAuth
//	@Security		APIKeyAuth
//	@Param			id		path	string	true	"Link UUID"
//	@Param			format	query	string	false	"Output format: png (default) or svg"
//	@Param			fg		query	string	false	"Foreground hex colour without '#' (default 000000)"
//	@Param			bg		query	string	false	"Background hex colour without '#' (default ffffff)"
//	@Param			size	query	int		false	"Image dimension in pixels 64–1024 (default 256)"
//	@Param			ec		query	string	false	"Error correction level: L, M, Q, H (default M)"
//	@Success		200
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Router			/api/v1/links/{id}/qr [get]
func (h *QRHandler) GetQRCode(c *gin.Context) {
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

	// Verify ownership
	link, svcErr := h.linkService.GetLink(ctx, id, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	// Parse customisation options
	opts := qrSvc.QROptions{
		ForegroundHex:   c.DefaultQuery("fg", "000000"),
		BackgroundHex:   c.DefaultQuery("bg", "ffffff"),
		ErrorCorrection: c.DefaultQuery("ec", "M"),
	}
	if sizeStr := c.Query("size"); sizeStr != "" {
		if sz, err := strconv.Atoi(sizeStr); err == nil {
			opts.Size = sz
		}
	}
	if opts.Size == 0 {
		opts.Size = 256
	}

	shortURL := fmt.Sprintf("%s/%s", h.appCfg.BaseDomain, link.Slug)

	format := strings.ToLower(c.DefaultQuery("format", "png"))
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	switch format {
	case "svg":
		svgData, err := h.qrService.GenerateSVG(shortURL, opts)
		if err != nil {
			transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
			return
		}
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"qr-%s.svg\"", link.Slug))
		c.Data(http.StatusOK, "image/svg+xml; charset=utf-8", svgData)

	default: // "png"
		pngData, err := h.qrService.GenerateCustomPNG(shortURL, opts)
		if err != nil {
			transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
			return
		}
		// Use inline so the <img> tag preview works; the browser still downloads
		// when the frontend <a> tag carries the `download` attribute.
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"qr-%s.png\"", link.Slug))
		c.Data(http.StatusOK, "image/png", pngData)
	}
}
