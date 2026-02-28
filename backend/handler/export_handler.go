package handler

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/middlewares"
	"github.com/shafikshaon/linkbee/repository"
)

// exportProfile is the user profile section of the data export.
type exportProfile struct {
	ID             string  `json:"id"`
	Email          string  `json:"email"`
	FirstName      string  `json:"first_name,omitempty"`
	LastName       string  `json:"last_name,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	Role           string  `json:"role"`
	Status         string  `json:"status"`
	EmailVerified  bool    `json:"email_verified"`
	AuthProvider   string  `json:"auth_provider"`
	CreatedAt      string  `json:"created_at"`
	LastLogin      *string `json:"last_login,omitempty"`
}

// exportLink is a single link entry in the data export.
type exportLink struct {
	ID             string   `json:"id"`
	Slug           string   `json:"slug"`
	ShortURL       string   `json:"short_url"`
	DestinationURL string   `json:"destination_url"`
	Title          string   `json:"title,omitempty"`
	ClickCount     int64    `json:"click_count"`
	IsActive       bool     `json:"is_active"`
	Tags           []string `json:"tags"`
	CreatedAt      string   `json:"created_at"`
	ExpiresAt      *string  `json:"expires_at,omitempty"`
}

// dataExport is the top-level export document.
type dataExport struct {
	ExportedAt  string          `json:"exported_at"`
	Profile     exportProfile   `json:"profile"`
	TotalLinks  int             `json:"total_links"`
	TotalClicks int64           `json:"total_clicks"`
	Links       []exportLink    `json:"links"`
}

type ExportHandler struct {
	userRepo repository.UserRepositoryI
	linkRepo repository.LinkRepositoryI
	appCfg   *config.AppConfig
}

func NewExportHandler(userRepo repository.UserRepositoryI, linkRepo repository.LinkRepositoryI, appCfg *config.AppConfig) *ExportHandler {
	return &ExportHandler{userRepo: userRepo, linkRepo: linkRepo, appCfg: appCfg}
}

// ExportData streams a JSON file containing all of the authenticated user's data.
func (h *ExportHandler) ExportData(c *gin.Context) {
	ctx := c.Request.Context()

	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHORIZED", "description": "unauthorized"})
		return
	}

	// ── Profile ──────────────────────────────────────────────────────────────
	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_SERVER_ERROR", "description": "failed to fetch profile"})
		return
	}

	prof := exportProfile{
		ID:            user.ID.String(),
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		AuthProvider:  user.AuthProvider,
		CreatedAt:     user.CreatedAt.UTC().Format(time.RFC3339),
	}
	if user.LastLogin != nil {
		s := user.LastLogin.UTC().Format(time.RFC3339)
		prof.LastLogin = &s
	}

	// ── Links ─────────────────────────────────────────────────────────────────
	links, err := h.linkRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_SERVER_ERROR", "description": "failed to fetch links"})
		return
	}

	exportLinks := make([]exportLink, 0, len(links))
	var totalClicks int64
	for _, l := range links {
		tags := []string{}
		if l.Tags != nil {
			tags = l.Tags
		}
		el := exportLink{
			ID:             l.ID.String(),
			Slug:           l.Slug,
			ShortURL:       fmt.Sprintf("%s/%s", h.appCfg.BaseDomain, l.Slug),
			DestinationURL: l.DestinationURL,
			Title:          l.Title,
			ClickCount:     l.ClickCount,
			IsActive:       l.IsActive,
			Tags:           tags,
			CreatedAt:      l.CreatedAt.UTC().Format(time.RFC3339),
		}
		if l.ExpiresAt != nil {
			s := l.ExpiresAt.UTC().Format(time.RFC3339)
			el.ExpiresAt = &s
		}
		exportLinks = append(exportLinks, el)
		totalClicks += l.ClickCount
	}

	// ── Build document ────────────────────────────────────────────────────────
	doc := dataExport{
		ExportedAt:  time.Now().UTC().Format(time.RFC3339),
		Profile:     prof,
		TotalLinks:  len(exportLinks),
		TotalClicks: totalClicks,
		Links:       exportLinks,
	}

	payload, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_SERVER_ERROR", "description": "failed to encode export"})
		return
	}

	filename := fmt.Sprintf("linkbee-data-%s.json", time.Now().UTC().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "application/json; charset=utf-8", payload)
}

// ExportLinksCSV streams a CSV file containing all of the authenticated user's links.
// Columns: slug, short_url, destination_url, title, click_count, is_active, is_starred,
//
//	redirect_type, tags, expires_at, max_clicks, password_protected,
//	utm_source, utm_medium, utm_campaign, health_status, created_at
func (h *ExportHandler) ExportLinksCSV(c *gin.Context) {
	ctx := c.Request.Context()

	rawID, _ := c.Get(middlewares.ContextKeyUserID)
	userID, ok := rawID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error_code": "UNAUTHORIZED", "description": "unauthorized"})
		return
	}

	links, err := h.linkRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error_code": "INTERNAL_SERVER_ERROR", "description": "failed to fetch links"})
		return
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Header row
	_ = w.Write([]string{
		"slug", "short_url", "destination_url", "title",
		"click_count", "is_active", "is_starred", "redirect_type",
		"tags", "expires_at", "max_clicks", "password_protected",
		"utm_source", "utm_medium", "utm_campaign", "health_status", "created_at",
	})

	for _, l := range links {
		expiresAt := ""
		if l.ExpiresAt != nil {
			expiresAt = l.ExpiresAt.UTC().Format(time.RFC3339)
		}
		maxClicks := ""
		if l.MaxClicks != nil {
			maxClicks = strconv.FormatInt(*l.MaxClicks, 10)
		}
		tags := ""
		if len(l.Tags) > 0 {
			tags = strings.Join(l.Tags, "|")
		}
		_ = w.Write([]string{
			l.Slug,
			fmt.Sprintf("%s/%s", h.appCfg.BaseDomain, l.Slug),
			l.DestinationURL,
			l.Title,
			strconv.FormatInt(l.ClickCount, 10),
			strconv.FormatBool(l.IsActive),
			strconv.FormatBool(l.IsStarred),
			strconv.Itoa(int(l.RedirectType)),
			tags,
			expiresAt,
			maxClicks,
			strconv.FormatBool(l.PasswordHash != ""),
			l.UTMSource,
			l.UTMMedium,
			l.UTMCampaign,
			l.HealthStatus,
			l.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	w.Flush()

	filename := fmt.Sprintf("links-%s.csv", time.Now().UTC().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}
