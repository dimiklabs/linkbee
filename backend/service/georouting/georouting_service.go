package georouting

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/response"
)

// GeoRuleResponse is the API representation of a single geo rule.
type GeoRuleResponse struct {
	ID             uuid.UUID `json:"id"`
	LinkID         uuid.UUID `json:"link_id"`
	CountryCode    string    `json:"country_code"`
	DestinationURL string    `json:"destination_url"`
	Priority       int       `json:"priority"`
}

type GeoRoutingServiceI interface {
	ListRules(ctx context.Context, linkID, userID uuid.UUID) ([]GeoRuleResponse, *dto.ServiceError)
	CreateRule(ctx context.Context, linkID, userID uuid.UUID, countryCode, destURL string, priority int) (*GeoRuleResponse, *dto.ServiceError)
	UpdateRule(ctx context.Context, ruleID, linkID, userID uuid.UUID, countryCode, destURL string, priority int) (*GeoRuleResponse, *dto.ServiceError)
	DeleteRule(ctx context.Context, ruleID, linkID, userID uuid.UUID) *dto.ServiceError
	ToggleGeoRouting(ctx context.Context, linkID, userID uuid.UUID, enabled bool) (*response.LinkResponse, *dto.ServiceError)
}

type geoRoutingService struct {
	geoRuleRepo repository.LinkGeoRuleRepositoryI
	linkRepo    repository.LinkRepositoryI
}

func NewGeoRoutingService(geoRuleRepo repository.LinkGeoRuleRepositoryI, linkRepo repository.LinkRepositoryI) GeoRoutingServiceI {
	return &geoRoutingService{
		geoRuleRepo: geoRuleRepo,
		linkRepo:    linkRepo,
	}
}

// ownerLink fetches a link and asserts ownership, returning a service error on failure.
func (s *geoRoutingService) ownerLink(ctx context.Context, linkID, userID uuid.UUID) (*model.Link, *dto.ServiceError) {
	link, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if link.UserID != userID {
		return nil, dto.NewForbiddenError(constant.ErrCodeForbidden, constant.ErrMsgForbidden)
	}
	return link, nil
}

func (s *geoRoutingService) ListRules(ctx context.Context, linkID, userID uuid.UUID) ([]GeoRuleResponse, *dto.ServiceError) {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return nil, svcErr
	}
	rules, err := s.geoRuleRepo.GetByLinkID(ctx, linkID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	out := make([]GeoRuleResponse, 0, len(rules))
	for _, r := range rules {
		out = append(out, toGeoRuleResponse(r))
	}
	return out, nil
}

func (s *geoRoutingService) CreateRule(ctx context.Context, linkID, userID uuid.UUID, countryCode, destURL string, priority int) (*GeoRuleResponse, *dto.ServiceError) {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return nil, svcErr
	}
	rule := &model.LinkGeoRule{
		LinkID:         linkID,
		CountryCode:    strings.ToUpper(countryCode),
		DestinationURL: destURL,
		Priority:       priority,
	}
	if err := s.geoRuleRepo.Create(ctx, rule); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	resp := toGeoRuleResponse(*rule)
	return &resp, nil
}

func (s *geoRoutingService) UpdateRule(ctx context.Context, ruleID, linkID, userID uuid.UUID, countryCode, destURL string, priority int) (*GeoRuleResponse, *dto.ServiceError) {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return nil, svcErr
	}
	// Fetch the existing rule to confirm it belongs to this link
	rules, err := s.geoRuleRepo.GetByLinkID(ctx, linkID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	var rule *model.LinkGeoRule
	for i := range rules {
		if rules[i].ID == ruleID {
			rule = &rules[i]
			break
		}
	}
	if rule == nil {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, "Geo rule not found")
	}
	if countryCode != "" {
		rule.CountryCode = strings.ToUpper(countryCode)
	}
	if destURL != "" {
		rule.DestinationURL = destURL
	}
	rule.Priority = priority

	if err := s.geoRuleRepo.Update(ctx, rule); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	resp := toGeoRuleResponse(*rule)
	return &resp, nil
}

func (s *geoRoutingService) DeleteRule(ctx context.Context, ruleID, linkID, userID uuid.UUID) *dto.ServiceError {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return svcErr
	}
	if err := s.geoRuleRepo.Delete(ctx, ruleID, linkID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeLinkNotFound, "Geo rule not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

func (s *geoRoutingService) ToggleGeoRouting(ctx context.Context, linkID, userID uuid.UUID, enabled bool) (*response.LinkResponse, *dto.ServiceError) {
	link, svcErr := s.ownerLink(ctx, linkID, userID)
	if svcErr != nil {
		return nil, svcErr
	}
	link.IsGeoRouting = enabled
	if err := s.linkRepo.Update(ctx, link); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	resp := toLinkResponseMinimal(link)
	return &resp, nil
}

func toGeoRuleResponse(r model.LinkGeoRule) GeoRuleResponse {
	return GeoRuleResponse{
		ID:             r.ID,
		LinkID:         r.LinkID,
		CountryCode:    r.CountryCode,
		DestinationURL: r.DestinationURL,
		Priority:       r.Priority,
	}
}

func toLinkResponseMinimal(link *model.Link) response.LinkResponse {
	return response.LinkResponse{
		ID:           link.ID,
		Slug:         link.Slug,
		IsGeoRouting: link.IsGeoRouting,
		UpdatedAt:    link.UpdatedAt,
	}
}
