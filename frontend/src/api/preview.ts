import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { LinkPreviewData } from '@/types/preview';

export const previewApi = {
  get: async (linkId: string): Promise<ApiResponse<LinkPreviewData>> => {
    const res = await apiClient.get(`/links/${linkId}/preview`);
    return res.data;
  },
};

export default previewApi;
