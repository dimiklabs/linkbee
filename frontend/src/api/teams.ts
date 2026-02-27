import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { Team, TeamMember, CreateTeamRequest, InviteMemberRequest } from '@/types/teams';

export default {
  listMyTeams(): Promise<{ data: ApiResponse<Team[]> }> {
    return apiClient.get('/teams');
  },
  createTeam(data: CreateTeamRequest): Promise<{ data: ApiResponse<Team> }> {
    return apiClient.post('/teams', data);
  },
  getTeam(id: string): Promise<{ data: ApiResponse<Team> }> {
    return apiClient.get(`/teams/${id}`);
  },
  updateTeam(id: string, data: Partial<CreateTeamRequest> & { avatar_url?: string }): Promise<{ data: ApiResponse<Team> }> {
    return apiClient.put(`/teams/${id}`, data);
  },
  deleteTeam(id: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.delete(`/teams/${id}`);
  },
  listMembers(teamId: string): Promise<{ data: ApiResponse<TeamMember[]> }> {
    return apiClient.get(`/teams/${teamId}/members`);
  },
  inviteMember(teamId: string, data: InviteMemberRequest): Promise<{ data: ApiResponse<TeamMember> }> {
    return apiClient.post(`/teams/${teamId}/members`, data);
  },
  updateMemberRole(teamId: string, userId: string, role: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.patch(`/teams/${teamId}/members/${userId}/role`, { role });
  },
  removeMember(teamId: string, userId: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.delete(`/teams/${teamId}/members/${userId}`);
  },
  leaveTeam(teamId: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.post(`/teams/${teamId}/leave`);
  },
  acceptInvite(token: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.post('/teams/invite/accept', { token });
  },
};
