package request

type CreatePixelRequest struct {
	PixelType    string `json:"pixel_type" binding:"required"`
	PixelID      string `json:"pixel_id"`
	CustomScript string `json:"custom_script"`
}

type TogglePixelTrackingRequest struct {
	Enabled bool `json:"enabled"`
}
