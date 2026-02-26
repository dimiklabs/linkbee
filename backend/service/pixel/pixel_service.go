package pixel

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

// PixelResponse is the outward-facing pixel representation.
type PixelResponse struct {
	ID           uuid.UUID `json:"id"`
	LinkID       uuid.UUID `json:"link_id"`
	PixelType    string    `json:"pixel_type"`
	PixelID      string    `json:"pixel_id,omitempty"`
	CustomScript string    `json:"custom_script,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// PixelServiceI manages retargeting pixels and renders the redirect page.
type PixelServiceI interface {
	List(ctx context.Context, linkID, userID uuid.UUID) ([]PixelResponse, *dto.ServiceError)
	Create(ctx context.Context, linkID, userID uuid.UUID, pixelType, pixelID, customScript string) (*PixelResponse, *dto.ServiceError)
	Delete(ctx context.Context, pixelID, linkID, userID uuid.UUID) *dto.ServiceError
	TogglePixelTracking(ctx context.Context, linkID, userID uuid.UUID, enabled bool) *dto.ServiceError
	// RenderRedirectPage generates the intermediate HTML page with pixel scripts
	// and a JS redirect to destURL. Returns empty string if no active pixels exist.
	RenderRedirectPage(ctx context.Context, linkID uuid.UUID, destURL string) (string, error)
}

type pixelService struct {
	pixelRepo repository.RetargetingPixelRepositoryI
	linkRepo  repository.LinkRepositoryI
}

func NewPixelService(pixelRepo repository.RetargetingPixelRepositoryI, linkRepo repository.LinkRepositoryI) PixelServiceI {
	return &pixelService{pixelRepo: pixelRepo, linkRepo: linkRepo}
}

func (s *pixelService) ownerLink(ctx context.Context, linkID, userID uuid.UUID) error {
	link, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	if link.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *pixelService) List(ctx context.Context, linkID, userID uuid.UUID) ([]PixelResponse, *dto.ServiceError) {
	if err := s.ownerLink(ctx, linkID, userID); err != nil {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}
	pixels, err := s.pixelRepo.GetByLinkID(ctx, linkID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	out := make([]PixelResponse, len(pixels))
	for i, p := range pixels {
		out[i] = toResponse(&p)
	}
	return out, nil
}

func (s *pixelService) Create(ctx context.Context, linkID, userID uuid.UUID, pixelType, pixelID, customScript string) (*PixelResponse, *dto.ServiceError) {
	if err := s.ownerLink(ctx, linkID, userID); err != nil {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}

	// Validate pixel type
	validTypes := map[string]bool{
		model.PixelTypeFacebook:  true,
		model.PixelTypeGoogleAds: true,
		model.PixelTypeTikTok:    true,
		model.PixelTypeLinkedIn:  true,
		model.PixelTypeCustom:    true,
	}
	if !validTypes[pixelType] {
		return nil, dto.NewBadRequestError(constant.ErrCodeValidationError, "invalid pixel_type; must be one of: facebook, google_ads, tiktok, linkedin, custom")
	}
	if pixelType != model.PixelTypeCustom && pixelID == "" {
		return nil, dto.NewBadRequestError(constant.ErrCodeValidationError, "pixel_id is required for this pixel type")
	}
	if pixelType == model.PixelTypeCustom && customScript == "" {
		return nil, dto.NewBadRequestError(constant.ErrCodeValidationError, "custom_script is required for custom pixel type")
	}

	p := &model.RetargetingPixel{
		LinkID:       linkID,
		PixelType:    pixelType,
		PixelID:      pixelID,
		CustomScript: customScript,
		IsActive:     true,
	}
	if err := s.pixelRepo.Create(ctx, p); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	resp := toResponse(p)
	return &resp, nil
}

func (s *pixelService) Delete(ctx context.Context, pixelID, linkID, userID uuid.UUID) *dto.ServiceError {
	if err := s.ownerLink(ctx, linkID, userID); err != nil {
		return dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}
	if err := s.pixelRepo.Delete(ctx, pixelID, linkID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Pixel not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

func (s *pixelService) TogglePixelTracking(ctx context.Context, linkID, userID uuid.UUID, enabled bool) *dto.ServiceError {
	link, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if link.UserID != userID {
		return dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}
	link.IsPixelTracking = enabled
	if err := s.linkRepo.Update(ctx, link); err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

// RenderRedirectPage generates an intermediate HTML page that fires all active
// tracking pixels for the link and then redirects to destURL via JavaScript.
func (s *pixelService) RenderRedirectPage(ctx context.Context, linkID uuid.UUID, destURL string) (string, error) {
	pixels, err := s.pixelRepo.GetActiveByLinkID(ctx, linkID)
	if err != nil {
		return "", err
	}

	var scripts strings.Builder
	for _, p := range pixels {
		scripts.WriteString(buildPixelScript(p))
		scripts.WriteString("\n")
	}

	// JSON-encode the destination URL so it's safe to embed in a JS string literal.
	destJSON, _ := json.Marshal(destURL)
	// HTML-escape the destination URL for use in the meta refresh attribute.
	destHTMLAttr := html.EscapeString(destURL)

	page := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="refresh" content="0;url=%s">
<title>Redirecting...</title>
<style>body{margin:0;background:#fff;font-family:sans-serif;display:flex;align-items:center;justify-content:center;height:100vh}</style>
</head>
<body>
%s
<script>window.location.replace(%s);</script>
</body>
</html>`, destHTMLAttr, scripts.String(), string(destJSON))

	return page, nil
}

