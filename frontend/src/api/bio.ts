import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type {
  BioPage,
  BioLinkItem,
  UpdateBioPageRequest,
  CreateBioLinkRequest,
  UpdateBioLinkRequest,
  ReorderBioLinksRequest,
} from '@/types/bio';

export const bioApi = {
  get: async (): Promise<ApiResponse<BioPage>> => {
    const res = await apiClient.get('/bio');
    return res.data;
  },

  update: async (data: UpdateBioPageRequest): Promise<ApiResponse<BioPage>> => {
    const res = await apiClient.put('/bio', data);
    return res.data;
  },

  getPublic: async (username: string): Promise<ApiResponse<BioPage>> => {
    const res = await apiClient.get(`/bio/public/${username}`);
    return res.data;
  },

  listLinks: async (): Promise<ApiResponse<BioLinkItem[]>> => {
    const res = await apiClient.get('/bio/links');
    return res.data;
  },

  createLink: async (data: CreateBioLinkRequest): Promise<ApiResponse<BioLinkItem>> => {
    const res = await apiClient.post('/bio/links', data);
    return res.data;
  },

  updateLink: async (id: string, data: UpdateBioLinkRequest): Promise<ApiResponse<BioLinkItem>> => {
    const res = await apiClient.put(`/bio/links/${id}`, data);
    return res.data;
  },

  deleteLink: async (id: string): Promise<void> => {
    await apiClient.delete(`/bio/links/${id}`);
  },

  reorderLinks: async (data: ReorderBioLinksRequest): Promise<void> => {
    await apiClient.patch('/bio/links/reorder', data);
  },
};

export default bioApi;
