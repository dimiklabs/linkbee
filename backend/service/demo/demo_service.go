package demo

import (
	"context"
	"fmt"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/response"
	"github.com/shafikshaon/shortlink/util"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

const demoRateLimitKeyPrefix = "demo:rl:"

type DemoServiceI interface {
	ShortenURL(ctx context.Context, destinationURL, ipAddress string) (*response.DemoShortenResponse, *dto.ServiceError)
}

type demoService struct {
	linkRepo repository.LinkRepositoryI
	cache    valkeycompat.Cmdable
	appCfg   *config.AppConfig
	linkCfg  *config.LinkConfig
}

func NewDemoService(linkRepo repository.LinkRepositoryI, cache valkeycompat.Cmdable, appCfg *config.AppConfig, linkCfg *config.LinkConfig) DemoServiceI {
	return &demoService{
		linkRepo: linkRepo,
		cache:    cache,
		appCfg:   appCfg,
		linkCfg:  linkCfg,
	}
}

func (s *demoService) ShortenURL(ctx context.Context, destinationURL, ipAddress string) (*response.DemoShortenResponse, *dto.ServiceError) {
	// Rate limit check: max N demo links per IP per day
	rateLimitKey := fmt.Sprintf("%s%s", demoRateLimitKeyPrefix, util.HashIP(ipAddress))
	count, _ := s.cache.Incr(ctx, rateLimitKey).Result()
	if count == 1 {
		// Set TTL of 24 hours on first request
		s.cache.Expire(ctx, rateLimitKey, 86400*1000000000) // 24h in nanoseconds
	}
	if count > int64(s.linkCfg.DemoRateLimitPerIP) {
		return nil, dto.NewTooManyRequestsError(constant.ErrCodeRateLimited, "Demo rate limit reached. Please sign up for a free account.")
	}

	// Generate slug
	var slug string
	for attempt := 0; attempt < 5; attempt++ {
		generated, err := util.GenerateSlug(s.linkCfg.SlugLength)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		exists, err := s.linkRepo.SlugExists(ctx, generated)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		if !exists {
			slug = generated
			break
		}
	}
	if slug == "" {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "Failed to generate unique slug")
	}

	// Create link without a user (demo link - user_id will be a zero UUID)
	link := &model.Link{
		Slug:           slug,
		DestinationURL: destinationURL,
		Title:          "Demo link",
		RedirectType:   int16(s.linkCfg.DefaultRedirectType),
		IsActive:       true,
	}

	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	shortURL := fmt.Sprintf("%s/%s", s.appCfg.BaseDomain, slug)
	return &response.DemoShortenResponse{
		ShortURL:       shortURL,
		Slug:           slug,
		DestinationURL: destinationURL,
	}, nil
}
