import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type {
  CreateLinkRequest,
  UpdateLinkRequest,
  LinkResponse,
  LinkListResponse,
  AnalyticsResponse,
  DemoShortenRequest,
  DemoShortenResponse,
  ImportLinksResponse,
} from '@/types/links';

export const linksApi = {
  list: async (page = 1, limit = 20, search = '', folderID = '', starred?: boolean, healthStatus?: string): Promise<ApiResponse<LinkListResponse>> => {
    const response = await apiClient.get('/links', { params: { page, limit, search, folder_id: folderID || undefined, starred, health_status: healthStatus || undefined } });
    return response.data;
  },

  create: async (data: CreateLinkRequest): Promise<ApiResponse<LinkResponse>> => {
    const response = await apiClient.post('/links', data);
    return response.data;
  },

  get: async (id: string): Promise<ApiResponse<LinkResponse>> => {
    const response = await apiClient.get(`/links/${id}`);
    return response.data;
  },

  update: async (id: string, data: UpdateLinkRequest): Promise<ApiResponse<LinkResponse>> => {
    const response = await apiClient.put(`/links/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.delete(`/links/${id}`);
    return response.data;
  },

  toggleStar: async (id: string): Promise<ApiResponse<LinkResponse>> => {
    const response = await apiClient.patch(`/links/${id}/star`);
    return response.data;
  },

  checkHealth: async (id: string): Promise<ApiResponse<LinkResponse>> => {
    const response = await apiClient.post(`/links/${id}/health-check`);
    return response.data;
  },

  importCSV: async (file: File): Promise<ApiResponse<ImportLinksResponse>> => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await apiClient.post('/links/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  getQRUrl: (id: string, opts?: { fg?: string; bg?: string; size?: number; ec?: string }): string => {
    const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';
    const token = localStorage.getItem('access_token');
    const params = new URLSearchParams({ token: token ?? '' });
    if (opts?.fg) params.set('fg', opts.fg.replace('#', ''));
    if (opts?.bg) params.set('bg', opts.bg.replace('#', ''));
    if (opts?.size) params.set('size', String(opts.size));
    if (opts?.ec) params.set('ec', opts.ec);
    return `${base}/links/${id}/qr?${params.toString()}`;
  },

  getLiveCountUrl: (id: string): string => {
    const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';
    const token = localStorage.getItem('access_token');
    return `${base}/links/${id}/analytics/live?token=${token}`;
  },

  getAnalytics: async (id: string, from?: string, to?: string, granularity = 'day'): Promise<ApiResponse<AnalyticsResponse>> => {
    const response = await apiClient.get(`/links/${id}/analytics`, {
      params: { from, to, granularity },
    });
    return response.data;
  },

  demoShorten: async (data: DemoShortenRequest): Promise<ApiResponse<DemoShortenResponse>> => {
    const response = await apiClient.post('/demo/shorten', data);
    return response.data;
  },
};

export default linksApi;
