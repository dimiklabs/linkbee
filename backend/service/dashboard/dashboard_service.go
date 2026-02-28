package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/response"
)

type DashboardServiceI interface {
	GetOverview(ctx context.Context, userID uuid.UUID) (*response.DashboardOverviewResponse, *dto.ServiceError)
	GetGlobalAnalytics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*response.GlobalAnalyticsResponse, *dto.ServiceError)
	GetGlobalAnalyticsComparison(ctx context.Context, userID uuid.UUID, from, to time.Time) (*response.GlobalAnalyticsComparisonResponse, *dto.ServiceError)
}

type dashboardService struct {
	linkRepo       repository.LinkRepositoryI
	clickEventRepo repository.ClickEventRepositoryI
	appCfg         *config.AppConfig
}

func NewDashboardService(linkRepo repository.LinkRepositoryI, clickEventRepo repository.ClickEventRepositoryI, appCfg *config.AppConfig) DashboardServiceI {
	return &dashboardService{
		linkRepo:       linkRepo,
		clickEventRepo: clickEventRepo,
		appCfg:         appCfg,
	}
}

func (s *dashboardService) GetOverview(ctx context.Context, userID uuid.UUID) (*response.DashboardOverviewResponse, *dto.ServiceError) {
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	from30d := now.AddDate(0, 0, -30)

	totalLinks, err := s.linkRepo.CountByUserID(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to count links for dashboard", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	totalClicks, err := s.clickEventRepo.GetTotalClicksByUserID(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get total clicks for dashboard", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	clicksToday, err := s.clickEventRepo.GetClicksInPeriod(ctx, userID, todayStart, now)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get today's clicks for dashboard", zap.Error(err))
	}

	clicks30Days, err := s.clickEventRepo.GetClicksInPeriod(ctx, userID, from30d, now)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get 30-day clicks for dashboard", zap.Error(err))
	}

	clicks7Days, err := s.clickEventRepo.GetClicksInPeriod(ctx, userID, now.AddDate(0, 0, -7), now)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get 7-day clicks for dashboard", zap.Error(err))
	}

	rawSeries, err := s.clickEventRepo.GetTimeSeriesByUserID(ctx, userID, from30d, now, "day")
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get 30-day time series for dashboard", zap.Error(err))
	}
	tsData := make([]response.TimeSeriesData, len(rawSeries))
	for i, pt := range rawSeries {
		tsData[i] = response.TimeSeriesData{Timestamp: pt.Timestamp, Count: pt.Count}
	}

	topLinks, err := s.linkRepo.GetTopByClicks(ctx, userID, 5)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get top links for dashboard", zap.Error(err))
	}

	recentLinks, err := s.linkRepo.GetRecentByUserID(ctx, userID, 5)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get recent links for dashboard", zap.Error(err))
	}

	return &response.DashboardOverviewResponse{
		TotalLinks:    totalLinks,
		TotalClicks:   totalClicks,
		ClicksToday:   clicksToday,
		Clicks30Days:  clicks30Days,
		Clicks7Days:   clicks7Days,
		TimeSeries30d: tsData,
		TopLinks:      s.toResponseList(topLinks),
		RecentLinks:   s.toResponseList(recentLinks),
	}, nil
}

func (s *dashboardService) toResponseList(links []model.Link) []*response.LinkResponse {
	out := make([]*response.LinkResponse, len(links))
	for i := range links {
		out[i] = s.toLinkResponse(&links[i])
	}
	return out
}

func (s *dashboardService) toLinkResponse(link *model.Link) *response.LinkResponse {
	shortURL := fmt.Sprintf("%s/%s", s.appCfg.BaseDomain, link.Slug)
	tags := []string{}
	if link.Tags != nil {
		tags = link.Tags
	}
	return &response.LinkResponse{
		ID:               link.ID,
		FolderID:         link.FolderID,
		Slug:             link.Slug,
		ShortURL:         shortURL,
		DestinationURL:   link.DestinationURL,
		Title:            link.Title,
		ClickCount:       link.ClickCount,
		RedirectType:     link.RedirectType,
		IsActive:         link.IsActive,
		IsStarred:        link.IsStarred,
		IsSplitTest:      link.IsSplitTest,
		IsGeoRouting:     link.IsGeoRouting,
		IsPixelTracking:  link.IsPixelTracking,
		HealthStatus:     link.HealthStatus,
		HealthStatusCode: link.HealthStatusCode,
		HealthCheckedAt:  link.HealthCheckedAt,
		Tags:             tags,
		HasPassword:      link.PasswordHash != "",
		ExpiresAt:        link.ExpiresAt,
		MaxClicks:        link.MaxClicks,
		UTMSource:        link.UTMSource,
		UTMMedium:        link.UTMMedium,
		UTMCampaign:      link.UTMCampaign,
		CreatedAt:        link.CreatedAt,
		UpdatedAt:        link.UpdatedAt,
	}
}

