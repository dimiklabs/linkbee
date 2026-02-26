package user

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

type UserServiceI interface {
	// Create
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, *dto.ServiceError)

	// Read
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, *dto.ServiceError)
	GetUserByEmail(ctx context.Context, email string) (*model.User, *dto.ServiceError)
	GetUserByPhone(ctx context.Context, phone string) (*model.User, *dto.ServiceError)
	GetUserList(ctx context.Context, req *dto.UserListRequest) (*dto.PaginatedResponse, *dto.ServiceError)

	// Update
	UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*model.User, *dto.ServiceError)
	ChangePassword(ctx context.Context, id uuid.UUID, req *dto.ChangePasswordRequest) *dto.ServiceError
	SetPassword(ctx context.Context, id uuid.UUID, newPassword string) *dto.ServiceError
	UpdateRole(ctx context.Context, id uuid.UUID, req *dto.UpdateRoleRequest) *dto.ServiceError
	UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateStatusRequest) *dto.ServiceError
	UpdateLastLogin(ctx context.Context, id uuid.UUID) *dto.ServiceError
	ActivateUser(ctx context.Context, id uuid.UUID) *dto.ServiceError
	DeactivateUser(ctx context.Context, id uuid.UUID) *dto.ServiceError
	BulkUpdateStatus(ctx context.Context, req *dto.BulkUpdateStatusRequest) *dto.ServiceError

	// Delete
	DeleteUser(ctx context.Context, id uuid.UUID) *dto.ServiceError
	HardDeleteUser(ctx context.Context, id uuid.UUID) *dto.ServiceError
	RestoreUser(ctx context.Context, id uuid.UUID) *dto.ServiceError
	ScheduleDeletion(ctx context.Context, id uuid.UUID, days int) *dto.ServiceError
	CancelScheduledDeletion(ctx context.Context, id uuid.UUID) *dto.ServiceError

	// Checks
	IsEmailAvailable(ctx context.Context, email string, excludeID *uuid.UUID) (bool, *dto.ServiceError)

	// Count
	GetUserCount(ctx context.Context) (int64, *dto.ServiceError)
	GetUserCountByStatus(ctx context.Context, status string) (int64, *dto.ServiceError)
}

type UserService struct {
	userRepo repository.UserRepositoryI
}

func NewUserService(userRepo repository.UserRepositoryI) UserServiceI {
	return &UserService{
		userRepo: userRepo,
	}
}

// --------------- Create ---------------

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Creating new user",
		zap.String("email", req.Email))

	email := strings.ToLower(strings.TrimSpace(req.Email))

	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check email existence",
			zap.String("email", email),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	if exists {
		logger.WarnCtx(ctx, "Email already exists",
			zap.String("email", email))
		return nil, dto.NewServiceError(constant.ErrCodeEmailAlreadyExists, constant.ErrMsgEmailAlreadyExists, http.StatusConflict)
	}

	if req.Phone != "" {
		phoneExists, err := s.userRepo.ExistsByPhone(ctx, req.Phone)
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to check phone existence",
				zap.String("phone", req.Phone),
				zap.Error(err))
			return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
		}
		if phoneExists {
			logger.WarnCtx(ctx, "Phone already exists",
				zap.String("phone", req.Phone))
			return nil, dto.NewServiceError(constant.ErrCodePhoneAlreadyExists, constant.ErrMsgPhoneAlreadyExists, http.StatusConflict)
		}
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to hash password", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &model.User{
		Email:     email,
		Password:  hashedPassword,
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
		Phone:     strings.TrimSpace(req.Phone),
		Status:    "active",
		Role:      role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.ErrorCtx(ctx, "Failed to create user",
			zap.String("email", email),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User created successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	return user, nil
}

// --------------- Read ---------------

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Getting user by ID",
		zap.String("user_id", id.String()))

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WarnCtx(ctx, "User not found",
				zap.String("user_id", id.String()))
			return nil, dto.NewServiceError(constant.ErrCodeUserNotFound, constant.ErrMsgUserNotFound, http.StatusNotFound)
		}
		logger.ErrorCtx(ctx, "Failed to get user by ID",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Getting user by email",
		zap.String("email", email))

	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WarnCtx(ctx, "User not found by email",
				zap.String("email", email))
			return nil, dto.NewServiceError(constant.ErrCodeUserNotFound, constant.ErrMsgUserNotFound, http.StatusNotFound)
		}
		logger.ErrorCtx(ctx, "Failed to get user by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return user, nil
}

