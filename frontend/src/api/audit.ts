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

  async exportAuditLogs(action?: string, from?: string, to?: string): Promise<void> {
    const response = await apiClient.get('/audit-logs/export', {
      params: { action, from, to },
      responseType: 'blob',
    });
    const url = window.URL.createObjectURL(new Blob([response.data as BlobPart]));
    const a = document.createElement('a');
    a.href = url;
    a.download = `audit-logs-${new Date().toISOString().split('T')[0]}.csv`;
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
  },
};
