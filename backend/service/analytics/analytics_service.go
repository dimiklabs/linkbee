package analytics

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/constant"
	"github.com/shafikshaon/shortlink/dto"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/repository"
	"github.com/shafikshaon/shortlink/response"
)

type AnalyticsServiceI interface {
	GetLinkAnalytics(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, from, to time.Time, granularity string) (*response.AnalyticsResponse, *dto.ServiceError)
	GetPeriodComparison(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, from, to time.Time) (*response.PeriodComparisonResponse, *dto.ServiceError)
	GetMultiLinkComparison(ctx context.Context, userID uuid.UUID, linkIDs []uuid.UUID, from, to time.Time) (*response.MultiLinkComparisonResponse, *dto.ServiceError)
}

type analyticsService struct {
	linkRepo       repository.LinkRepositoryI
	clickEventRepo repository.ClickEventRepositoryI
}

func NewAnalyticsService(linkRepo repository.LinkRepositoryI, clickEventRepo repository.ClickEventRepositoryI) AnalyticsServiceI {
	return &analyticsService{
		linkRepo:       linkRepo,
		clickEventRepo: clickEventRepo,
	}
}

func (s *analyticsService) GetLinkAnalytics(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, from, to time.Time, granularity string) (*response.AnalyticsResponse, *dto.ServiceError) {
	// Verify ownership
	link, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	if link.UserID != userID {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}

	totalClicks, err := s.clickEventRepo.GetClickCountByLinkID(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get click count", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	uniqueClicks, err := s.clickEventRepo.GetUniqueClickCountByLinkID(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get unique click count", zap.Error(err))
	}

	timeSeries, err := s.clickEventRepo.GetTimeSeriesData(ctx, linkID, from, to, granularity)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get time series data", zap.Error(err))
	}

	referrers, err := s.clickEventRepo.GetTopReferrers(ctx, linkID, 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get referrers", zap.Error(err))
	}

	devices, err := s.clickEventRepo.GetDeviceBreakdown(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get device breakdown", zap.Error(err))
	}

	countries, err := s.clickEventRepo.GetCountryBreakdown(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get country breakdown", zap.Error(err))
	}

	browsers, err := s.clickEventRepo.GetBrowserBreakdown(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get browser breakdown", zap.Error(err))
	}

	osBreakdown, err := s.clickEventRepo.GetOSBreakdown(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get OS breakdown", zap.Error(err))
	}

	sources, err := s.clickEventRepo.GetSourceBreakdown(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get source breakdown", zap.Error(err))
	}

	heatmap, err := s.clickEventRepo.GetHeatmapData(ctx, linkID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get heatmap data", zap.Error(err))
	}

	utmSources, err := s.clickEventRepo.GetUTMBreakdown(ctx, linkID, "utm_source", 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM source breakdown", zap.Error(err))
	}

	utmMediums, err := s.clickEventRepo.GetUTMBreakdown(ctx, linkID, "utm_medium", 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM medium breakdown", zap.Error(err))
	}

	utmCampaigns, err := s.clickEventRepo.GetUTMBreakdown(ctx, linkID, "utm_campaign", 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM campaign breakdown", zap.Error(err))
	}

	utmContents, err := s.clickEventRepo.GetUTMBreakdown(ctx, linkID, "utm_content", 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM content breakdown", zap.Error(err))
	}

	utmTerms, err := s.clickEventRepo.GetUTMBreakdown(ctx, linkID, "utm_term", 10)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get UTM term breakdown", zap.Error(err))
	}

	cities, err := s.clickEventRepo.GetCityBreakdown(ctx, linkID, 20)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get city breakdown", zap.Error(err))
	}

	referrerCategories, err := s.clickEventRepo.GetReferrerCategoryBreakdown(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get referrer category breakdown", zap.Error(err))
	}

	firstTimeVisitors, returningVisitors, err := s.clickEventRepo.GetReturnVisitorStats(ctx, linkID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get return visitor stats", zap.Error(err))
	}

	tsData := make([]response.TimeSeriesData, len(timeSeries))
	for i, ts := range timeSeries {
		tsData[i] = response.TimeSeriesData{Timestamp: ts.Timestamp, Count: ts.Count}
	}

	refData := make([]response.ReferrerData, len(referrers))
	for i, r := range referrers {
		refData[i] = response.ReferrerData{Referrer: r.Referrer, Count: r.Count}
	}

	devData := make([]response.DeviceData, len(devices))
	for i, d := range devices {
		devData[i] = response.DeviceData{DeviceType: d.DeviceType, Count: d.Count}
	}

	countryData := make([]response.CountryData, len(countries))
	for i, c := range countries {
		countryData[i] = response.CountryData{Country: c.Country, Count: c.Count}
	}

	browserData := make([]response.BrowserData, len(browsers))
	for i, b := range browsers {
		browserData[i] = response.BrowserData{Browser: b.Browser, Count: b.Count}
	}

	osData := make([]response.OSData, len(osBreakdown))
	for i, o := range osBreakdown {
		osData[i] = response.OSData{OS: o.OS, Count: o.Count}
	}

	sourceData := make([]response.SourceData, len(sources))
	for i, s := range sources {
		sourceData[i] = response.SourceData{Source: s.Source, Count: s.Count}
	}

	heatmapData := make([]response.HeatmapData, len(heatmap))
	for i, h := range heatmap {
		heatmapData[i] = response.HeatmapData{DayOfWeek: h.DayOfWeek, Hour: h.Hour, Count: h.Count}
	}

	utmSourceData := make([]response.UTMData, len(utmSources))
	for i, u := range utmSources {
		utmSourceData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	utmMediumData := make([]response.UTMData, len(utmMediums))
	for i, u := range utmMediums {
		utmMediumData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	utmCampaignData := make([]response.UTMData, len(utmCampaigns))
	for i, u := range utmCampaigns {
		utmCampaignData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	utmContentData := make([]response.UTMData, len(utmContents))
	for i, u := range utmContents {
		utmContentData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	utmTermData := make([]response.UTMData, len(utmTerms))
	for i, u := range utmTerms {
		utmTermData[i] = response.UTMData{Value: u.Value, Count: u.Count}
	}

	cityData := make([]response.CityData, len(cities))
	for i, c := range cities {
		cityData[i] = response.CityData{City: c.City, Country: c.Country, Count: c.Count}
	}

	refCatData := make([]response.ReferrerCategoryData, len(referrerCategories))
	for i, rc := range referrerCategories {
		refCatData[i] = response.ReferrerCategoryData{Category: rc.Category, Count: rc.Count}
	}

	return &response.AnalyticsResponse{
		LinkID:            linkID,
		TotalClicks:       totalClicks,
		UniqueClicks:      uniqueClicks,
		FirstTimeVisitors: firstTimeVisitors,
		ReturningVisitors: returningVisitors,
		TimeSeries:   tsData,
		Referrers:    refData,
		Devices:      devData,
		Countries:    countryData,
		Browsers:     browserData,
		OSBreakdown:  osData,
		Sources:      sourceData,
		Heatmap:      heatmapData,
		UTMSources:   utmSourceData,
		UTMMediums:   utmMediumData,
		UTMCampaigns: utmCampaignData,
		UTMContents:        utmContentData,
		UTMTerms:           utmTermData,
		Cities:             cityData,
		ReferrerCategories: refCatData,
	}, nil
}

func (s *analyticsService) GetPeriodComparison(ctx context.Context, linkID uuid.UUID, userID uuid.UUID, from, to time.Time) (*response.PeriodComparisonResponse, *dto.ServiceError) {
	// Verify ownership
	link, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	if link.UserID != userID {
		return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
	}

	// Current period metrics
	curClicks, err := s.clickEventRepo.GetClicksInPeriodByLinkID(ctx, linkID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get current period clicks", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	curUnique, err := s.clickEventRepo.GetUniqueClicksInPeriodByLinkID(ctx, linkID, from, to)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get current period unique clicks", zap.Error(err))
	}

	// Previous period: same duration shifted back
	duration := to.Sub(from)
	prevFrom := from.Add(-duration)
	prevTo := from

	prevClicks, err := s.clickEventRepo.GetClicksInPeriodByLinkID(ctx, linkID, prevFrom, prevTo)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get previous period clicks", zap.Error(err))
	}
	prevUnique, err := s.clickEventRepo.GetUniqueClicksInPeriodByLinkID(ctx, linkID, prevFrom, prevTo)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to get previous period unique clicks", zap.Error(err))
	}

	return &response.PeriodComparisonResponse{
		LinkID: linkID,
		Current: response.PeriodMetrics{
			From:         from.Format(time.RFC3339),
			To:           to.Format(time.RFC3339),
			TotalClicks:  curClicks,
			UniqueClicks: curUnique,
		},
		Previous: response.PeriodMetrics{
			From:         prevFrom.Format(time.RFC3339),
			To:           prevTo.Format(time.RFC3339),
			TotalClicks:  prevClicks,
			UniqueClicks: prevUnique,
		},
		Clicks:       periodChange(curClicks, prevClicks),
		UniqueClicks: periodChange(curUnique, prevUnique),
	}, nil
}

