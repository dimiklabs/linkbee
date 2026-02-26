package request

// CreateGeoRuleRequest is the body for adding a new geo routing rule.
type CreateGeoRuleRequest struct {
	CountryCode    string `json:"country_code" binding:"required,len=2"`
	DestinationURL string `json:"destination_url" binding:"required,url"`
	Priority       int    `json:"priority" binding:"min=0"`
}

// UpdateGeoRuleRequest is the body for editing an existing geo routing rule.
type UpdateGeoRuleRequest struct {
	CountryCode    string `json:"country_code" binding:"omitempty,len=2"`
	DestinationURL string `json:"destination_url" binding:"omitempty,url"`
	Priority       int    `json:"priority" binding:"min=0"`
}

// ToggleGeoRoutingRequest enables or disables geo routing for a link.
type ToggleGeoRoutingRequest struct {
	Enabled bool `json:"enabled"`
}
