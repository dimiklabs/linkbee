import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { APIKey, CreateAPIKeyRequest, CreateAPIKeyResponse } from '@/types/apikeys';

export const apiKeysApi = {
  list: async (): Promise<ApiResponse<APIKey[]>> => {
    const response = await apiClient.get('/api-keys');
    return response.data;
  },

  create: async (data: CreateAPIKeyRequest): Promise<ApiResponse<CreateAPIKeyResponse>> => {
    const response = await apiClient.post('/api-keys', data);
    return response.data;
  },

  revoke: async (id: string): Promise<void> => {
    await apiClient.delete(`/api-keys/${id}`);
  },
};

export default apiKeysApi;
