package split

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/request"
	"github.com/shafikshaon/shortlink/response"
)

type SplitServiceI interface {
	ListVariants(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]*response.LinkVariantResponse, *dto.ServiceError)
	CreateVariant(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, req *request.CreateVariantRequest) (*response.LinkVariantResponse, *dto.ServiceError)
	UpdateVariant(ctx context.Context, variantID uuid.UUID, linkID uuid.UUID, userID uuid.UUID, req *request.UpdateVariantRequest) (*response.LinkVariantResponse, *dto.ServiceError)
	DeleteVariant(ctx context.Context, variantID uuid.UUID, linkID uuid.UUID, userID uuid.UUID) *dto.ServiceError
	ToggleSplitTest(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, enabled bool) (*response.LinkResponse, *dto.ServiceError)
}

type splitService struct {
	variantRepo repository.LinkVariantRepositoryI
	linkRepo    repository.LinkRepositoryI
}

func NewSplitService(variantRepo repository.LinkVariantRepositoryI, linkRepo repository.LinkRepositoryI) SplitServiceI {
	return &splitService{variantRepo: variantRepo, linkRepo: linkRepo}
}

func (s *splitService) ownerLink(ctx context.Context, linkID, userID uuid.UUID) (*model.Link, *dto.ServiceError) {
	link, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if link.UserID != userID {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}
	return link, nil
}

func (s *splitService) ListVariants(ctx context.Context, linkID uuid.UUID, userID uuid.UUID) ([]*response.LinkVariantResponse, *dto.ServiceError) {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return nil, svcErr
	}
	variants, err := s.variantRepo.GetByLinkID(ctx, linkID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	result := make([]*response.LinkVariantResponse, len(variants))
	for i := range variants {
		result[i] = toVariantResponse(&variants[i])
	}
	return result, nil
}

func (s *splitService) CreateVariant(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, req *request.CreateVariantRequest) (*response.LinkVariantResponse, *dto.ServiceError) {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return nil, svcErr
	}

	variant := &model.LinkVariant{
		LinkID:         linkID,
		Name:           req.Name,
		DestinationURL: req.DestinationURL,
		Weight:         req.Weight,
	}
	if err := s.variantRepo.Create(ctx, variant); err != nil {
		logger.ErrorCtx(ctx, "Failed to create variant", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return toVariantResponse(variant), nil
}

func (s *splitService) UpdateVariant(ctx context.Context, variantID uuid.UUID, linkID uuid.UUID, userID uuid.UUID, req *request.UpdateVariantRequest) (*response.LinkVariantResponse, *dto.ServiceError) {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return nil, svcErr
	}

	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError("VARIANT_NOT_FOUND", "Variant not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if variant.LinkID != linkID {
		return nil, dto.NewNotFoundError("VARIANT_NOT_FOUND", "Variant not found")
	}

	if req.Name != "" {
		variant.Name = req.Name
	}
	if req.DestinationURL != "" {
		variant.DestinationURL = req.DestinationURL
	}
	if req.Weight > 0 {
		variant.Weight = req.Weight
	}
	variant.UpdatedAt = time.Now()

	if err := s.variantRepo.Update(ctx, variant); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return toVariantResponse(variant), nil
}

func (s *splitService) DeleteVariant(ctx context.Context, variantID uuid.UUID, linkID uuid.UUID, userID uuid.UUID) *dto.ServiceError {
	if _, svcErr := s.ownerLink(ctx, linkID, userID); svcErr != nil {
		return svcErr
	}
	if err := s.variantRepo.Delete(ctx, variantID, linkID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError("VARIANT_NOT_FOUND", "Variant not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

func (s *splitService) ToggleSplitTest(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, enabled bool) (*response.LinkResponse, *dto.ServiceError) {
	link, svcErr := s.ownerLink(ctx, linkID, userID)
	if svcErr != nil {
		return nil, svcErr
	}

	link.IsSplitTest = enabled
	link.UpdatedAt = time.Now()

	if err := s.linkRepo.Update(ctx, link); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return toLinkResponseMinimal(link), nil
}

func toVariantResponse(v *model.LinkVariant) *response.LinkVariantResponse {
	return &response.LinkVariantResponse{
		ID:             v.ID,
		LinkID:         v.LinkID,
		Name:           v.Name,
		DestinationURL: v.DestinationURL,
		Weight:         v.Weight,
		ClickCount:     v.ClickCount,
		CreatedAt:      v.CreatedAt,
		UpdatedAt:      v.UpdatedAt,
	}
}

// toLinkResponseMinimal builds a minimal LinkResponse for ToggleSplitTest.
// Full response (with ShortURL etc.) is built by the link service; here we only
// need the IsSplitTest field update confirmed.
func toLinkResponseMinimal(link *model.Link) *response.LinkResponse {
	return &response.LinkResponse{
		ID:          link.ID,
		Slug:        link.Slug,
		IsSplitTest: link.IsSplitTest,
		IsActive:    link.IsActive,
		IsStarred:   link.IsStarred,
		UpdatedAt:   link.UpdatedAt,
	}
}
