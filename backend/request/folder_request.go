package request

// CreateFolderRequest is the request body for creating a link folder.
type CreateFolderRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Color string `json:"color" binding:"omitempty,max=20"`
}

// UpdateFolderRequest is the request body for renaming / recoloring a folder.
type UpdateFolderRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=100"`
	Color string `json:"color" binding:"omitempty,max=20"`
}
