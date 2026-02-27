import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type {
  CreateReportRequest,
  UpdateReportRequest,
  ReportResponse,
  ReportDelivery,
} from '@/types/reports';

export const reportsApi = {
  list: async (): Promise<ApiResponse<ReportResponse[]>> => {
    const response = await apiClient.get('/reports');
    return response.data;
  },

  get: async (id: string): Promise<ApiResponse<ReportResponse>> => {
    const response = await apiClient.get(`/reports/${id}`);
    return response.data;
  },

  create: async (data: CreateReportRequest): Promise<ApiResponse<ReportResponse>> => {
    const response = await apiClient.post('/reports', data);
    return response.data;
  },

  update: async (id: string, data: UpdateReportRequest): Promise<ApiResponse<ReportResponse>> => {
    const response = await apiClient.put(`/reports/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.delete(`/reports/${id}`);
    return response.data;
  },

  sendNow: async (id: string): Promise<ApiResponse<null>> => {
    const response = await apiClient.post(`/reports/${id}/send`);
    return response.data;
  },

  getDeliveries: async (id: string): Promise<ApiResponse<ReportDelivery[]>> => {
    const response = await apiClient.get(`/reports/${id}/deliveries`);
    return response.data;
  },
};

export default reportsApi;
