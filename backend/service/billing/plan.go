package billing

const (
	PlanFree     = "free"
	PlanPro      = "pro"
	PlanBusiness = "business"
)

// PlanInfo describes a subscription plan and its feature limits.
// A limit value of -1 means unlimited.
type PlanInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MaxLinks    int    `json:"max_links"`
	MaxAPIKeys  int    `json:"max_api_keys"`
	MaxWebhooks int    `json:"max_webhooks"`
	HasWebhooks bool   `json:"has_webhooks"`
}

// PlanCatalog is the single source of truth for plan limits.
var PlanCatalog = map[string]PlanInfo{
	PlanFree: {
		ID:          PlanFree,
		Name:        "Free",
		MaxLinks:    5,
		MaxAPIKeys:  1,
		MaxWebhooks: 0,
		HasWebhooks: false,
	},
	PlanPro: {
		ID:          PlanPro,
		Name:        "Pro",
		MaxLinks:    100,
		MaxAPIKeys:  5,
		MaxWebhooks: 5,
		HasWebhooks: true,
	},
	PlanBusiness: {
		ID:          PlanBusiness,
		Name:        "Business",
		MaxLinks:    -1,
		MaxAPIKeys:  -1,
		MaxWebhooks: -1,
		HasWebhooks: true,
	},
}

// GetPlan returns the PlanInfo for the given planID, defaulting to Free.
func GetPlan(planID string) PlanInfo {
	if p, ok := PlanCatalog[planID]; ok {
		return p
	}
	return PlanCatalog[PlanFree]
}
