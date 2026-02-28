package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
)

type TeamRepositoryI interface {
	// Team CRUD
	CreateTeam(ctx context.Context, team *model.Team) error
	GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error)
	GetTeamBySlug(ctx context.Context, slug string) (*model.Team, error)
	UpdateTeam(ctx context.Context, team *model.Team) error
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	ListTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error)

	// Member management
	AddMember(ctx context.Context, member *model.TeamMember) error
	GetMember(ctx context.Context, teamID, userID uuid.UUID) (*model.TeamMember, error)
	GetMemberByInviteToken(ctx context.Context, token string) (*model.TeamMember, error)
	UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role string) error
	UpdateMemberStatus(ctx context.Context, id uuid.UUID, status string, joinedAt *time.Time) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error)
	IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error)
	IsOwner(ctx context.Context, teamID, userID uuid.UUID) (bool, error)
	CountMembers(ctx context.Context, teamID uuid.UUID) (int64, error)
}

type TeamRepository struct {
	masterDB  *gorm.DB
	replicaDB *gorm.DB
}

func NewTeamRepository(masterDB, replicaDB *gorm.DB) TeamRepositoryI {
	return &TeamRepository{
		masterDB:  masterDB,
		replicaDB: replicaDB,
	}
}

// --------------- Team CRUD ---------------

func (r *TeamRepository) CreateTeam(ctx context.Context, team *model.Team) error {
	logger.DebugCtx(ctx, "Creating team", zap.String("name", team.Name))

	if err := r.masterDB.WithContext(ctx).Create(team).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to create team", zap.String("name", team.Name), zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Team created", zap.String("team_id", team.ID.String()), zap.String("name", team.Name))
	return nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id uuid.UUID) (*model.Team, error) {
	logger.DebugCtx(ctx, "Fetching team by ID", zap.String("team_id", id.String()))

	var team model.Team
	if err := r.replicaDB.WithContext(ctx).Where("id = ?", id).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) GetTeamBySlug(ctx context.Context, slug string) (*model.Team, error) {
	logger.DebugCtx(ctx, "Fetching team by slug", zap.String("slug", slug))

	var team model.Team
	if err := r.replicaDB.WithContext(ctx).Where("slug = ?", slug).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, team *model.Team) error {
	logger.DebugCtx(ctx, "Updating team", zap.String("team_id", team.ID.String()))

	team.UpdatedAt = time.Now()
	if err := r.masterDB.WithContext(ctx).Save(team).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update team", zap.String("team_id", team.ID.String()), zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Team updated", zap.String("team_id", team.ID.String()))
	return nil
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	logger.DebugCtx(ctx, "Deleting team", zap.String("team_id", id.String()))

	result := r.masterDB.WithContext(ctx).Where("id = ?", id).Delete(&model.Team{})
	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to delete team", zap.String("team_id", id.String()), zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.InfoCtx(ctx, "Team deleted", zap.String("team_id", id.String()))
	return nil
}

func (r *TeamRepository) ListTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Team, error) {
	logger.DebugCtx(ctx, "Listing teams for user", zap.String("user_id", userID.String()))

	// Get teams where user is owner
	var ownedTeams []model.Team
	if err := r.replicaDB.WithContext(ctx).Where("owner_id = ?", userID).Find(&ownedTeams).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch owned teams", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}

	// Get team IDs where user is an active member (not owner)
	var memberTeamIDs []uuid.UUID
	if err := r.replicaDB.WithContext(ctx).Model(&model.TeamMember{}).
		Select("team_id").
		Where("user_id = ? AND status = 'active' AND deleted_at IS NULL", userID).
		Pluck("team_id", &memberTeamIDs).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch member team IDs", zap.String("user_id", userID.String()), zap.Error(err))
		return nil, err
	}

	// Collect owned team IDs
	ownedIDs := make(map[uuid.UUID]bool)
	for _, t := range ownedTeams {
		ownedIDs[t.ID] = true
	}

	// Filter out teams where user is already included as owner
	var remainingIDs []uuid.UUID
	for _, id := range memberTeamIDs {
		if !ownedIDs[id] {
			remainingIDs = append(remainingIDs, id)
		}
	}

	result := ownedTeams
	if len(remainingIDs) > 0 {
		var memberTeams []model.Team
		if err := r.replicaDB.WithContext(ctx).Where("id IN ?", remainingIDs).Find(&memberTeams).Error; err != nil {
			logger.ErrorCtx(ctx, "Failed to fetch member teams", zap.String("user_id", userID.String()), zap.Error(err))
			return nil, err
		}
		result = append(result, memberTeams...)
	}

	return result, nil
}

