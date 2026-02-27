package response

// GlobalAnalyticsResponse is the account-wide analytics summary.
type GlobalAnalyticsResponse struct {
	From               string                  `json:"from"`
	To                 string                  `json:"to"`
	TotalClicks        int64                   `json:"total_clicks"`
	UniqueClicks       int64                   `json:"unique_clicks"`
	TimeSeries         []TimeSeriesData        `json:"time_series"`
	TopCountries       []CountryData           `json:"top_countries"`
	DeviceBreakdown    []DeviceData            `json:"device_breakdown"`
	TopBrowsers        []BrowserData           `json:"top_browsers"`
	TopReferrers       []ReferrerData          `json:"top_referrers"`
	OSBreakdown        []OSData                `json:"os_breakdown"`
	TopCities          []CityData              `json:"top_cities"`
	SourceBreakdown    []SourceData            `json:"source_breakdown"`
	ReferrerCategories []ReferrerCategoryData  `json:"referrer_categories"`
	UTMSources         []UTMData               `json:"utm_sources"`
	UTMMediums         []UTMData               `json:"utm_mediums"`
	UTMCampaigns       []UTMData               `json:"utm_campaigns"`
}

// GlobalAnalyticsPeriodData holds aggregated stats for a single time window.
type GlobalAnalyticsPeriodData struct {
	From         string `json:"from"`
	To           string `json:"to"`
	TotalClicks  int64  `json:"total_clicks"`
	UniqueClicks int64  `json:"unique_clicks"`
}

// TrendData describes the directional change between two periods.
type TrendData struct {
	Trend         string  `json:"trend"`          // "up", "down", "flat"
	PercentChange float64 `json:"percent_change"` // absolute value of the percentage change
}

// GlobalAnalyticsComparisonResponse is returned by GET /dashboard/analytics/comparison.
type GlobalAnalyticsComparisonResponse struct {
	Current      GlobalAnalyticsPeriodData `json:"current"`
	Previous     GlobalAnalyticsPeriodData `json:"previous"`
	Clicks       TrendData                 `json:"clicks"`
	UniqueClicks TrendData                 `json:"unique_clicks"`
}

// DashboardOverviewResponse is the aggregated summary returned to the dashboard home page.
type DashboardOverviewResponse struct {
	TotalLinks    int64            `json:"total_links"`
	TotalClicks   int64            `json:"total_clicks"`
	ClicksToday   int64            `json:"clicks_today"`
	Clicks30Days  int64            `json:"clicks_30_days"`
	Clicks7Days   int64            `json:"clicks_7_days"`
	TimeSeries30d []TimeSeriesData `json:"time_series_30d"`
	TopLinks      []*LinkResponse  `json:"top_links"`
	RecentLinks   []*LinkResponse  `json:"recent_links"`
}
