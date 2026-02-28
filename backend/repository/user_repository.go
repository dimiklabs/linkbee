package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
)

type UserRepositoryI interface {
	// Create
	Create(ctx context.Context, user *model.User) error

	// Read - single
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	GetByGitHubID(ctx context.Context, githubID string) (*model.User, error)

	// Read - list
	GetAll(ctx context.Context, offset, limit int) ([]model.User, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]model.User, int64, error)
	GetByRole(ctx context.Context, role string, offset, limit int) ([]model.User, int64, error)
	GetActiveUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error)
	GetInactiveUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error)
	SearchByName(ctx context.Context, name string, offset, limit int) ([]model.User, int64, error)
	SearchAll(ctx context.Context, query string, offset, limit int) ([]model.User, int64, error)

	// Update
	Update(ctx context.Context, user *model.User) error
	UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error
	UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateRole(ctx context.Context, id uuid.UUID, role string) error
	UpdateProfile(ctx context.Context, id uuid.UUID, firstName, lastName, phone string) error
	SetActive(ctx context.Context, id uuid.UUID) error
	SetInactive(ctx context.Context, id uuid.UUID) error
	BulkUpdateStatus(ctx context.Context, ids []uuid.UUID, status string) error
	ScheduleDeletion(ctx context.Context, id uuid.UUID, deletionDate time.Time) error
	CancelScheduledDeletion(ctx context.Context, id uuid.UUID) error

	// Delete
	Delete(ctx context.Context, id uuid.UUID) error
	HardDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error

	// Checks
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	IsEmailAvailable(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)

	// Count
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)

	// Time series
	GetRegistrationTimeSeries(ctx context.Context, from, to time.Time) ([]TimeSeriesPoint, error)
}

type UserRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewUserRepository(masterDB, replicaDB *gorm.DB) UserRepositoryI {
	return &UserRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

// --------------- Create ---------------

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	logger.DebugCtx(ctx, "Creating user",
		zap.String("email", user.Email))

	if err := r.masterDB.WithContext(ctx).Create(user).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create user",
			zap.String("email", user.Email),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User created successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))
	return nil
}