// periodChange computes the delta between current and previous values.
func periodChange(current, previous int64) response.PeriodChange {
	diff := current - previous
	var pct float64
	if previous > 0 {
		pct = math.Round((float64(diff)/float64(previous))*1000) / 10 // 1 decimal
	} else if current > 0 {
		pct = 100.0
	}

	trend := "stable"
	if pct > 5.0 {
		trend = "up"
	} else if pct < -5.0 {
		trend = "down"
	}

	return response.PeriodChange{
		CountChange:   diff,
		PercentChange: pct,
		Trend:         trend,
	}
}

func (s *analyticsService) GetMultiLinkComparison(ctx context.Context, userID uuid.UUID, linkIDs []uuid.UUID, from, to time.Time) (*response.MultiLinkComparisonResponse, *dto.ServiceError) {
	if len(linkIDs) < 2 || len(linkIDs) > 5 {
		return nil, dto.NewBadRequestError(constant.ErrCodeBadRequest, "provide between 2 and 5 link IDs")
	}

	spanDays := int64(math.Round(to.Sub(from).Hours() / 24))
	if spanDays < 1 {
		spanDays = 1
	}

	metrics := make([]response.LinkComparisonMetric, 0, len(linkIDs))

	for _, linkID := range linkIDs {
		link, err := s.linkRepo.GetByID(ctx, linkID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
			}
			return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
		}
		if link.UserID != userID {
			return nil, dto.NewNotFoundError(constant.ErrCodeLinkNotFound, constant.ErrMsgLinkNotFound)
		}

		clicks, _ := s.clickEventRepo.GetClicksInPeriodByLinkID(ctx, linkID, from, to)
		unique, _ := s.clickEventRepo.GetUniqueClicksInPeriodByLinkID(ctx, linkID, from, to)

		clicksPerDay := math.Round(float64(clicks)/float64(spanDays)*100) / 100

		topReferrer := ""
		if refs, _ := s.clickEventRepo.GetTopReferrers(ctx, linkID, 1); len(refs) > 0 {
			topReferrer = refs[0].Referrer
		}

		topCountry := ""
		if countries, _ := s.clickEventRepo.GetCountryBreakdown(ctx, linkID); len(countries) > 0 {
			topCountry = countries[0].Country
		}

		topBrowser := ""
		if browsers, _ := s.clickEventRepo.GetBrowserBreakdown(ctx, linkID); len(browsers) > 0 {
			topBrowser = browsers[0].Browser
		}

		topDevice := ""
		if devices, _ := s.clickEventRepo.GetDeviceBreakdown(ctx, linkID); len(devices) > 0 {
			topDevice = devices[0].DeviceType
		}

		metrics = append(metrics, response.LinkComparisonMetric{
			LinkID:       linkID,
			Slug:         link.Slug,
			Title:        link.Title,
			TotalClicks:  clicks,
			UniqueClicks: unique,
			ClicksPerDay: clicksPerDay,
			TopReferrer:  topReferrer,
			TopCountry:   topCountry,
			TopBrowser:   topBrowser,
			TopDevice:    topDevice,
		})
	}

	return &response.MultiLinkComparisonResponse{
		Links:    metrics,
		From:     from.Format(time.RFC3339),
		To:       to.Format(time.RFC3339),
		SpanDays: spanDays,
	}, nil
}
