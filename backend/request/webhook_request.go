package request

type CreateWebhookRequest struct {
	URL    string   `json:"url" binding:"required,url"`
	Events []string `json:"events" binding:"required,min=1"`
}

type UpdateWebhookRequest struct {
	URL      string   `json:"url"`
	Events   []string `json:"events"`
	IsActive *bool    `json:"is_active"`
}