// --------------- Read - single ---------------

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	logger.DebugCtx(ctx, "Fetching user by ID",
		zap.String("user_id", id.String()))

	var user model.User
	if err := r.replicaDB.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch user by ID",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	logger.DebugCtx(ctx, "Fetching user by email",
		zap.String("email", email))

	var user model.User
	if err := r.replicaDB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch user by email",
			zap.String("email", email),
			zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	logger.DebugCtx(ctx, "Fetching user by phone",
		zap.String("phone", phone))

	var user model.User
	if err := r.replicaDB.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch user by phone",
			zap.String("phone", phone),
			zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	logger.DebugCtx(ctx, "Fetching user by Google ID")

	var user model.User
	if err := r.replicaDB.WithContext(ctx).Where("google_id = ?", googleID).First(&user).Error; err != nil {
		logger.DebugCtx(ctx, "User not found by Google ID",
			zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByGitHubID(ctx context.Context, githubID string) (*model.User, error) {
	logger.DebugCtx(ctx, "Fetching user by GitHub ID")

	var user model.User
	if err := r.replicaDB.WithContext(ctx).Where("github_id = ?", githubID).First(&user).Error; err != nil {
		logger.DebugCtx(ctx, "User not found by GitHub ID",
			zap.Error(err))
		return nil, err
	}
	return &user, nil
}


// --------------- Read - list ---------------

func (r *UserRepository) GetAll(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	logger.DebugCtx(ctx, "Fetching all users",
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	var users []model.User
	var total int64

	if err := r.replicaDB.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count users", zap.Error(err))
		return nil, 0, err
	}

	if err := r.replicaDB.WithContext(ctx).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch users", zap.Error(err))
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) GetByStatus(ctx context.Context, status string, offset, limit int) ([]model.User, int64, error) {
	logger.DebugCtx(ctx, "Fetching users by status",
		zap.String("status", status),
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	var users []model.User
	var total int64

	query := r.replicaDB.WithContext(ctx).Model(&model.User{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count users by status", zap.String("status", status), zap.Error(err))
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch users by status", zap.String("status", status), zap.Error(err))
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) GetByRole(ctx context.Context, role string, offset, limit int) ([]model.User, int64, error) {
	logger.DebugCtx(ctx, "Fetching users by role",
		zap.String("role", role),
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	var users []model.User
	var total int64

	query := r.replicaDB.WithContext(ctx).Model(&model.User{}).Where("role = ?", role)

	if err := query.Count(&total).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count users by role", zap.String("role", role), zap.Error(err))
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch users by role", zap.String("role", role), zap.Error(err))
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	return r.GetByStatus(ctx, "active", offset, limit)
}

func (r *UserRepository) GetInactiveUsers(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	return r.GetByStatus(ctx, "inactive", offset, limit)
}

func (r *UserRepository) SearchByName(ctx context.Context, name string, offset, limit int) ([]model.User, int64, error) {
	logger.DebugCtx(ctx, "Searching users by name",
		zap.String("name", name),
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	var users []model.User
	var total int64

	searchPattern := "%" + name + "%"
	query := r.replicaDB.WithContext(ctx).Model(&model.User{}).
		Where("first_name ILIKE ? OR last_name ILIKE ?", searchPattern, searchPattern)

	if err := query.Count(&total).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count users by name search", zap.String("name", name), zap.Error(err))
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to search users by name", zap.String("name", name), zap.Error(err))
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) SearchAll(ctx context.Context, query string, offset, limit int) ([]model.User, int64, error) {
	logger.DebugCtx(ctx, "Searching users",
		zap.String("query", query),
		zap.Int("offset", offset),
		zap.Int("limit", limit))

	var users []model.User
	var total int64

	pattern := "%" + query + "%"
	q := r.replicaDB.WithContext(ctx).Model(&model.User{}).
		Where("email ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?", pattern, pattern, pattern)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// --------------- Update ---------------

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	logger.DebugCtx(ctx, "Updating user",
		zap.String("user_id", user.ID.String()))

	if err := r.masterDB.WithContext(ctx).Save(user).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User updated successfully",
		zap.String("user_id", user.ID.String()))
	return nil
}

func (r *UserRepository) UpdateFields(ctx context.Context, id uuid.UUID, fields map[string]interface{}) error {
	logger.DebugCtx(ctx, "Updating user fields",
		zap.String("user_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(fields).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user fields",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User fields updated successfully",
		zap.String("user_id", id.String()))
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error {
	logger.DebugCtx(ctx, "Updating user password",
		zap.String("user_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user password",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User password updated successfully",
		zap.String("user_id", id.String()))
	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Updating user last login",
		zap.String("user_id", id.String()))

	now := time.Now()
	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login", now).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user last login",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}
	return nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	logger.DebugCtx(ctx, "Updating user status",
		zap.String("user_id", id.String()),
		zap.String("status", status))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user status",
			zap.String("user_id", id.String()),
			zap.String("status", status),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User status updated",
		zap.String("user_id", id.String()),
		zap.String("status", status))
	return nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, id uuid.UUID, role string) error {
	logger.DebugCtx(ctx, "Updating user role",
		zap.String("user_id", id.String()),
		zap.String("role", role))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("role", role).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user role",
			zap.String("user_id", id.String()),
			zap.String("role", role),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User role updated",
		zap.String("user_id", id.String()),
		zap.String("role", role))
	return nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id uuid.UUID, firstName, lastName, phone string) error {
	logger.DebugCtx(ctx, "Updating user profile",
		zap.String("user_id", id.String()))

	updates := map[string]any{
		"first_name": firstName,
		"last_name":  lastName,
		"phone":      phone,
	}

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update user profile",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User profile updated",
		zap.String("user_id", id.String()))
	return nil
}

func (r *UserRepository) SetActive(ctx context.Context, id uuid.UUID) error {
	return r.UpdateStatus(ctx, id, "active")
}

func (r *UserRepository) SetInactive(ctx context.Context, id uuid.UUID) error {
	return r.UpdateStatus(ctx, id, "inactive")
}

func (r *UserRepository) BulkUpdateStatus(ctx context.Context, ids []uuid.UUID, status string) error {
	logger.DebugCtx(ctx, "Bulk updating user status",
		zap.Int("count", len(ids)),
		zap.String("status", status))

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id IN ?", ids).
		Update("status", status).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to bulk update user status",
			zap.Int("count", len(ids)),
			zap.String("status", status),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Bulk user status updated",
		zap.Int("count", len(ids)),
		zap.String("status", status))
	return nil
}

func (r *UserRepository) ScheduleDeletion(ctx context.Context, id uuid.UUID, deletionDate time.Time) error {
	logger.InfoCtx(ctx, "Scheduling user deletion",
		zap.String("user_id", id.String()),
		zap.Time("deletion_date", deletionDate))

	updates := map[string]any{
		"status":                "inactive",
		"scheduled_deletion_at": deletionDate,
	}

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to schedule user deletion",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User deletion scheduled",
		zap.String("user_id", id.String()),
		zap.Time("deletion_date", deletionDate))
	return nil
}

func (r *UserRepository) CancelScheduledDeletion(ctx context.Context, id uuid.UUID) error {
	logger.InfoCtx(ctx, "Cancelling scheduled deletion",
		zap.String("user_id", id.String()))

	updates := map[string]any{
		"status":                "active",
		"scheduled_deletion_at": nil,
	}

	if err := r.masterDB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to cancel scheduled deletion",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Scheduled deletion cancelled, user reactivated",
		zap.String("user_id", id.String()))
	return nil
}

// --------------- Delete ---------------

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Soft deleting user",
		zap.String("user_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.User{}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to soft delete user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User soft deleted",
		zap.String("user_id", id.String()))
	return nil
}