func (s *UserService) GetUserByPhone(ctx context.Context, phone string) (*model.User, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Getting user by phone",
		zap.String("phone", phone))

	user, err := s.userRepo.GetByPhone(ctx, strings.TrimSpace(phone))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WarnCtx(ctx, "User not found by phone",
				zap.String("phone", phone))
			return nil, dto.NewServiceError(constant.ErrCodeUserNotFound, constant.ErrMsgUserNotFound, http.StatusNotFound)
		}
		logger.ErrorCtx(ctx, "Failed to get user by phone",
			zap.String("phone", phone),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return user, nil
}

func (s *UserService) GetUserList(ctx context.Context, req *dto.UserListRequest) (*dto.PaginatedResponse, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Getting user list",
		zap.Int("page", req.Page),
		zap.Int("limit", req.Limit))

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}
	offset := (req.Page - 1) * req.Limit

	var users []model.User
	var total int64
	var err error

	switch {
	case req.Search != "":
		users, total, err = s.userRepo.SearchByName(ctx, req.Search, offset, req.Limit)
	case req.Status != "":
		users, total, err = s.userRepo.GetByStatus(ctx, req.Status, offset, req.Limit)
	case req.Role != "":
		users, total, err = s.userRepo.GetByRole(ctx, req.Role, offset, req.Limit)
	default:
		users, total, err = s.userRepo.GetAll(ctx, offset, req.Limit)
	}

	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get user list", zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	logger.DebugCtx(ctx, "User list fetched",
		zap.Int64("total", total),
		zap.Int("page", req.Page),
		zap.Int("total_pages", totalPages))

	return &dto.PaginatedResponse{
		Data:       users,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// --------------- Update ---------------

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*model.User, *dto.ServiceError) {
	logger.InfoCtx(ctx, "Updating user profile",
		zap.String("user_id", id.String()))

	user, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return nil, svcErr
	}

	user.FirstName = strings.TrimSpace(req.FirstName)
	user.LastName = strings.TrimSpace(req.LastName)
	user.Phone = strings.TrimSpace(req.Phone)

	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.ErrorCtx(ctx, "Failed to update user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return nil, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User profile updated",
		zap.String("user_id", id.String()))
	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, id uuid.UUID, req *dto.ChangePasswordRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Changing user password",
		zap.String("user_id", id.String()))

	user, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		logger.WarnCtx(ctx, "Invalid old password provided",
			zap.String("user_id", id.String()))
		return dto.NewServiceError(constant.ErrCodeInvalidPassword, constant.ErrMsgInvalidPassword, http.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword)); err == nil {
		logger.WarnCtx(ctx, "New password is same as old password",
			zap.String("user_id", id.String()))
		return dto.NewServiceError(constant.ErrCodeSamePassword, constant.ErrMsgSamePassword, http.StatusBadRequest)
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to hash new password", zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	if err := s.userRepo.UpdatePassword(ctx, id, hashedPassword); err != nil {
		logger.ErrorCtx(ctx, "Failed to update password in database",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User password changed successfully",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) SetPassword(ctx context.Context, id uuid.UUID, newPassword string) *dto.ServiceError {
	logger.InfoCtx(ctx, "Setting user password",
		zap.String("user_id", id.String()))

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to hash new password", zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	if err := s.userRepo.UpdatePassword(ctx, id, hashedPassword); err != nil {
		logger.ErrorCtx(ctx, "Failed to set password in database",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User password set successfully",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) UpdateRole(ctx context.Context, id uuid.UUID, req *dto.UpdateRoleRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Updating user role",
		zap.String("user_id", id.String()),
		zap.String("role", req.Role))

	_, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	if err := s.userRepo.UpdateRole(ctx, id, req.Role); err != nil {
		logger.ErrorCtx(ctx, "Failed to update user role",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User role updated",
		zap.String("user_id", id.String()),
		zap.String("role", req.Role))
	return nil
}

func (s *UserService) UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateStatusRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Updating user status",
		zap.String("user_id", id.String()),
		zap.String("status", req.Status))

	_, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	if err := s.userRepo.UpdateStatus(ctx, id, req.Status); err != nil {
		logger.ErrorCtx(ctx, "Failed to update user status",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User status updated",
		zap.String("user_id", id.String()),
		zap.String("status", req.Status))
	return nil
}

func (s *UserService) UpdateLastLogin(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.DebugCtx(ctx, "Updating user last login",
		zap.String("user_id", id.String()))

	if err := s.userRepo.UpdateLastLogin(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to update last login",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return nil
}

func (s *UserService) ActivateUser(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Activating user",
		zap.String("user_id", id.String()))

	_, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	if err := s.userRepo.SetActive(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to activate user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User activated",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) DeactivateUser(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Deactivating user",
		zap.String("user_id", id.String()))

	_, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	if err := s.userRepo.SetInactive(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to deactivate user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User deactivated",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) BulkUpdateStatus(ctx context.Context, req *dto.BulkUpdateStatusRequest) *dto.ServiceError {
	logger.InfoCtx(ctx, "Bulk updating user status",
		zap.Int("count", len(req.UserIDs)),
		zap.String("status", req.Status))

	ids := make([]uuid.UUID, 0, len(req.UserIDs))
	for _, idStr := range req.UserIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			logger.WarnCtx(ctx, "Invalid user ID in bulk update",
				zap.String("invalid_id", idStr),
				zap.Error(err))
			return dto.NewServiceError(constant.ErrCodeBadRequest, "Invalid user ID: "+idStr, http.StatusBadRequest)
		}
		ids = append(ids, id)
	}

	if err := s.userRepo.BulkUpdateStatus(ctx, ids, req.Status); err != nil {
		logger.ErrorCtx(ctx, "Failed to bulk update user status",
			zap.Int("count", len(ids)),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "Bulk user status updated",
		zap.Int("count", len(ids)),
		zap.String("status", req.Status))
	return nil
}

// --------------- Delete ---------------

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Deleting user",
		zap.String("user_id", id.String()))

	_, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User deleted",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) HardDeleteUser(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.WarnCtx(ctx, "Hard deleting user",
		zap.String("user_id", id.String()))

	if err := s.userRepo.HardDelete(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to hard delete user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.WarnCtx(ctx, "User hard deleted",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) RestoreUser(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Restoring user",
		zap.String("user_id", id.String()))

	if err := s.userRepo.Restore(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to restore user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User restored",
		zap.String("user_id", id.String()))
	return nil
}

