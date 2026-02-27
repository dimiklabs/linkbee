import type { LinkResponse, TimeSeriesPoint, CountryPoint, DevicePoint, BrowserPoint, ReferrerPoint, OSPoint, CityPoint, SourcePoint, ReferrerCategoryPoint, UTMPoint } from './links';

export interface GlobalAnalyticsResponse {
  from: string;
  to: string;
  total_clicks: number;
  unique_clicks: number;
  time_series: TimeSeriesPoint[];
  top_countries: CountryPoint[];
  device_breakdown: DevicePoint[];
  top_browsers: BrowserPoint[];
  top_referrers: ReferrerPoint[];
  os_breakdown: OSPoint[];
  top_cities: CityPoint[];
  source_breakdown: SourcePoint[];
  referrer_categories: ReferrerCategoryPoint[];
  utm_sources: UTMPoint[];
  utm_mediums: UTMPoint[];
  utm_campaigns: UTMPoint[];
}

export interface GlobalAnalyticsPeriodData {
  from: string;
  to: string;
  total_clicks: number;
  unique_clicks: number;
}

export interface TrendData {
  trend: 'up' | 'down' | 'flat';
  percent_change: number;
}

export interface GlobalAnalyticsComparisonResponse {
  current: GlobalAnalyticsPeriodData;
  previous: GlobalAnalyticsPeriodData;
  clicks: TrendData;
  unique_clicks: TrendData;
}

export interface DashboardOverviewResponse {
  total_links: number;
  total_clicks: number;
  clicks_today: number;
  clicks_30_days: number;
  clicks_7_days: number;
  time_series_30d: TimeSeriesPoint[];
  top_links: LinkResponse[];
  recent_links: LinkResponse[];
}
