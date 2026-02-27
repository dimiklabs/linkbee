export interface AdminStats {
  total_users: number;
  active_users: number;
  inactive_users: number;
  total_links: number;
  total_clicks: number;
}

export interface AdminUser {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  status: string;
  role: string;
  auth_provider: string;
  email_verified: boolean;
  created_at: string;
  last_login?: string;
}

export interface AdminUsersResponse {
  users: AdminUser[];
  total: number;
  page: number;
  limit: number;
}

export interface GrowthTimeSeriesPoint {
  timestamp: string;
  count: number;
}

export interface GrowthTimeSeriesResponse {
  users: GrowthTimeSeriesPoint[];
  links: GrowthTimeSeriesPoint[];
}

export interface ImpersonationResponse {
  access_token: string;
  expires_in: number;
  target_user_id: string;
  target_email: string;
}
