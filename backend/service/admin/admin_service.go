package admin

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
	util "github.com/shafikshaon/shortlink/util"
)

type StatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
	TotalLinks    int64 `json:"total_links"`
	TotalClicks   int64 `json:"total_clicks"`
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

type GrowthTimeSeriesResponse struct {
	Users []TimeSeriesData `json:"users"`
	Links []TimeSeriesData `json:"links"`
}

type TimeSeriesData struct {
	Timestamp string `json:"timestamp"`
	Count     int64  `json:"count"`
}

type ImpersonationResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TargetUserID string `json:"target_user_id"`
	TargetEmail  string `json:"target_email"`
}

type AdminServiceI interface {
	GetStats(ctx context.Context) (*StatsResponse, *dto.ServiceError)
	ListUsers(ctx context.Context, search string, page, limit int) (*UsersListResponse, *dto.ServiceError)
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status string) *dto.ServiceError
	GetGrowthTimeSeries(ctx context.Context) (*GrowthTimeSeriesResponse, error)
	UpdateUserRole(ctx context.Context, adminID uuid.UUID, targetUserID uuid.UUID, role string) *dto.ServiceError
	ImpersonateUser(ctx context.Context, adminID uuid.UUID, targetUserID uuid.UUID, cfg *config.AppConfig) (*ImpersonationResponse, *dto.ServiceError)
}

type adminService struct {
	userRepo       repository.UserRepositoryI
	linkRepo       repository.LinkRepositoryI
	clickEventRepo repository.ClickEventRepositoryI
}

func NewAdminService(userRepo repository.UserRepositoryI, linkRepo repository.LinkRepositoryI, clickEventRepo repository.ClickEventRepositoryI) AdminServiceI {
	return &adminService{userRepo: userRepo, linkRepo: linkRepo, clickEventRepo: clickEventRepo}
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

	totalClicks, err := s.clickEventRepo.GetTotalClicks(ctx)
	if err != nil {
		totalClicks = 0
	}

	return &StatsResponse{
		TotalUsers:    total,
		ActiveUsers:   active,
		InactiveUsers: inactive,
		TotalLinks:    totalLinks,
		TotalClicks:   totalClicks,
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

func (s *adminService) UpdateUserRole(ctx context.Context, adminID uuid.UUID, targetUserID uuid.UUID, role string) *dto.ServiceError {
	if role != "admin" && role != "user" {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "role must be admin or user")
	}
	if adminID == targetUserID {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "cannot change your own role")
	}
	if err := s.userRepo.UpdateRole(ctx, targetUserID, role); err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to update user role")
	}
	return nil
}

func (s *adminService) ImpersonateUser(ctx context.Context, adminID uuid.UUID, targetUserID uuid.UUID, cfg *config.AppConfig) (*ImpersonationResponse, *dto.ServiceError) {
	if adminID == targetUserID {
		return nil, dto.NewBadRequestError(constant.ErrCodeBadRequest, "cannot impersonate yourself")
	}
	targetUser, err := s.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "user not found")
	}
	jwtCfg := util.JWTConfig{
		Secret:              cfg.JWTSecret,
		Issuer:              cfg.JWTIssuer,
		AccessExpiryMinutes: cfg.JWTAccessExpiry,
		RefreshExpiryDays:   cfg.JWTRefreshExpiry,
	}
	tokens, genErr := util.GenerateTokenPair(&jwtCfg, targetUser.ID.String(), targetUser.Email, targetUser.Role)
	if genErr != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, "failed to generate impersonation token")
	}
	return &ImpersonationResponse{
		AccessToken:  tokens.AccessToken,
		ExpiresIn:    tokens.ExpiresIn,
		TargetUserID: targetUser.ID.String(),
		TargetEmail:  targetUser.Email,
	}, nil
}

func (s *adminService) GetGrowthTimeSeries(ctx context.Context) (*GrowthTimeSeriesResponse, error) {
	to := time.Now()
	from := to.AddDate(0, 0, -30)

	userPoints, err := s.userRepo.GetRegistrationTimeSeries(ctx, from, to)
	if err != nil {
		return nil, err
	}

	linkPoints, err := s.linkRepo.GetCreationTimeSeries(ctx, from, to)
	if err != nil {
		return nil, err
	}

	userSeries := make([]TimeSeriesData, len(userPoints))
	for i, p := range userPoints {
		userSeries[i] = TimeSeriesData{
			Timestamp: p.Timestamp.UTC().Format("2006-01-02"),
			Count:     p.Count,
		}
	}

	linkSeries := make([]TimeSeriesData, len(linkPoints))
	for i, p := range linkPoints {
		linkSeries[i] = TimeSeriesData{
			Timestamp: p.Timestamp.UTC().Format("2006-01-02"),
			Count:     p.Count,
		}
	}

	return &GrowthTimeSeriesResponse{
		Users: userSeries,
		Links: linkSeries,
	}, nil
}
