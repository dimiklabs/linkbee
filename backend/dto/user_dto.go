package dto

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role,omitempty"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type LoginOptions struct {
	RememberMe  bool
	LoginMethod string
	UserAgent   string
	IPAddress   string
}

type UserListRequest struct {
	Page   int    `json:"page" form:"page"`
	Limit  int    `json:"limit" form:"limit"`
	Search string `json:"search" form:"search"`
	Role   string `json:"role" form:"role"`
	Status string `json:"status" form:"status"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type BulkUpdateStatusRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`
	Status  string   `json:"status" binding:"required"`
}

type PaginatedResponse struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
}
