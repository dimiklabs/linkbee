package demo

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go/valkeycompat"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/response"
	"github.com/shafikshaon/linkbee/util"
)

const (
	demoRateLimitKeyPrefix = "demo:rl:"
	demoSlugCounterKey     = "linkbee:slug:counter"
)

type DemoServiceI interface {
	ShortenURL(ctx context.Context, destinationURL, ipAddress string) (*response.DemoShortenResponse, *dto.ServiceError)
}

type demoService struct {
	linkRepo repository.LinkRepositoryI
	cache    valkeycompat.Cmdable
	slugGen  *util.SlugGenerator
	appCfg   *config.AppConfig
	linkCfg  *config.LinkConfig
}

func NewDemoService(linkRepo repository.LinkRepositoryI, cache valkeycompat.Cmdable, slugGen *util.SlugGenerator, appCfg *config.AppConfig, linkCfg *config.LinkConfig) DemoServiceI {
	return &demoService{
		linkRepo: linkRepo,
		cache:    cache,
		slugGen:  slugGen,
		appCfg:   appCfg,
		linkCfg:  linkCfg,
	}
}

func (s *demoService) ShortenURL(ctx context.Context, destinationURL, ipAddress string) (*response.DemoShortenResponse, *dto.ServiceError) {
	// Rate limit check: max N demo links per IP per day
	rateLimitKey := fmt.Sprintf("%s%s", demoRateLimitKeyPrefix, util.HashIP(ipAddress))
	count, _ := s.cache.Incr(ctx, rateLimitKey).Result()
	if count == 1 {
		s.cache.Expire(ctx, rateLimitKey, 86400*1000000000) // 24h in nanoseconds
	}
	if count > int64(s.linkCfg.DemoRateLimitPerIP) {
		return nil, dto.NewTooManyRequestsError(constant.ErrCodeRateLimited, "Demo rate limit reached. Please sign up for a free account.")
	}

	// Generate slug via shared counter + shuffled-alphabet encoding
	counter, err := s.cache.Incr(ctx, demoSlugCounterKey).Result()
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	slug := s.slugGen.FromCounter(counter)

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
