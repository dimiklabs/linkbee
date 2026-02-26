export interface CreateLinkRequest {
  destination_url: string;
  slug?: string;
  title?: string;
  password?: string;
  expires_at?: string;
  max_clicks?: number;
  redirect_type?: 301 | 302;
  tags?: string[];
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
  folder_id?: string | null;
}

export interface UpdateLinkRequest {
  destination_url?: string;
  title?: string;
  password?: string;
  expires_at?: string | null;
  max_clicks?: number;
  redirect_type?: 301 | 302;
  tags?: string[];
  is_active?: boolean;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
  folder_id?: string | null;
}

export interface LinkResponse {
  id: string;
  folder_id?: string;
  slug: string;
  short_url: string;
  destination_url: string;
  title?: string;
  click_count: number;
  redirect_type: number;
  is_active: boolean;
  is_starred: boolean;
  is_split_test: boolean;
  health_status: string;
  health_status_code?: number;
  health_checked_at?: string;
  tags?: string[];
  has_password: boolean;
  expires_at?: string;
  max_clicks?: number;
  utm_source?: string;
  utm_medium?: string;
  utm_campaign?: string;
  created_at: string;
  updated_at: string;
}

export interface LinkListResponse {
  links: LinkResponse[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface TimeSeriesPoint {
  timestamp: string;
  count: number;
}

export interface ReferrerPoint {
  referrer: string;
  count: number;
}

export interface DevicePoint {
  device_type: string;
  count: number;
}

export interface CountryPoint {
  country: string;
  count: number;
}

export interface BrowserPoint {
  browser: string;
  count: number;
}

export interface OSPoint {
  os: string;
  count: number;
}

export interface HeatmapPoint {
  day_of_week: number; // 0 = Sunday … 6 = Saturday
  hour: number;        // 0–23
  count: number;
}

export interface AnalyticsResponse {
  link_id: string;
  total_clicks: number;
  unique_clicks: number;
  time_series: TimeSeriesPoint[];
  referrers: ReferrerPoint[];
  devices: DevicePoint[];
  countries: CountryPoint[];
  browsers: BrowserPoint[];
  os_breakdown: OSPoint[];
  heatmap: HeatmapPoint[];
}

export interface DemoShortenRequest {
  destination_url: string;
}

export interface DemoShortenResponse {
  short_url: string;
  slug: string;
  destination_url: string;
}

export interface LinkVariant {
  id: string;
  link_id: string;
  name: string;
  destination_url: string;
  weight: number;
  click_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateVariantRequest {
  name: string;
  destination_url: string;
  weight: number;
}

export interface UpdateVariantRequest {
  name?: string;
  destination_url?: string;
  weight?: number;
}

export interface ImportLinkError {
  row: number;
  url: string;
  error: string;
}

export interface ImportLinksResponse {
  total: number;
  created: number;
  failed: number;
  errors?: ImportLinkError[];
}
