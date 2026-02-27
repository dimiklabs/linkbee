import apiClient from './client';
import type { AdminStats, AdminUsersResponse, GrowthTimeSeriesResponse, ImpersonationResponse } from '@/types/admin';
import type { ApiResponse } from '@/types/auth';

export default {
  getStats(): Promise<{ data: ApiResponse<AdminStats> }> {
    return apiClient.get('/admin/stats');
  },

  getGrowthTimeSeries(): Promise<{ data: { data: GrowthTimeSeriesResponse } }> {
    return apiClient.get('/admin/growth');
  },

  listUsers(params: { page?: number; limit?: number; search?: string }): Promise<{ data: ApiResponse<AdminUsersResponse> }> {
    return apiClient.get('/admin/users', { params });
  },

  updateUserStatus(id: string, status: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.patch(`/admin/users/${id}/status`, { status });
  },

  updateUserRole(id: string, role: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.patch(`/admin/users/${id}/role`, { role });
  },

  impersonateUser(id: string): Promise<{ data: ApiResponse<ImpersonationResponse> }> {
    return apiClient.post(`/admin/users/${id}/impersonate`);
  },
};
