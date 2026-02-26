import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { RetargetingPixel, CreatePixelRequest } from '@/types/pixels';

export const pixelsApi = {
  list: async (linkId: string): Promise<ApiResponse<RetargetingPixel[]>> => {
    const response = await apiClient.get(`/links/${linkId}/pixels`);
    return response.data;
  },

  create: async (linkId: string, data: CreatePixelRequest): Promise<ApiResponse<RetargetingPixel>> => {
    const response = await apiClient.post(`/links/${linkId}/pixels`, data);
    return response.data;
  },

  delete: async (linkId: string, pixelId: string): Promise<void> => {
    await apiClient.delete(`/links/${linkId}/pixels/${pixelId}`);
  },

  toggle: async (linkId: string, enabled: boolean): Promise<void> => {
    await apiClient.patch(`/links/${linkId}/pixel-tracking`, { enabled });
  },
};

export default pixelsApi;