func (r *UserRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	logger.WarnCtx(ctx, "Hard deleting user",
		zap.String("user_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Unscoped().
		Where("id = ?", id).
		Delete(&model.User{}).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to hard delete user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User hard deleted",
		zap.String("user_id", id.String()))
	return nil
}

func (r *UserRepository) Restore(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Restoring soft-deleted user",
		zap.String("user_id", id.String()))

	if err := r.masterDB.WithContext(ctx).
		Unscoped().
		Model(&model.User{}).
		Where("id = ?", id).
		Update("deleted_at", nil).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to restore user",
			zap.String("user_id", id.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "User restored",
		zap.String("user_id", id.String()))
	return nil
}

// --------------- Checks ---------------

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	logger.DebugCtx(ctx, "Checking if email exists",
		zap.String("email", email))

	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to check email existence",
			zap.String("email", email),
			zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	logger.DebugCtx(ctx, "Checking if phone exists",
		zap.String("phone", phone))

	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.User{}).
		Where("phone = ?", phone).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to check phone existence",
			zap.String("phone", phone),
			zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) IsEmailAvailable(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	logger.DebugCtx(ctx, "Checking email availability",
		zap.String("email", email))

	var count int64
	query := r.replicaDB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to check email availability",
			zap.String("email", email),
			zap.Error(err))
		return false, err
	}
	return count == 0, nil
}

// --------------- Count ---------------

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	logger.DebugCtx(ctx, "Counting total users")

	var count int64
	if err := r.replicaDB.WithContext(ctx).Model(&model.User{}).Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count users", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	logger.DebugCtx(ctx, "Counting users by status",
		zap.String("status", status))

	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.User{}).
		Where("status = ?", status).
		Count(&count).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to count users by status",
			zap.String("status", status),
			zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetRegistrationTimeSeries(ctx context.Context, from, to time.Time) ([]TimeSeriesPoint, error) {
	type result struct {
		Timestamp time.Time
		Count     int64
	}

	var rows []result
	if err := r.replicaDB.WithContext(ctx).
		Raw(`SELECT date_trunc('day', created_at) AS timestamp, COUNT(*) AS count
		     FROM users
		     WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL
		     GROUP BY 1
		     ORDER BY 1`, from, to).
		Scan(&rows).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to get registration time series", zap.Error(err))
		return nil, err
	}

	points := make([]TimeSeriesPoint, len(rows))
	for i, row := range rows {
		points[i] = TimeSeriesPoint{Timestamp: row.Timestamp, Count: row.Count}
	}
	return points, nil
}
