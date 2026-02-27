import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { DeliveriesResponse, Webhook, CreateWebhookRequest, UpdateWebhookRequest, WebhookDelivery } from '@/types/webhooks';

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

  getSecret: async (id: string): Promise<ApiResponse<{ secret: string }>> => {
    const response = await apiClient.get(`/webhooks/${id}/secret`);
    return response.data;
  },

  test: async (id: string): Promise<ApiResponse<WebhookDelivery>> => {
    const response = await apiClient.post(`/webhooks/${id}/test`);
    return response.data;
  },

  getDeliveries: async (id: string, page = 1, limit = 20): Promise<ApiResponse<DeliveriesResponse>> => {
    const response = await apiClient.get(`/webhooks/${id}/deliveries`, { params: { page, limit } });
    return response.data;
  },

  resendDelivery: async (webhookId: string, deliveryId: string): Promise<ApiResponse<WebhookDelivery>> => {
    const response = await apiClient.post(`/webhooks/${webhookId}/deliveries/${deliveryId}/resend`);
    return response.data;
  },
};

export default webhooksApi;
