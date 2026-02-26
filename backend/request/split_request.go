package request

// CreateVariantRequest is the body for adding a new A/B variant to a link.
type CreateVariantRequest struct {
	Name           string `json:"name" binding:"required,min=1,max=100"`
	DestinationURL string `json:"destination_url" binding:"required,url"`
	Weight         int    `json:"weight" binding:"required,min=1,max=1000"`
}

// UpdateVariantRequest is the body for editing an existing variant.
type UpdateVariantRequest struct {
	Name           string `json:"name" binding:"omitempty,min=1,max=100"`
	DestinationURL string `json:"destination_url" binding:"omitempty,url"`
	Weight         int    `json:"weight" binding:"omitempty,min=1,max=1000"`
}

// ToggleSplitTestRequest enables or disables split testing for a link.
type ToggleSplitTestRequest struct {
	Enabled bool `json:"enabled"`
}
