package admin

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

type StatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
	TotalLinks    int64 `json:"total_links"`
}

type UsersListResponse struct {
	Users []UserItem `json:"users"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

type UserItem struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	FirstName     string  `json:"first_name,omitempty"`
	LastName      string  `json:"last_name,omitempty"`
	Status        string  `json:"status"`
	Role          string  `json:"role"`
	AuthProvider  string  `json:"auth_provider"`
	EmailVerified bool    `json:"email_verified"`
	CreatedAt     string  `json:"created_at"`
	LastLogin     *string `json:"last_login,omitempty"`
}

type AdminServiceI interface {
	GetStats(ctx context.Context) (*StatsResponse, *dto.ServiceError)
	ListUsers(ctx context.Context, search string, page, limit int) (*UsersListResponse, *dto.ServiceError)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) *dto.ServiceError
}

type adminService struct {
	userRepo repository.UserRepositoryI
	linkRepo repository.LinkRepositoryI
}

func NewAdminService(userRepo repository.UserRepositoryI, linkRepo repository.LinkRepositoryI) AdminServiceI {
	return &adminService{userRepo: userRepo, linkRepo: linkRepo}
}

func (s *adminService) GetStats(ctx context.Context) (*StatsResponse, *dto.ServiceError) {
	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count users")
	}
	active, err := s.userRepo.CountByStatus(ctx, "active")
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count active users")
	}
	inactive, err := s.userRepo.CountByStatus(ctx, "inactive")
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to count inactive users")
	}

	totalLinks, err := s.linkRepo.Count(ctx)
	if err != nil {
		totalLinks = 0
	}

	return &StatsResponse{
		TotalUsers:    total,
		ActiveUsers:   active,
		InactiveUsers: inactive,
		TotalLinks:    totalLinks,
	}, nil
}

func (s *adminService) ListUsers(ctx context.Context, search string, page, limit int) (*UsersListResponse, *dto.ServiceError) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var users []model.User
	var total int64
	var err error

	search = strings.TrimSpace(search)
	if search != "" {
		users, total, err = s.userRepo.SearchAll(ctx, search, offset, limit)
	} else {
		users, total, err = s.userRepo.GetAll(ctx, offset, limit)
	}
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to list users")
	}

	items := make([]UserItem, 0, len(users))
	for _, u := range users {
		item := UserItem{
			ID:            u.ID.String(),
			Email:         u.Email,
			FirstName:     u.FirstName,
			LastName:      u.LastName,
			Status:        u.Status,
			Role:          u.Role,
			AuthProvider:  u.AuthProvider,
			EmailVerified: u.EmailVerified,
			CreatedAt:     u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		}
		if u.LastLogin != nil {
			s := u.LastLogin.UTC().Format("2006-01-02T15:04:05Z")
			item.LastLogin = &s
		}
		items = append(items, item)
	}

	return &UsersListResponse{
		Users: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *adminService) UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) *dto.ServiceError {
	if status != "active" && status != "inactive" && status != "banned" {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "status must be active, inactive, or banned")
	}
	if err := s.userRepo.UpdateStatus(ctx, userID, status); err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to update user status")
	}
	return nil
}
