package team

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

// ─── Request / Response types ────────────────────────────────────────────────

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Slug        string `json:"slug" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarURL   string `json:"avatar_url"`
}

type InviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required,oneof=admin member"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member"`
}

type AcceptInviteRequest struct {
	Token string `json:"token" binding:"required"`
}

type TeamResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	OwnerID     string `json:"owner_id"`
	Description string `json:"description,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	MemberCount int64  `json:"member_count"`
	CreatedAt   string `json:"created_at"`
}

type MemberResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	Role        string  `json:"role"`
	Status      string  `json:"status"`
	JoinedAt    *string `json:"joined_at,omitempty"`
	InviteEmail string  `json:"invite_email,omitempty"`
}

// ─── Interface ────────────────────────────────────────────────────────────────

type TeamServiceI interface {
	CreateTeam(ctx context.Context, ownerID uuid.UUID, req *CreateTeamRequest) (*TeamResponse, *dto.ServiceError)
	GetTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*TeamResponse, *dto.ServiceError)
	UpdateTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, req *UpdateTeamRequest) (*TeamResponse, *dto.ServiceError)
	DeleteTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) *dto.ServiceError
	ListMyTeams(ctx context.Context, userID uuid.UUID) ([]TeamResponse, *dto.ServiceError)

	InviteMember(ctx context.Context, teamID uuid.UUID, inviterID uuid.UUID, req *InviteMemberRequest) (*MemberResponse, *dto.ServiceError)
	AcceptInvite(ctx context.Context, token string, userID uuid.UUID) *dto.ServiceError
	ListMembers(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) ([]MemberResponse, *dto.ServiceError)
	UpdateMemberRole(ctx context.Context, teamID uuid.UUID, updaterID uuid.UUID, memberUserID uuid.UUID, req *UpdateMemberRoleRequest) *dto.ServiceError
	RemoveMember(ctx context.Context, teamID uuid.UUID, removerID uuid.UUID, memberUserID uuid.UUID) *dto.ServiceError
	LeaveTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) *dto.ServiceError
}

// ─── Implementation ───────────────────────────────────────────────────────────

type teamService struct {
	teamRepo repository.TeamRepositoryI
	userRepo repository.UserRepositoryI
}

func NewTeamService(teamRepo repository.TeamRepositoryI, userRepo repository.UserRepositoryI) TeamServiceI {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func generateInviteToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ─── Team CRUD ────────────────────────────────────────────────────────────────

func (s *teamService) CreateTeam(ctx context.Context, ownerID uuid.UUID, req *CreateTeamRequest) (*TeamResponse, *dto.ServiceError) {
	// Check if slug is already taken
	existing, slugErr := s.teamRepo.GetTeamBySlug(ctx, req.Slug)
	if slugErr != nil && slugErr != gorm.ErrRecordNotFound {
		logger.ErrorCtx(ctx, "Failed to check team slug", zap.Error(slugErr))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if existing != nil {
		return nil, dto.NewConflictError(constant.ErrCodeConflict, "Team slug is already taken. Please choose another.")
	}

	team := &model.Team{
		Name:        req.Name,
		Slug:        req.Slug,
		OwnerID:     ownerID,
		Description: req.Description,
	}

	if err := s.teamRepo.CreateTeam(ctx, team); err != nil {
		logger.ErrorCtx(ctx, "Failed to create team", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	memberCount, _ := s.teamRepo.CountMembers(ctx, team.ID)

	return toTeamResponse(team, memberCount), nil
}

func (s *teamService) GetTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) (*TeamResponse, *dto.ServiceError) {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Check access: owner or active member
	if team.OwnerID != userID {
		isMember, err := s.teamRepo.IsMember(ctx, teamID, userID)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		if !isMember {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
	}

	memberCount, _ := s.teamRepo.CountMembers(ctx, team.ID)
	return toTeamResponse(team, memberCount), nil
}

func (s *teamService) UpdateTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID, req *UpdateTeamRequest) (*TeamResponse, *dto.ServiceError) {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Only owner or admin can update
	if team.OwnerID != userID {
		member, err := s.teamRepo.GetMember(ctx, teamID, userID)
		if err != nil || member.Status != "active" || member.Role == "member" {
			return nil, dto.NewForbiddenError(constant.ErrCodeForbidden, "Only team owners and admins can update team details")
		}
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = req.Description
	}
	if req.AvatarURL != "" {
		team.AvatarURL = req.AvatarURL
	}

	if err := s.teamRepo.UpdateTeam(ctx, team); err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	memberCount, _ := s.teamRepo.CountMembers(ctx, team.ID)
	return toTeamResponse(team, memberCount), nil
}

func (s *teamService) DeleteTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) *dto.ServiceError {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Only owner can delete
	if team.OwnerID != userID {
		return dto.NewForbiddenError(constant.ErrCodeForbidden, "Only the team owner can delete the team")
	}

	if err := s.teamRepo.DeleteTeam(ctx, teamID); err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return nil
}

func (s *teamService) ListMyTeams(ctx context.Context, userID uuid.UUID) ([]TeamResponse, *dto.ServiceError) {
	teams, err := s.teamRepo.ListTeamsByUserID(ctx, userID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	responses := make([]TeamResponse, 0, len(teams))
	for _, t := range teams {
		memberCount, _ := s.teamRepo.CountMembers(ctx, t.ID)
		tc := t
		responses = append(responses, *toTeamResponse(&tc, memberCount))
	}

	return responses, nil
}

// ─── Member management ───────────────────────────────────────────────────────

func (s *teamService) InviteMember(ctx context.Context, teamID uuid.UUID, inviterID uuid.UUID, req *InviteMemberRequest) (*MemberResponse, *dto.ServiceError) {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Only owner or admin can invite
	if team.OwnerID != inviterID {
		member, err := s.teamRepo.GetMember(ctx, teamID, inviterID)
		if err != nil || member.Status != "active" || member.Role == "member" {
			return nil, dto.NewForbiddenError(constant.ErrCodeForbidden, "Only team owners and admins can invite members")
		}
	}

	// Check if email is already invited or a member via email lookup
	// Look up the user by email first
	existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		// Check if the user is already a member
		existingMember, memberErr := s.teamRepo.GetMember(ctx, teamID, existingUser.ID)
		if memberErr == nil && existingMember != nil {
			return nil, dto.NewConflictError(constant.ErrCodeConflict, "User is already a member of this team")
		}
	}

	// Generate invite token
	token, err := generateInviteToken()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate invite token", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	member := &model.TeamMember{
		TeamID:      teamID,
		UserID:      uuid.Nil, // will be set when accepted
		Role:        req.Role,
		InvitedBy:   inviterID,
		InviteEmail: req.Email,
		InviteToken: token,
		Status:      "pending",
	}

	// If user already exists, pre-fill UserID
	if existingUser != nil {
		member.UserID = existingUser.ID
	}

	if err := s.teamRepo.AddMember(ctx, member); err != nil {
		logger.ErrorCtx(ctx, "Failed to add team member", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	logger.InfoCtx(ctx, "Team invite created",
		zap.String("team_id", teamID.String()),
		zap.String("invite_email", req.Email),
		zap.String("invite_token", token))

	return toMemberResponse(member, nil), nil
}

func (s *teamService) AcceptInvite(ctx context.Context, token string, userID uuid.UUID) *dto.ServiceError {
	member, err := s.teamRepo.GetMemberByInviteToken(ctx, token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Invalid or expired invite token")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	if member.Status != "pending" {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "Invite has already been accepted or declined")
	}

	now := time.Now()

	// Update member: set status=active, joinedAt=now
	if err := s.teamRepo.UpdateMemberStatus(ctx, member.ID, "active", &now); err != nil {
		logger.ErrorCtx(ctx, "Failed to accept invite", zap.Error(err))
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	logger.InfoCtx(ctx, "Team invite accepted",
		zap.String("member_id", member.ID.String()),
		zap.String("user_id", userID.String()))

	return nil
}

func (s *teamService) ListMembers(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) ([]MemberResponse, *dto.ServiceError) {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Check access
	if team.OwnerID != userID {
		isMember, err := s.teamRepo.IsMember(ctx, teamID, userID)
		if err != nil {
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		if !isMember {
			return nil, dto.NewForbiddenError(constant.ErrCodeForbidden, "You are not a member of this team")
		}
	}

	members, err := s.teamRepo.ListMembers(ctx, teamID)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	responses := make([]MemberResponse, 0, len(members))
	for _, m := range members {
		mc := m
		var user *model.User
		if mc.UserID != uuid.Nil {
			user, _ = s.userRepo.GetByID(ctx, mc.UserID)
		}
		responses = append(responses, *toMemberResponse(&mc, user))
	}

	return responses, nil
}

func (s *teamService) UpdateMemberRole(ctx context.Context, teamID uuid.UUID, updaterID uuid.UUID, memberUserID uuid.UUID, req *UpdateMemberRoleRequest) *dto.ServiceError {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Only owner or admin can update roles
	if team.OwnerID != updaterID {
		updater, err := s.teamRepo.GetMember(ctx, teamID, updaterID)
		if err != nil || updater.Status != "active" || updater.Role == "member" {
			return dto.NewForbiddenError(constant.ErrCodeForbidden, "Only team owners and admins can change member roles")
		}
	}

	// Cannot change owner's role
	if memberUserID == team.OwnerID {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "Cannot change the role of the team owner")
	}

	if err := s.teamRepo.UpdateMemberRole(ctx, teamID, memberUserID, req.Role); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Member not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return nil
}

func (s *teamService) RemoveMember(ctx context.Context, teamID uuid.UUID, removerID uuid.UUID, memberUserID uuid.UUID) *dto.ServiceError {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Cannot remove the owner
	if memberUserID == team.OwnerID {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "Cannot remove the team owner")
	}

	// Only owner or admin can remove members
	if team.OwnerID != removerID {
		remover, err := s.teamRepo.GetMember(ctx, teamID, removerID)
		if err != nil || remover.Status != "active" || remover.Role == "member" {
			return dto.NewForbiddenError(constant.ErrCodeForbidden, "Only team owners and admins can remove members")
		}
	}

	if err := s.teamRepo.RemoveMember(ctx, teamID, memberUserID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Member not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return nil
}

func (s *teamService) LeaveTeam(ctx context.Context, teamID uuid.UUID, userID uuid.UUID) *dto.ServiceError {
	team, err := s.teamRepo.GetTeamByID(ctx, teamID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "Team not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	// Owner cannot leave — they must delete the team
	if team.OwnerID == userID {
		return dto.NewBadRequestError(constant.ErrCodeBadRequest, "Team owner cannot leave the team. Delete the team instead.")
	}

	if err := s.teamRepo.RemoveMember(ctx, teamID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError(constant.ErrCodeNotFound, "You are not a member of this team")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	return nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func toTeamResponse(team *model.Team, memberCount int64) *TeamResponse {
	return &TeamResponse{
		ID:          team.ID.String(),
		Name:        team.Name,
		Slug:        team.Slug,
		OwnerID:     team.OwnerID.String(),
		Description: team.Description,
		AvatarURL:   team.AvatarURL,
		MemberCount: memberCount,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
	}
}

func toMemberResponse(member *model.TeamMember, user *model.User) *MemberResponse {
	resp := &MemberResponse{
		ID:          member.ID.String(),
		UserID:      member.UserID.String(),
		Role:        member.Role,
		Status:      member.Status,
		InviteEmail: member.InviteEmail,
	}

	if member.JoinedAt != nil {
		t := member.JoinedAt.Format(time.RFC3339)
		resp.JoinedAt = &t
	}

	if user != nil {
		resp.Email = user.Email
		resp.FirstName = user.FirstName
		resp.LastName = user.LastName
	} else {
		resp.Email = member.InviteEmail
	}

	return resp
}
