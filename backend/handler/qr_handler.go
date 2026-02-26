package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/middlewares"
	linkSvc "github.com/shafikshaon/shortlink/service/link"
	qrSvc "github.com/shafikshaon/shortlink/service/qr"
	"github.com/shafikshaon/shortlink/transport"
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

// GetQRCode handles GET /api/v1/links/:id/qr
//
// Query parameters:
//
//	format – "png" (default, inline for preview) | "svg" (attachment download)
//	fg     – foreground hex without '#', default "000000"
//	bg     – background hex without '#', default "ffffff"
//	size   – pixel dimension (64–1024), default 256
//	ec     – error correction: L / M / Q / H, default "M"
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
