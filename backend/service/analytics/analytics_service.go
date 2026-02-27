package analytics

import (
	"context"
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

	return &response.AnalyticsResponse{
		LinkID:       linkID,
		TotalClicks:  totalClicks,
		UniqueClicks: uniqueClicks,
		TimeSeries:   tsData,
		Referrers:    refData,
		Devices:      devData,
		Countries:    countryData,
		Browsers:     browserData,
		OSBreakdown:  osData,
		Heatmap:      heatmapData,
		UTMSources:   utmSourceData,
		UTMMediums:   utmMediumData,
		UTMCampaigns: utmCampaignData,
	}, nil
}
