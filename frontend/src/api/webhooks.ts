import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { Webhook, CreateWebhookRequest, UpdateWebhookRequest } from '@/types/webhooks';

export const webhooksApi = {
  list: async (): Promise<ApiResponse<Webhook[]>> => {
    const response = await apiClient.get('/webhooks');
    return response.data;
  },

  create: async (data: CreateWebhookRequest): Promise<ApiResponse<Webhook>> => {
    const response = await apiClient.post('/webhooks', data);
    return response.data;
  },

  update: async (id: string, data: UpdateWebhookRequest): Promise<ApiResponse<Webhook>> => {
    const response = await apiClient.put(`/webhooks/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/webhooks/${id}`);
  },
};

export default webhooksApi;
