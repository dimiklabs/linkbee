import apiClient from './client';
import type { ApiResponse } from '@/types/auth';
import type { AuditLogsResponse } from '@/types/audit';

export default {
  list(params: {
    page?: number
    limit?: number
    action?: string
    resource_type?: string
    from?: string
    to?: string
  }): Promise<{ data: ApiResponse<AuditLogsResponse> }> {
    return apiClient.get('/audit-logs', { params });
  },
};
