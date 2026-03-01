package billing

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/repository"
)

// PlanEnforcerI checks whether a user is allowed to create a new resource
// given their current subscription plan.
type PlanEnforcerI interface {
	CheckLinkLimit(ctx context.Context, userID uuid.UUID) *dto.ServiceError
	CheckAPIKeyLimit(ctx context.Context, userID uuid.UUID) *dto.ServiceError
	CheckWebhookLimit(ctx context.Context, userID uuid.UUID) *dto.ServiceError
	CheckCustomSlug(ctx context.Context, userID uuid.UUID) *dto.ServiceError
	CheckAnalytics(ctx context.Context, userID uuid.UUID) *dto.ServiceError
}

type planEnforcer struct {
	subRepo     repository.SubscriptionRepositoryI
	linkRepo    repository.LinkRepositoryI
	apiKeyRepo  repository.APIKeyRepositoryI
	webhookRepo repository.WebhookRepositoryI
}

func NewPlanEnforcer(
	subRepo repository.SubscriptionRepositoryI,
	linkRepo repository.LinkRepositoryI,
	apiKeyRepo repository.APIKeyRepositoryI,
	webhookRepo repository.WebhookRepositoryI,
) PlanEnforcerI {
	return &planEnforcer{
		subRepo:     subRepo,
		linkRepo:    linkRepo,
		apiKeyRepo:  apiKeyRepo,
		webhookRepo: webhookRepo,
	}
}

// plan returns the PlanInfo for the user, defaulting to Free on any error.
func (e *planEnforcer) plan(ctx context.Context, userID uuid.UUID) PlanInfo {
	sub, err := e.subRepo.GetByUserID(ctx, userID)
	if err != nil || sub == nil {
		return GetPlan(PlanFree)
	}
	return GetPlan(sub.PlanID)
}

func (e *planEnforcer) CheckLinkLimit(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	plan := e.plan(ctx, userID)
	if plan.MaxLinks == -1 {
		return nil
	}
	count, err := e.linkRepo.CountByUserID(ctx, userID)
	if err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count links")
	}
	if count >= int64(plan.MaxLinks) {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			fmt.Sprintf("your %s plan allows up to %d links; upgrade to create more", plan.Name, plan.MaxLinks),
		)
	}
	return nil
}

func (e *planEnforcer) CheckAPIKeyLimit(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	plan := e.plan(ctx, userID)
	if plan.MaxAPIKeys == -1 {
		return nil
	}
	if plan.MaxAPIKeys == 0 {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			fmt.Sprintf("API keys are not available on the %s plan; upgrade to Pro or Growth", plan.Name),
		)
	}
	count, err := e.apiKeyRepo.CountByUserID(ctx, userID)
	if err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count API keys")
	}
	if count >= int64(plan.MaxAPIKeys) {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			fmt.Sprintf("your %s plan allows up to %d API keys; upgrade to create more", plan.Name, plan.MaxAPIKeys),
		)
	}
	return nil
}

func (e *planEnforcer) CheckAnalytics(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	plan := e.plan(ctx, userID)
	if plan.ID == PlanFree {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			fmt.Sprintf("analytics are not available on the %s plan; upgrade to Pro or Growth", plan.Name),
		)
	}
	return nil
}

func (e *planEnforcer) CheckCustomSlug(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	plan := e.plan(ctx, userID)
	if plan.ID == PlanFree {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			"custom slugs are not available on the Free plan; upgrade to Pro or Growth",
		)
	}
	return nil
}

func (e *planEnforcer) CheckWebhookLimit(ctx context.Context, userID uuid.UUID) *dto.ServiceError {
	plan := e.plan(ctx, userID)
	if !plan.HasWebhooks {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			fmt.Sprintf("webhooks are not available on the %s plan; upgrade to use webhooks", plan.Name),
		)
	}
	if plan.MaxWebhooks == -1 {
		return nil
	}
	count, err := e.webhookRepo.CountByUserID(ctx, userID)
	if err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count webhooks")
	}
	if count >= int64(plan.MaxWebhooks) {
		return dto.NewForbiddenError(
			constant.ErrCodePlanLimitReached,
			fmt.Sprintf("your %s plan allows up to %d webhooks; upgrade to create more", plan.Name, plan.MaxWebhooks),
		)
	}
	return nil
}
