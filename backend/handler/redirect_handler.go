package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/repository"
	clickSvc "github.com/shafikshaon/shortlink/service/click"
	redirectSvc "github.com/shafikshaon/shortlink/service/redirect"
	"github.com/shafikshaon/shortlink/transport"
)

type RedirectHandler struct {
	redirectService redirectSvc.RedirectServiceI
	clickService    clickSvc.ClickServiceI
	linkRepo        repository.LinkRepositoryI
}

func NewRedirectHandler(
	redirectService redirectSvc.RedirectServiceI,
	clickService clickSvc.ClickServiceI,
	linkRepo repository.LinkRepositoryI,
) *RedirectHandler {
	return &RedirectHandler{
		redirectService: redirectService,
		clickService:    clickService,
		linkRepo:        linkRepo,
	}
}

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
		pwd := c.Query("pwd")
		if pwd == "" {
			transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "This link is password protected")
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
			transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, "Invalid link password")
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
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")
	go h.clickService.EnqueueClick(context.Background(), linkID, ip, userAgent, referrer, source)

	// Redirect
	redirectCode := http.StatusFound // 302 default
	if link.RedirectType == 301 {
		redirectCode = http.StatusMovedPermanently
	}
	c.Redirect(redirectCode, link.DestinationURL)
}
