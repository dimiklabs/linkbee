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
  is_geo_routing: boolean;
  is_pixel_tracking: boolean;
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

export interface CityPoint {
  city: string;
  country: string;
  count: number;
}

export interface ReferrerCategoryPoint {
  category: 'direct' | 'search' | 'social' | 'email' | 'referral';
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

export interface UTMPoint {
  value: string;
  count: number;
}

export interface SourcePoint {
  source: string;
  count: number;
}

export interface AnalyticsResponse {
  link_id: string;
  total_clicks: number;
  unique_clicks: number;
  first_time_visitors: number;
  returning_visitors: number;
  time_series: TimeSeriesPoint[];
  referrers: ReferrerPoint[];
  devices: DevicePoint[];
  countries: CountryPoint[];
  browsers: BrowserPoint[];
  os_breakdown: OSPoint[];
  sources: SourcePoint[];
  heatmap: HeatmapPoint[];
  utm_sources: UTMPoint[];
  utm_mediums: UTMPoint[];
  utm_campaigns: UTMPoint[];
  utm_contents: UTMPoint[];
  utm_terms: UTMPoint[];
  cities: CityPoint[];
  referrer_categories: ReferrerCategoryPoint[];
}

export interface PeriodMetrics {
  from: string;
  to: string;
  total_clicks: number;
  unique_clicks: number;
}

export interface PeriodChange {
  count_change: number;
  percent_change: number;
  trend: 'up' | 'down' | 'stable';
}

export interface PeriodComparisonResponse {
  link_id: string;
  current: PeriodMetrics;
  previous: PeriodMetrics;
  clicks: PeriodChange;
  unique_clicks: PeriodChange;
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

export interface LinkGeoRule {
  id: string;
  link_id: string;
  country_code: string;
  destination_url: string;
  priority: number;
}

export interface CreateGeoRuleRequest {
  country_code: string;
  destination_url: string;
  priority: number;
}

export interface UpdateGeoRuleRequest {
  country_code?: string;
  destination_url?: string;
  priority?: number;
}

export interface LinkComparisonMetric {
  link_id: string;
  slug: string;
  short_url: string;
  title?: string;
  total_clicks: number;
  unique_clicks: number;
  clicks_per_day: number;
  top_referrer?: string;
  top_country?: string;
  top_browser?: string;
  top_device?: string;
}

export interface MultiLinkComparisonResponse {
  links: LinkComparisonMetric[];
  from: string;
  to: string;
  span_days: number;
}

export interface CloneLinkRequest {
  new_title?: string;
  new_slug?: string;
}

export interface BulkLinkActionRequest {
  ids: string[];
  action: 'delete' | 'activate' | 'deactivate' | 'move_folder' | 'add_tags' | 'remove_tags';
  folder_id?: string | null;
  tags?: string[];
}

export interface BulkActionResponse {
  affected: number;
  action: string;
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
