package request

// CreateReportRequest is the body for creating a new scheduled report.
type CreateReportRequest struct {
	Name      string   `json:"name"      binding:"required,min=1,max=255"`
	LinkIDs   []string `json:"link_ids"  binding:"required,min=1,max=10"`
	Frequency string   `json:"frequency" binding:"required,oneof=daily weekly monthly"`
}

// UpdateReportRequest is the body for updating an existing report.
type UpdateReportRequest struct {
	Name      string   `json:"name"      binding:"omitempty,min=1,max=255"`
	LinkIDs   []string `json:"link_ids"  binding:"omitempty,min=1,max=10"`
	Frequency string   `json:"frequency" binding:"omitempty,oneof=daily weekly monthly"`
	IsActive  *bool    `json:"is_active"`
}
