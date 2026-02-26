import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { LinkGeoRule, CreateGeoRuleRequest, UpdateGeoRuleRequest, LinkResponse } from '@/types/links';

export const geoApi = {
  list: async (linkId: string): Promise<ApiResponse<LinkGeoRule[]>> => {
    const response = await apiClient.get(`/links/${linkId}/geo-rules`);
    return response.data;
  },

  create: async (linkId: string, data: CreateGeoRuleRequest): Promise<ApiResponse<LinkGeoRule>> => {
    const response = await apiClient.post(`/links/${linkId}/geo-rules`, data);
    return response.data;
  },

  update: async (linkId: string, ruleId: string, data: UpdateGeoRuleRequest): Promise<ApiResponse<LinkGeoRule>> => {
    const response = await apiClient.put(`/links/${linkId}/geo-rules/${ruleId}`, data);
    return response.data;
  },

  delete: async (linkId: string, ruleId: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.delete(`/links/${linkId}/geo-rules/${ruleId}`);
    return response.data;
  },

  toggleGeoRouting: async (linkId: string, enabled: boolean): Promise<ApiResponse<LinkResponse>> => {
    const response = await apiClient.patch(`/links/${linkId}/geo-routing`, { enabled });
    return response.data;
  },
};

export default geoApi;
