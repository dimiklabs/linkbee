package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/repository"
	clickSvc "github.com/shafikshaon/shortlink/service/click"
	geoSvc "github.com/shafikshaon/shortlink/service/geo"
	pixelSvc "github.com/shafikshaon/shortlink/service/pixel"
	redirectSvc "github.com/shafikshaon/shortlink/service/redirect"
	webhookSvc "github.com/shafikshaon/shortlink/service/webhook"
	"github.com/shafikshaon/shortlink/transport"
)

type RedirectHandler struct {
	redirectService redirectSvc.RedirectServiceI
	clickService    clickSvc.ClickServiceI
	geoService      geoSvc.GeoServiceI
	pixelService    pixelSvc.PixelServiceI
	webhookService  webhookSvc.WebhookServiceI
	linkRepo        repository.LinkRepositoryI
}

func NewRedirectHandler(
	redirectService redirectSvc.RedirectServiceI,
	clickService clickSvc.ClickServiceI,
	geoService geoSvc.GeoServiceI,
	pixelService pixelSvc.PixelServiceI,
	webhookService webhookSvc.WebhookServiceI,
	linkRepo repository.LinkRepositoryI,
) *RedirectHandler {
	return &RedirectHandler{
		redirectService: redirectService,
		clickService:    clickService,
		geoService:      geoService,
		pixelService:    pixelService,
		webhookService:  webhookService,
		linkRepo:        linkRepo,
	}
}

// Redirect godoc
//
//	@Summary		Redirect short URL
//	@Description	Resolves a slug and redirects to the destination URL. Tracks the click asynchronously.
//	@Tags			redirect
//	@Param			slug	path	string	true	"Short URL slug"
//	@Param			pwd		query	string	false	"Password for protected links"
//	@Param			source	query	string	false	"Click source (e.g. 'qr')"
//	@Success		302
//	@Success		301
//	@Failure		401	{object}	transport.ErrorResponse
//	@Failure		404	{object}	transport.ErrorResponse
//	@Failure		410	{object}	transport.ErrorResponse
//	@Router			/{slug} [get]
// Redirect handles GET /:slug — cache-first lookup + async click tracking.
func (h *RedirectHandler) Redirect(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	link, svcErr := h.redirectService.GetCachedLink(ctx, slug)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	if !link.IsActive {
		transport.RespondWithError(c, http.StatusNotFound, constant.ErrCodeLinkDisabled, constant.ErrMsgLinkDisabled)
		return
	}

	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		transport.RespondWithError(c, http.StatusGone, constant.ErrCodeLinkExpired, constant.ErrMsgLinkExpired)
		return
	}

	if link.MaxClicks != nil && link.ClickCount >= *link.MaxClicks {
		transport.RespondWithError(c, http.StatusGone, constant.ErrCodeLinkExpired, "Link click limit reached")
		return
	}

	// Password-protected link
	if link.PasswordHash != "" {
		unlockURL := "/unlock/" + slug
		pwd := c.Query("pwd")
		if pwd == "" {
			c.Redirect(http.StatusFound, unlockURL)
			return
		}

		// Cache stores a "cached" sentinel — fetch the real bcrypt hash from DB
		passwordHash := link.PasswordHash
		if passwordHash == "cached" {
			fullLink, err := h.linkRepo.GetBySlug(ctx, slug)
			if err != nil {
				transport.RespondWithError(c, http.StatusInternalServerError, constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
				return
			}
			passwordHash = fullLink.PasswordHash
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(pwd)); err != nil {
			c.Redirect(http.StatusFound, unlockURL+"?error=invalid")
			return
		}
	}

	// Determine click source
	source := "web"
	if c.Query("source") == "qr" {
		source = "qr"
	}

	// Enqueue click asynchronously using background context so it outlives the request
	linkID := link.ID
	userID := link.UserID
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")
	go h.clickService.EnqueueClick(context.Background(), linkID, ip, userAgent, referrer, source)
	h.webhookService.Trigger(userID, webhookSvc.EventLinkClicked, map[string]string{
		"link_id": linkID.String(),
		"ip":      ip,
		"source":  source,
	})

	// Determine destination URL: split test → geo routing → default
	destURL := link.DestinationURL
	if link.IsSplitTest {
		if variantURL, svcErr := h.redirectService.PickSplitTestVariant(ctx, linkID); svcErr == nil && variantURL != "" {
			destURL = variantURL
		}
	} else if link.IsGeoRouting {
		country := h.detectCountry(c)
		if geoURL, svcErr := h.redirectService.ApplyGeoRouting(ctx, linkID, country); svcErr == nil && geoURL != "" {
			destURL = geoURL
		}
	}

	// If pixel tracking is enabled, serve an intermediate HTML page that fires
	// all tracking pixels and then redirects via JS — no direct HTTP redirect.
	if link.IsPixelTracking {
		if html, err := h.pixelService.RenderRedirectPage(ctx, linkID, destURL); err == nil && html != "" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
			return
		}
		// Fall through to direct redirect if pixel rendering fails
	}

	redirectCode := http.StatusFound // 302 default
	if link.RedirectType == 301 {
		redirectCode = http.StatusMovedPermanently
	}
	c.Redirect(redirectCode, destURL)
}

// detectCountry extracts the visitor's ISO 3166-1 alpha-2 country code from
// CDN / proxy headers, then falls back to MaxMind DB lookup via geoService.
func (h *RedirectHandler) detectCountry(c *gin.Context) string {
	headers := map[string]string{
		"CF-IPCountry":    c.GetHeader("CF-IPCountry"),
		"X-GeoIP-Country": c.GetHeader("X-GeoIP-Country"),
		"X-Country-Code":  c.GetHeader("X-Country-Code"),
	}
	country := h.geoService.GetCountryCode(c.ClientIP(), headers)
	return strings.ToUpper(country)
}