// --------------- Member management ---------------

func (r *TeamRepository) AddMember(ctx context.Context, member *model.TeamMember) error {
	logger.DebugCtx(ctx, "Adding team member",
		zap.String("team_id", member.TeamID.String()),
		zap.String("invite_email", member.InviteEmail))

	if err := r.masterDB.WithContext(ctx).Create(member).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to add team member",
			zap.String("team_id", member.TeamID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Team member added", zap.String("member_id", member.ID.String()))
	return nil
}

func (r *TeamRepository) GetMember(ctx context.Context, teamID, userID uuid.UUID) (*model.TeamMember, error) {
	logger.DebugCtx(ctx, "Fetching team member",
		zap.String("team_id", teamID.String()),
		zap.String("user_id", userID.String()))

	var member model.TeamMember
	if err := r.replicaDB.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamRepository) GetMemberByInviteToken(ctx context.Context, token string) (*model.TeamMember, error) {
	logger.DebugCtx(ctx, "Fetching team member by invite token")

	var member model.TeamMember
	if err := r.replicaDB.WithContext(ctx).
		Where("invite_token = ?", token).
		First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamRepository) UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role string) error {
	logger.DebugCtx(ctx, "Updating member role",
		zap.String("team_id", teamID.String()),
		zap.String("user_id", userID.String()),
		zap.String("role", role))

	result := r.masterDB.WithContext(ctx).
		Model(&model.TeamMember{}).
		Where("team_id = ? AND user_id = ? AND deleted_at IS NULL", teamID, userID).
		Updates(map[string]interface{}{
			"role":       role,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to update member role",
			zap.String("team_id", teamID.String()),
			zap.String("user_id", userID.String()),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *TeamRepository) UpdateMemberStatus(ctx context.Context, id uuid.UUID, status string, joinedAt *time.Time) error {
	logger.DebugCtx(ctx, "Updating member status",
		zap.String("member_id", id.String()),
		zap.String("status", status))

	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if joinedAt != nil {
		updates["joined_at"] = joinedAt
	}

	if err := r.masterDB.WithContext(ctx).
		Model(&model.TeamMember{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to update member status",
			zap.String("member_id", id.String()),
			zap.Error(err))
		return err
	}

	return nil
}

func (r *TeamRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	logger.DebugCtx(ctx, "Removing team member",
		zap.String("team_id", teamID.String()),
		zap.String("user_id", userID.String()))

	result := r.masterDB.WithContext(ctx).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&model.TeamMember{})

	if result.Error != nil {
		logger.ErrorCtx(ctx, "Failed to remove team member",
			zap.String("team_id", teamID.String()),
			zap.String("user_id", userID.String()),
			zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.InfoCtx(ctx, "Team member removed",
		zap.String("team_id", teamID.String()),
		zap.String("user_id", userID.String()))
	return nil
}

func (r *TeamRepository) ListMembers(ctx context.Context, teamID uuid.UUID) ([]model.TeamMember, error) {
	logger.DebugCtx(ctx, "Listing team members", zap.String("team_id", teamID.String()))

	var members []model.TeamMember
	if err := r.replicaDB.WithContext(ctx).
		Where("team_id = ? AND deleted_at IS NULL", teamID).
		Order("created_at ASC").
		Find(&members).Error; err != nil {
		logger.ErrorCtx(ctx, "Failed to list team members", zap.String("team_id", teamID.String()), zap.Error(err))
		return nil, err
	}
	return members, nil
}

func (r *TeamRepository) IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.TeamMember{}).
		Where("team_id = ? AND user_id = ? AND status = 'active' AND deleted_at IS NULL", teamID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TeamRepository) IsOwner(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.Team{}).
		Where("id = ? AND owner_id = ?", teamID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TeamRepository) CountMembers(ctx context.Context, teamID uuid.UUID) (int64, error) {
	var count int64
	if err := r.replicaDB.WithContext(ctx).
		Model(&model.TeamMember{}).
		Where("team_id = ? AND status = 'active' AND deleted_at IS NULL", teamID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
