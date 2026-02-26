package redirect

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const (
	linkCachePrefix = "link:"
)

// CachedLink is the minimal link data cached in Valkey.
type CachedLink struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	DestURL      string     `json:"dest_url"`
	RedirectType int16      `json:"redirect_type"`
	IsActive     bool       `json:"is_active"`
	IsSplitTest  bool       `json:"is_split_test"`
	IsGeoRouting bool       `json:"is_geo_routing"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	MaxClicks    *int64     `json:"max_clicks,omitempty"`
	ClickCount   int64      `json:"click_count"`
	HasPassword  bool       `json:"has_password"`
}

type RedirectServiceI interface {
	GetCachedLink(ctx context.Context, slug string) (*model.Link, *dto.ServiceError)
	PickSplitTestVariant(ctx context.Context, linkID uuid.UUID) (string, *dto.ServiceError)
	ApplyGeoRouting(ctx context.Context, linkID uuid.UUID, countryCode string) (string, *dto.ServiceError)
}

type redirectService struct {
	linkRepo    repository.LinkRepositoryI
	variantRepo repository.LinkVariantRepositoryI
	geoRuleRepo repository.LinkGeoRuleRepositoryI
	cache       valkeycompat.Cmdable
	cacheTTL    time.Duration
}

func NewRedirectService(
	linkRepo repository.LinkRepositoryI,
	variantRepo repository.LinkVariantRepositoryI,
	geoRuleRepo repository.LinkGeoRuleRepositoryI,
	cache valkeycompat.Cmdable,
	cacheTTLSeconds int,
) RedirectServiceI {
	return &redirectService{
		linkRepo:    linkRepo,
		variantRepo: variantRepo,
		geoRuleRepo: geoRuleRepo,
		cache:       cache,
		cacheTTL:    time.Duration(cacheTTLSeconds) * time.Second,
	}
}

func (s *redirectService) GetCachedLink(ctx context.Context, slug string) (*model.Link, *dto.ServiceError) {
	cacheKey := fmt.Sprintf("%s%s", linkCachePrefix, slug)

	// Try Valkey cache first
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var cl CachedLink
		if jsonErr := json.Unmarshal([]byte(cached), &cl); jsonErr == nil {
			link := &model.Link{
				DestinationURL: cl.DestURL,
				RedirectType:   cl.RedirectType,
				IsActive:       cl.IsActive,
				IsSplitTest:    cl.IsSplitTest,
				IsGeoRouting:   cl.IsGeoRouting,
				ExpiresAt:      cl.ExpiresAt,
				MaxClicks:      cl.MaxClicks,
				ClickCount:     cl.ClickCount,
				PasswordHash:   "",
			}
			if cl.HasPassword {
				link.PasswordHash = "cached" // sentinel — actual hash not cached
			}
			_ = link.ID.UnmarshalText([]byte(cl.ID))
			_ = link.UserID.UnmarshalText([]byte(cl.UserID))
			return link, nil
		}
	}

	// Cache miss — fetch from PostgreSQL
	link, dbErr := s.linkRepo.GetBySlug(ctx, slug)
	if dbErr != nil {
		if dbErr == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		logger.ErrorCtx(ctx, "Failed to fetch link from DB", zap.String("slug", slug), zap.Error(dbErr))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Warm the cache
	s.warmCache(ctx, cacheKey, link)

	return link, nil
}

// PickSplitTestVariant selects a variant by weighted random and increments its click count.
// Returns the chosen destination URL, or "" if no variants exist (caller falls back to original URL).
func (s *redirectService) PickSplitTestVariant(ctx context.Context, linkID uuid.UUID) (string, *dto.ServiceError) {
	variants, err := s.variantRepo.GetByLinkID(ctx, linkID)
	if err != nil {
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if len(variants) == 0 {
		return "", nil
	}

	// Weighted random selection
	total := 0
	for _, v := range variants {
		total += v.Weight
	}

	r := rand.Intn(total)
	cumulative := 0
	var chosen *model.LinkVariant
	for i := range variants {
		cumulative += variants[i].Weight
		if r < cumulative {
			chosen = &variants[i]
			break
		}
	}
	if chosen == nil {
		chosen = &variants[len(variants)-1]
	}

	// Increment variant click count asynchronously
	variantID := chosen.ID
	go func() {
		_ = s.variantRepo.IncrementClickCount(context.Background(), variantID)
	}()

	return chosen.DestinationURL, nil
}

// ApplyGeoRouting looks up the first geo rule matching the given country code.
// Returns the destination URL for that rule, or "" if no rule matches (caller falls back).
func (s *redirectService) ApplyGeoRouting(ctx context.Context, linkID uuid.UUID, countryCode string) (string, *dto.ServiceError) {
	if countryCode == "" {
		return "", nil
	}
	rule, err := s.geoRuleRepo.GetByLinkIDAndCountry(ctx, linkID, countryCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil // no rule for this country — fall through to default URL
		}
		return "", dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return rule.DestinationURL, nil
}

func (s *redirectService) warmCache(ctx context.Context, key string, link *model.Link) {
	cl := CachedLink{
		ID:           link.ID.String(),
		UserID:       link.UserID.String(),
		DestURL:      link.DestinationURL,
		RedirectType: link.RedirectType,
		IsActive:     link.IsActive,
		IsSplitTest:  link.IsSplitTest,
		IsGeoRouting: link.IsGeoRouting,
		ExpiresAt:    link.ExpiresAt,
		MaxClicks:    link.MaxClicks,
		ClickCount:   link.ClickCount,
		HasPassword:  link.PasswordHash != "",
	}

	data, err := json.Marshal(cl)
	if err != nil {
		logger.Error("Failed to marshal link for cache", zap.Error(err))
		return
	}

	if err := s.cache.Set(ctx, key, string(data), s.cacheTTL).Err(); err != nil {
		logger.Error("Failed to warm link cache", zap.Error(err))
	}
}