func (s *UserService) ScheduleDeletion(ctx context.Context, id uuid.UUID, days int) *dto.ServiceError {
	logger.InfoCtx(ctx, "Scheduling user deletion",
		zap.String("user_id", id.String()),
		zap.Int("days", days))

	_, svcErr := s.GetUserByID(ctx, id)
	if svcErr != nil {
		return svcErr
	}

	deletionDate := time.Now().AddDate(0, 0, days)

	if err := s.userRepo.ScheduleDeletion(ctx, id, deletionDate); err != nil {
		logger.ErrorCtx(ctx, "Failed to schedule user deletion",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "User deletion scheduled",
		zap.String("user_id", id.String()),
		zap.Time("deletion_date", deletionDate))
	return nil
}

func (s *UserService) CancelScheduledDeletion(ctx context.Context, id uuid.UUID) *dto.ServiceError {
	logger.InfoCtx(ctx, "Cancelling scheduled deletion",
		zap.String("user_id", id.String()))

	if err := s.userRepo.CancelScheduledDeletion(ctx, id); err != nil {
		logger.ErrorCtx(ctx, "Failed to cancel scheduled deletion",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}

	logger.InfoCtx(ctx, "Scheduled deletion cancelled, user reactivated",
		zap.String("user_id", id.String()))
	return nil
}

// --------------- Checks ---------------

func (s *UserService) IsEmailAvailable(ctx context.Context, email string, excludeID *uuid.UUID) (bool, *dto.ServiceError) {
	logger.DebugCtx(ctx, "Checking email availability",
		zap.String("email", email))

	available, err := s.userRepo.IsEmailAvailable(ctx, strings.ToLower(strings.TrimSpace(email)), excludeID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to check email availability",
			zap.String("email", email),
			zap.Error(err))
		return false, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return available, nil
}

// --------------- Count ---------------

func (s *UserService) GetUserCount(ctx context.Context) (int64, *dto.ServiceError) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to count users", zap.Error(err))
		return 0, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return count, nil
}

func (s *UserService) GetUserCountByStatus(ctx context.Context, status string) (int64, *dto.ServiceError) {
	count, err := s.userRepo.CountByStatus(ctx, status)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to count users by status",
			zap.String("status", status),
			zap.Error(err))
		return 0, dto.NewServiceError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer, http.StatusInternalServerError)
	}
	return count, nil
}

// --------------- Helpers ---------------

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
