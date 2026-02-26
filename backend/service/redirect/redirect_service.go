package redirect

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

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
	DestURL      string     `json:"dest_url"`
	RedirectType int16      `json:"redirect_type"`
	IsActive     bool       `json:"is_active"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	MaxClicks    *int64     `json:"max_clicks,omitempty"`
	ClickCount   int64      `json:"click_count"`
	HasPassword  bool       `json:"has_password"`
}

type RedirectServiceI interface {
	GetCachedLink(ctx context.Context, slug string) (*model.Link, *dto.ServiceError)
}

type redirectService struct {
	linkRepo repository.LinkRepositoryI
	cache    valkeycompat.Cmdable
	cacheTTL time.Duration
}

func NewRedirectService(linkRepo repository.LinkRepositoryI, cache valkeycompat.Cmdable, cacheTTLSeconds int) RedirectServiceI {
	return &redirectService{
		linkRepo: linkRepo,
		cache:    cache,
		cacheTTL: time.Duration(cacheTTLSeconds) * time.Second,
	}
}

func (s *redirectService) GetCachedLink(ctx context.Context, slug string) (*model.Link, *dto.ServiceError) {
	cacheKey := fmt.Sprintf("%s%s", linkCachePrefix, slug)

	// Try Valkey cache first
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var cl CachedLink
		if jsonErr := json.Unmarshal([]byte(cached), &cl); jsonErr == nil {
			// Reconstruct minimal link from cache
			link := &model.Link{
				DestinationURL: cl.DestURL,
				RedirectType:   cl.RedirectType,
				IsActive:       cl.IsActive,
				ExpiresAt:      cl.ExpiresAt,
				MaxClicks:      cl.MaxClicks,
				ClickCount:     cl.ClickCount,
				PasswordHash:   "",
			}
			if cl.HasPassword {
				link.PasswordHash = "cached" // sentinel — actual hash not cached
			}
			// Parse UUID
			_ = link.ID.UnmarshalText([]byte(cl.ID))
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

func (s *redirectService) warmCache(ctx context.Context, key string, link *model.Link) {
	cl := CachedLink{
		ID:           link.ID.String(),
		DestURL:      link.DestinationURL,
		RedirectType: link.RedirectType,
		IsActive:     link.IsActive,
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
