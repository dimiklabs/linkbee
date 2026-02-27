export interface CreateReportRequest {
  name: string;
  link_ids: string[];
  frequency: 'daily' | 'weekly' | 'monthly';
}

export interface UpdateReportRequest {
  name?: string;
  link_ids?: string[];
  frequency?: 'daily' | 'weekly' | 'monthly';
  is_active?: boolean;
}

export interface ReportResponse {
  id: string;
  user_id: string;
  name: string;
  link_ids: string[];
  frequency: 'daily' | 'weekly' | 'monthly';
  next_run_at?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ReportDelivery {
  id: string;
  report_id: string;
  status: 'sent' | 'failed';
  failure_reason?: string;
  delivered_at?: string;
  created_at: string;
}