// buildPixelScript returns the HTML/JS snippet for a single pixel.
func buildPixelScript(p model.RetargetingPixel) string {
	switch p.PixelType {
	case model.PixelTypeFacebook:
		return buildFacebookScript(p.PixelID)
	case model.PixelTypeGoogleAds:
		return buildGoogleAdsScript(p.PixelID)
	case model.PixelTypeTikTok:
		return buildTikTokScript(p.PixelID)
	case model.PixelTypeLinkedIn:
		return buildLinkedInScript(p.PixelID)
	case model.PixelTypeCustom:
		// Custom scripts are inserted verbatim — user takes full responsibility.
		return p.CustomScript
	default:
		return ""
	}
}

func jsStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b) // produces a JSON-safe quoted string, e.g. "\"abc\""
}

func buildFacebookScript(pixelID string) string {
	id := jsStr(pixelID)
	idHTML := html.EscapeString(pixelID)
	return fmt.Sprintf(`<script>!function(f,b,e,v,n,t,s){if(f.fbq)return;n=f.fbq=function(){n.callMethod?n.callMethod.apply(n,arguments):n.queue.push(arguments)};if(!f._fbq)f._fbq=n;n.push=n;n.loaded=!0;n.version='2.0';n.queue=[];t=b.createElement(e);t.async=!0;t.src=v;s=b.getElementsByTagName(e)[0];s.parentNode.insertBefore(t,s)}(window,document,'script','https://connect.facebook.net/en_US/fbevents.js');fbq('init',%s);fbq('track','PageView');</script><noscript><img height="1" width="1" style="display:none" src="https://www.facebook.com/tr?id=%s&amp;ev=PageView&amp;noscript=1"/></noscript>`, id, idHTML)
}

func buildGoogleAdsScript(pixelID string) string {
	id := jsStr(pixelID)
	idHTML := html.EscapeString(pixelID)
	return fmt.Sprintf(`<script async src="https://www.googletagmanager.com/gtag/js?id=%s"></script><script>window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}gtag('js',new Date());gtag('config',%s);</script>`, idHTML, id)
}

func buildTikTokScript(pixelID string) string {
	id := jsStr(pixelID)
	return fmt.Sprintf(`<script>!function(w,d,t){w.TiktokAnalyticsObject=t;var ttq=w[t]=w[t]||[];ttq.methods=["page","track","identify","instances","debug","on","off","once","ready","alias","group","enableCookie","disableCookie"],ttq.setAndDefer=function(t,e){t[e]=function(){t.push([e].concat(Array.prototype.slice.call(arguments,0)))}};for(var i=0;i<ttq.methods.length;i++)ttq.setAndDefer(ttq,ttq.methods[i]);ttq.instance=function(t){for(var e=ttq._i[t]||[],n=0;n<ttq.methods.length;n++)ttq.setAndDefer(e,ttq.methods[n]);return e},ttq.load=function(e,n){var i="https://analytics.tiktok.com/i18n/pixel/events.js";ttq._i=ttq._i||{},ttq._i[e]=[],ttq._t=ttq._t||{},ttq._t[e]=+new Date,ttq._o=ttq._o||{},ttq._o[e]=n||{};var r=document.createElement("script");r.type="text/javascript",r.async=!0,r.src=i+"?sdkid="+e+"&lib="+t;var a=document.getElementsByTagName("script")[0];a.parentNode.insertBefore(r,a)};ttq.load(%s);ttq.page();}(window,document,'ttq');</script>`, id)
}

func buildLinkedInScript(pixelID string) string {
	id := jsStr(pixelID)
	return fmt.Sprintf(`<script>_linkedin_partner_id=%s;window._linkedin_data_partner_ids=window._linkedin_data_partner_ids||[];window._linkedin_data_partner_ids.push(_linkedin_partner_id);</script><script>(function(l){if(!l){window.lintrk=function(a,b){window.lintrk.q.push([a,b])};window.lintrk.q=[]}var s=document.getElementsByTagName("script")[0];var b=document.createElement("script");b.type="text/javascript";b.async=true;b.src="https://snap.licdn.com/li.lms-analytics/insight.min.js";s.parentNode.insertBefore(b,s);})(window.lintrk);</script><noscript><img height="1" width="1" style="display:none;" alt="" src="https://px.ads.linkedin.com/collect/?pid=%s&amp;fmt=gif"/></noscript>`, id, html.EscapeString(pixelID))
}

func toResponse(p *model.RetargetingPixel) PixelResponse {
	return PixelResponse{
		ID:           p.ID,
		LinkID:       p.LinkID,
		PixelType:    p.PixelType,
		PixelID:      p.PixelID,
		CustomScript: p.CustomScript,
		IsActive:     p.IsActive,
		CreatedAt:    p.CreatedAt,
	}
}
