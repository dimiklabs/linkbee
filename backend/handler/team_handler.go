package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/middlewares"
	teamSvc "github.com/shafikshaon/linkbee/service/team"
	"github.com/shafikshaon/linkbee/transport"
)

type TeamHandler struct {
	teamService teamSvc.TeamServiceI
}

func NewTeamHandler(teamService teamSvc.TeamServiceI) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// ListMyTeams returns all teams the authenticated user belongs to.
func (h *TeamHandler) ListMyTeams(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	result, svcErr := h.teamService.ListMyTeams(ctx, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Teams retrieved successfully", result)
}

// CreateTeam creates a new team owned by the authenticated user.
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req teamSvc.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	result, svcErr := h.teamService.CreateTeam(ctx, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusCreated, "Team created successfully", result)
}

// GetTeam returns a single team by ID.
func (h *TeamHandler) GetTeam(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	result, svcErr := h.teamService.GetTeam(ctx, teamID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Team retrieved successfully", result)
}

// UpdateTeam updates a team's details.
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	var req teamSvc.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	result, svcErr := h.teamService.UpdateTeam(ctx, teamID, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Team updated successfully", result)
}

// DeleteTeam permanently soft-deletes a team (owner only).
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	if svcErr := h.teamService.DeleteTeam(ctx, teamID, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Team deleted successfully", nil)
}

// InviteMember invites a user to the team by email.
func (h *TeamHandler) InviteMember(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	var req teamSvc.InviteMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	result, svcErr := h.teamService.InviteMember(ctx, teamID, userID, &req)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusCreated, "Invitation sent successfully", result)
}

// AcceptInvite accepts a team invite using a token.
func (h *TeamHandler) AcceptInvite(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	var req teamSvc.AcceptInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	if svcErr := h.teamService.AcceptInvite(ctx, req.Token, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Invitation accepted successfully", nil)
}

// ListMembers lists all members of a team.
func (h *TeamHandler) ListMembers(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	result, svcErr := h.teamService.ListMembers(ctx, teamID, userID)
	if svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Members retrieved successfully", result)
}

// UpdateMemberRole updates the role of a team member.
func (h *TeamHandler) UpdateMemberRole(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	memberUserID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid user ID")
		return
	}

	var req teamSvc.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeValidationError, err.Error())
		return
	}

	if svcErr := h.teamService.UpdateMemberRole(ctx, teamID, userID, memberUserID, &req); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Member role updated successfully", nil)
}

// RemoveMember removes a member from a team.
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	memberUserID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid user ID")
		return
	}

	if svcErr := h.teamService.RemoveMember(ctx, teamID, userID, memberUserID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Member removed successfully", nil)
}

// LeaveTeam allows a member to leave a team.
func (h *TeamHandler) LeaveTeam(c *gin.Context) {
	ctx := c.Request.Context()
	userIDStr, _ := c.Get(middlewares.ContextKeyUserID)
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		transport.RespondWithError(c, http.StatusUnauthorized, constant.ErrCodeUnauthorized, constant.ErrMsgUnauthorized)
		return
	}

	teamID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		transport.RespondWithError(c, http.StatusBadRequest, constant.ErrCodeBadRequest, "Invalid team ID")
		return
	}

	if svcErr := h.teamService.LeaveTeam(ctx, teamID, userID); svcErr != nil {
		transport.RespondWithError(c, svcErr.StatusCode, svcErr.ErrorCode, svcErr.Description)
		return
	}

	transport.RespondWithSuccess(c, http.StatusOK, "Left team successfully", nil)
}
