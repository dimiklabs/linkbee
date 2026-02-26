import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { FolderResponse, CreateFolderRequest, UpdateFolderRequest } from '@/types/folders';

export const foldersApi = {
  list: async (): Promise<ApiResponse<FolderResponse[]>> => {
    const response = await apiClient.get('/folders');
    return response.data;
  },

  create: async (data: CreateFolderRequest): Promise<ApiResponse<FolderResponse>> => {
    const response = await apiClient.post('/folders', data);
    return response.data;
  },

  update: async (id: string, data: UpdateFolderRequest): Promise<ApiResponse<FolderResponse>> => {
    const response = await apiClient.put(`/folders/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.delete(`/folders/${id}`);
    return response.data;
  },
};

export default foldersApi;
