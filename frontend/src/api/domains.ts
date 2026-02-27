import apiClient from './client';
import type { DomainResponse } from '@/types/domains';
import type { ApiResponse } from '@/types/auth';

export default {
  list(): Promise<{ data: ApiResponse<DomainResponse[]> }> {
    return apiClient.get('/domains');
  },

  add(domain: string): Promise<{ data: ApiResponse<DomainResponse> }> {
    return apiClient.post('/domains', { domain });
  },

  verify(id: string): Promise<{ data: ApiResponse<DomainResponse> }> {
    return apiClient.post(`/domains/${id}/verify`);
  },

  remove(id: string): Promise<{ data: ApiResponse<null> }> {
    return apiClient.delete(`/domains/${id}`);
  },
};
