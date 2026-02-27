import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { DashboardOverviewResponse, GlobalAnalyticsResponse, GlobalAnalyticsComparisonResponse } from '@/types/dashboard';

export const dashboardApi = {
  getOverview: async (): Promise<ApiResponse<DashboardOverviewResponse>> => {
    const response = await apiClient.get('/dashboard/overview');
    return response.data;
  },

  getGlobalAnalytics: async (from?: string, to?: string): Promise<ApiResponse<GlobalAnalyticsResponse>> => {
    const response = await apiClient.get('/dashboard/analytics', { params: { from, to } });
    return response.data;
  },

  getGlobalAnalyticsComparison: async (from?: string, to?: string): Promise<ApiResponse<GlobalAnalyticsComparisonResponse>> => {
    const response = await apiClient.get('/dashboard/analytics/comparison', { params: { from, to } });
    return response.data;
  },
};

export default dashboardApi;