func (s *dashboardService) GetGlobalAnalytics(ctx context.Context, userID uuid.UUID, from, to time.Time) (*response.GlobalAnalyticsResponse, *dto.ServiceError) {
	totalClicks, err := s.clickEventRepo.GetClicksInPeriod(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get total clicks for global analytics", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	uniqueClicks, err := s.clickEventRepo.GetUniqueClicksInPeriodByUserID(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get unique clicks for global analytics", zap.Error(err))
	}

	rawSeries, err := s.clickEventRepo.GetTimeSeriesByUserID(ctx, userID, from, to, "day")
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get time series for global analytics", zap.Error(err))
	}
	tsData := make([]response.TimeSeriesData, len(rawSeries))
	for i, pt := range rawSeries {
		tsData[i] = response.TimeSeriesData{Timestamp: pt.Timestamp, Count: pt.Count}
	}

	countries, err := s.clickEventRepo.GetCountryBreakdownByUserID(ctx, userID, from, to, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get country breakdown for global analytics", zap.Error(err))
	}
	countryData := make([]response.CountryData, len(countries))
	for i, c := range countries {
		countryData[i] = response.CountryData{Country: c.Country, Count: c.Count}
	}

	devices, err := s.clickEventRepo.GetDeviceBreakdownByUserID(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get device breakdown for global analytics", zap.Error(err))
	}
	deviceData := make([]response.DeviceData, len(devices))
	for i, d := range devices {
		deviceData[i] = response.DeviceData{DeviceType: d.DeviceType, Count: d.Count}
	}

	browsers, err := s.clickEventRepo.GetBrowserBreakdownByUserID(ctx, userID, from, to, 8)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get browser breakdown for global analytics", zap.Error(err))
	}
	browserData := make([]response.BrowserData, len(browsers))
	for i, b := range browsers {
		browserData[i] = response.BrowserData{Browser: b.Browser, Count: b.Count}
	}

	referrers, err := s.clickEventRepo.GetTopReferrersByUserID(ctx, userID, from, to, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get referrers for global analytics", zap.Error(err))
	}
	referrerData := make([]response.ReferrerData, len(referrers))
	for i, r := range referrers {
		referrerData[i] = response.ReferrerData{Referrer: r.Referrer, Count: r.Count}
	}

	// OS breakdown
	osPoints, err := s.clickEventRepo.GetOSBreakdownByUserID(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get OS breakdown for global analytics", zap.Error(err))
	}
	osData := make([]response.OSData, len(osPoints))
	for i, o := range osPoints {
		osData[i] = response.OSData{OS: o.OS, Count: o.Count}
	}

	// City breakdown
	cityPoints, err := s.clickEventRepo.GetCityBreakdownByUserID(ctx, userID, from, to, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get city breakdown for global analytics", zap.Error(err))
	}
	cityData := make([]response.CityData, len(cityPoints))
	for i, c := range cityPoints {
		cityData[i] = response.CityData{City: c.City, Country: c.Country, Count: c.Count}
	}

	// Source breakdown
	sourcePoints, err := s.clickEventRepo.GetSourceBreakdownByUserID(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get source breakdown for global analytics", zap.Error(err))
	}
	sourceData := make([]response.SourceData, len(sourcePoints))
	for i, s := range sourcePoints {
		sourceData[i] = response.SourceData{Source: s.Source, Count: s.Count}
	}

	// Referrer category breakdown
	refCatPoints, err := s.clickEventRepo.GetReferrerCategoryBreakdownByUserID(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get referrer category breakdown for global analytics", zap.Error(err))
	}
	refCatData := make([]response.ReferrerCategoryData, len(refCatPoints))
	for i, rc := range refCatPoints {
		refCatData[i] = response.ReferrerCategoryData{Category: rc.Category, Count: rc.Count}
	}

	// UTM breakdowns
	utmSrcPoints, err := s.clickEventRepo.GetUTMBreakdownByUserID(ctx, userID, "utm_source", from, to, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM source breakdown for global analytics", zap.Error(err))
	}
	utmSrcData := make([]response.UTMData, len(utmSrcPoints))
	for i, u := range utmSrcPoints {
		utmSrcData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	utmMedPoints, err := s.clickEventRepo.GetUTMBreakdownByUserID(ctx, userID, "utm_medium", from, to, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM medium breakdown for global analytics", zap.Error(err))
	}
	utmMedData := make([]response.UTMData, len(utmMedPoints))
	for i, u := range utmMedPoints {
		utmMedData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	utmCamPoints, err := s.clickEventRepo.GetUTMBreakdownByUserID(ctx, userID, "utm_campaign", from, to, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM campaign breakdown for global analytics", zap.Error(err))
	}
	utmCamData := make([]response.UTMData, len(utmCamPoints))
	for i, u := range utmCamPoints {
		utmCamData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	return &response.GlobalAnalyticsResponse{
		From:               from.Format(time.RFC3339),
		To:                 to.Format(time.RFC3339),
		TotalClicks:        totalClicks,
		UniqueClicks:       uniqueClicks,
		TimeSeries:         tsData,
		TopCountries:       countryData,
		DeviceBreakdown:    deviceData,
		TopBrowsers:        browserData,
		TopReferrers:       referrerData,
		OSBreakdown:        osData,
		TopCities:          cityData,
		SourceBreakdown:    sourceData,
		ReferrerCategories: refCatData,
		UTMSources:         utmSrcData,
		UTMMediums:         utmMedData,
		UTMCampaigns:       utmCamData,
	}, nil
}

func (s *dashboardService) GetGlobalAnalyticsComparison(ctx context.Context, userID uuid.UUID, from, to time.Time) (*response.GlobalAnalyticsComparisonResponse, *dto.ServiceError) {
	prevDuration := to.Sub(from)
	prevTo := from
	prevFrom := from.Add(-prevDuration)

	// Current period
	currentTotal, err := s.clickEventRepo.GetClicksInPeriod(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get current total clicks for comparison", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	currentUnique, err := s.clickEventRepo.GetUniqueClicksInPeriodByUserID(ctx, userID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get current unique clicks for comparison", zap.Error(err))
	}

	// Previous period
	prevTotal, err := s.clickEventRepo.GetClicksInPeriod(ctx, userID, prevFrom, prevTo)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get previous total clicks for comparison", zap.Error(err))
	}

	prevUnique, err := s.clickEventRepo.GetUniqueClicksInPeriodByUserID(ctx, userID, prevFrom, prevTo)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get previous unique clicks for comparison", zap.Error(err))
	}

	clicksTrend := calcTrend(prevTotal, currentTotal)
	uniqueTrend := calcTrend(prevUnique, currentUnique)

	return &response.GlobalAnalyticsComparisonResponse{
		Current: response.GlobalAnalyticsPeriodData{
			From:         from.Format(time.RFC3339),
			To:           to.Format(time.RFC3339),
			TotalClicks:  currentTotal,
			UniqueClicks: currentUnique,
		},
		Previous: response.GlobalAnalyticsPeriodData{
			From:         prevFrom.Format(time.RFC3339),
			To:           prevTo.Format(time.RFC3339),
			TotalClicks:  prevTotal,
			UniqueClicks: prevUnique,
		},
		Clicks:       clicksTrend,
		UniqueClicks: uniqueTrend,
	}, nil
}

// calcTrend computes the trend direction and percent change between two values.
func calcTrend(prev, current int64) response.TrendData {
	if prev == 0 {
		if current > 0 {
			return response.TrendData{Trend: "up", PercentChange: 100}
		}
		return response.TrendData{Trend: "flat", PercentChange: 0}
	}
	pct := float64(current-prev) / float64(prev) * 100
	trend := "flat"
	if pct > 0.5 {
		trend = "up"
	} else if pct < -0.5 {
		trend = "down"
	}
	if pct < 0 {
		pct = -pct
	}
	return response.TrendData{Trend: trend, PercentChange: pct}
}
