package request

type UpdateBioPageRequest struct {
	Username    string `json:"username"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
	Theme       string `json:"theme"`
	IsPublished bool   `json:"is_published"`
}

type CreateBioLinkRequest struct {
	Title string `json:"title" binding:"required"`
	URL   string `json:"url" binding:"required,url"`
}

type UpdateBioLinkRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	IsActive *bool  `json:"is_active"`
}

type ReorderBioLinksRequest struct {
	IDs []string `json:"ids" binding:"required,min=1"`
}
