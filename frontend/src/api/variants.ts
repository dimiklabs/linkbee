import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { LinkVariant, CreateVariantRequest, UpdateVariantRequest } from '@/types/links';

export const variantsApi = {
  list: async (linkId: string): Promise<ApiResponse<LinkVariant[]>> => {
    const response = await apiClient.get(`/links/${linkId}/variants`);
    return response.data;
  },

  create: async (linkId: string, data: CreateVariantRequest): Promise<ApiResponse<LinkVariant>> => {
    const response = await apiClient.post(`/links/${linkId}/variants`, data);
    return response.data;
  },

  update: async (linkId: string, variantId: string, data: UpdateVariantRequest): Promise<ApiResponse<LinkVariant>> => {
    const response = await apiClient.put(`/links/${linkId}/variants/${variantId}`, data);
    return response.data;
  },

  delete: async (linkId: string, variantId: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.delete(`/links/${linkId}/variants/${variantId}`);
    return response.data;
  },

  toggleSplitTest: async (linkId: string, enabled: boolean): Promise<ApiResponse<{ is_split_test: boolean }>> => {
    const response = await apiClient.patch(`/links/${linkId}/split-test`, { enabled });
    return response.data;
  },
};

export default variantsApi;
