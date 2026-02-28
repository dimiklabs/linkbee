package reporting

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/constant"
	"github.com/shafikshaon/linkbee/dto"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/request"
	emailSvc "github.com/shafikshaon/linkbee/service/email"
)

type ReportingServiceI interface {
	CreateReport(ctx context.Context, userID uuid.UUID, req request.CreateReportRequest) (*model.AnalyticsReport, *dto.ServiceError)
	UpdateReport(ctx context.Context, userID, id uuid.UUID, req request.UpdateReportRequest) (*model.AnalyticsReport, *dto.ServiceError)
	DeleteReport(ctx context.Context, userID, id uuid.UUID) *dto.ServiceError
	GetReport(ctx context.Context, userID, id uuid.UUID) (*model.AnalyticsReport, *dto.ServiceError)
	ListReports(ctx context.Context, userID uuid.UUID) ([]*model.AnalyticsReport, *dto.ServiceError)
	SendReportNow(ctx context.Context, userID, id uuid.UUID) *dto.ServiceError
	GetDeliveries(ctx context.Context, userID, id uuid.UUID) ([]*model.ReportDelivery, *dto.ServiceError)
	// ProcessDueReports is called by the background worker.
	ProcessDueReports(ctx context.Context)
}

type reportingService struct {
	reportRepo     repository.AnalyticsReportRepositoryI
	userRepo       repository.UserRepositoryI
	clickEventRepo repository.ClickEventRepositoryI
	linkRepo       repository.LinkRepositoryI
	emailSvc       emailSvc.EmailServiceI
	appCfg         *config.AppConfig
	emailCfg       *config.EmailConfig
}

func NewReportingService(
	reportRepo repository.AnalyticsReportRepositoryI,
	userRepo repository.UserRepositoryI,
	clickEventRepo repository.ClickEventRepositoryI,
	linkRepo repository.LinkRepositoryI,
	emailSvc emailSvc.EmailServiceI,
	appCfg *config.AppConfig,
	emailCfg *config.EmailConfig,
) ReportingServiceI {
	return &reportingService{
		reportRepo:     reportRepo,
		userRepo:       userRepo,
		clickEventRepo: clickEventRepo,
		linkRepo:       linkRepo,
		emailSvc:       emailSvc,
		appCfg:         appCfg,
		emailCfg:       emailCfg,
	}
}

// ── CRUD ──────────────────────────────────────────────────────────────────────

func (s *reportingService) CreateReport(ctx context.Context, userID uuid.UUID, req request.CreateReportRequest) (*model.AnalyticsReport, *dto.ServiceError) {
	t := nextRunAt(req.Frequency, time.Now().UTC())
	report := &model.AnalyticsReport{
		UserID:    userID,
		Name:      req.Name,
		LinkIDs:   pq.StringArray(req.LinkIDs),
		Frequency: req.Frequency,
		NextRunAt: &t,
		IsActive:  true,
	}
	if err := s.reportRepo.Create(ctx, report); err != nil {
		logger.ErrorCtx(ctx, "Failed to create report", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return report, nil
}

func (s *reportingService) UpdateReport(ctx context.Context, userID, id uuid.UUID, req request.UpdateReportRequest) (*model.AnalyticsReport, *dto.ServiceError) {
	report, err := s.reportRepo.GetByID(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError("REPORT_NOT_FOUND", "report not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	if req.Name != "" {
		report.Name = req.Name
	}
	if len(req.LinkIDs) > 0 {
		report.LinkIDs = pq.StringArray(req.LinkIDs)
	}
	if req.Frequency != "" {
		report.Frequency = req.Frequency
		t := nextRunAt(req.Frequency, time.Now().UTC())
		report.NextRunAt = &t
	}
	if req.IsActive != nil {
		report.IsActive = *req.IsActive
	}

	if err := s.reportRepo.Update(ctx, report); err != nil {
		logger.ErrorCtx(ctx, "Failed to update report", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return report, nil
}

func (s *reportingService) DeleteReport(ctx context.Context, userID, id uuid.UUID) *dto.ServiceError {
	if err := s.reportRepo.Delete(ctx, id, userID); err != nil {
		logger.ErrorCtx(ctx, "Failed to delete report", zap.Error(err))
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return nil
}

func (s *reportingService) GetReport(ctx context.Context, userID, id uuid.UUID) (*model.AnalyticsReport, *dto.ServiceError) {
	report, err := s.reportRepo.GetByID(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError("REPORT_NOT_FOUND", "report not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return report, nil
}

func (s *reportingService) ListReports(ctx context.Context, userID uuid.UUID) ([]*model.AnalyticsReport, *dto.ServiceError) {
	reports, err := s.reportRepo.ListByUserID(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to list reports", zap.Error(err))
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return reports, nil
}

func (s *reportingService) GetDeliveries(ctx context.Context, userID, id uuid.UUID) ([]*model.ReportDelivery, *dto.ServiceError) {
	// Verify ownership first
	if _, err := s.reportRepo.GetByID(ctx, id, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFoundError("REPORT_NOT_FOUND", "report not found")
		}
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	deliveries, err := s.reportRepo.ListDeliveries(ctx, id, 20)
	if err != nil {
		return nil, dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	return deliveries, nil
}

// ── Delivery ──────────────────────────────────────────────────────────────────

func (s *reportingService) SendReportNow(ctx context.Context, userID, id uuid.UUID) *dto.ServiceError {
	report, err := s.reportRepo.GetByID(ctx, id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFoundError("REPORT_NOT_FOUND", "report not found")
		}
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, constant.ErrMsgInternalServer)
	}

	to := time.Now().UTC()
	from := periodFrom(report.Frequency, to)
	if sendErr := s.deliver(ctx, report, user.Email, from, to); sendErr != nil {
		return dto.NewInternalError(constant.ErrCodeInternalServer, "failed to send report")
	}
	return nil
}

func (s *reportingService) ProcessDueReports(ctx context.Context) {
	now := time.Now().UTC()
	reports, err := s.reportRepo.GetDueReports(ctx, now)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch due reports", zap.Error(err))
		return
	}

	for _, report := range reports {
		user, uErr := s.userRepo.GetByID(ctx, report.UserID)
		if uErr != nil {
			logger.ErrorCtx(ctx, "Failed to get user for report", zap.String("report_id", report.ID.String()), zap.Error(uErr))
			continue
		}

		to := now
		from := periodFrom(report.Frequency, to)

		if sendErr := s.deliver(ctx, report, user.Email, from, to); sendErr != nil {
			logger.ErrorCtx(ctx, "Failed to deliver report", zap.String("report_id", report.ID.String()), zap.Error(sendErr))
		}

		// Advance next_run_at regardless of delivery success
		t := nextRunAt(report.Frequency, now)
		report.NextRunAt = &t
		if updateErr := s.reportRepo.Update(ctx, report); updateErr != nil {
			logger.ErrorCtx(ctx, "Failed to update report next_run_at", zap.String("report_id", report.ID.String()), zap.Error(updateErr))
		}
	}
}

// linkStat holds per-link click data collected during report generation.
type linkStat struct {
	Title        string
	Slug         string
	ShortURL     string
	TotalClicks  int64
	UniqueClicks int64
}

// deliver collects stats, renders email HTML and records the delivery.
func (s *reportingService) deliver(ctx context.Context, report *model.AnalyticsReport, toEmail string, from, to time.Time) error {
	stats := make([]linkStat, 0, len(report.LinkIDs))
	var grandTotal, grandUnique int64

	for _, idStr := range report.LinkIDs {
		linkID, parseErr := uuid.Parse(idStr)
		if parseErr != nil {
			continue
		}
		link, linkErr := s.linkRepo.GetByID(ctx, linkID)
		if linkErr != nil {
			continue
		}
		total, _ := s.clickEventRepo.GetClicksInPeriodByLinkID(ctx, linkID, from, to)
		unique, _ := s.clickEventRepo.GetUniqueClicksInPeriodByLinkID(ctx, linkID, from, to)
		grandTotal += total
		grandUnique += unique
		stats = append(stats, linkStat{
			Title:        link.Title,
			Slug:         link.Slug,
			ShortURL:     fmt.Sprintf("%s/%s", s.appCfg.BaseDomain, link.Slug),
			TotalClicks:  total,
			UniqueClicks: unique,
		})
	}

	body, renderErr := renderReportEmail(report.Name, report.Frequency, from, to, stats, grandTotal, grandUnique, s.emailCfg.FromName)
	if renderErr != nil {
		s.recordDelivery(ctx, report.ID, "failed", renderErr.Error())
		return renderErr
	}

	subject := fmt.Sprintf("[%s] %s – %s analytics report", s.emailCfg.FromName, report.Name, capitalize(report.Frequency))
	if sendErr := s.emailSvc.SendHTML(ctx, toEmail, subject, body); sendErr != nil {
		s.recordDelivery(ctx, report.ID, "failed", sendErr.Error())
		return sendErr
	}

	s.recordDelivery(ctx, report.ID, "sent", "")
	return nil
}

func (s *reportingService) recordDelivery(ctx context.Context, reportID uuid.UUID, status, reason string) {
	now := time.Now().UTC()
	d := &model.ReportDelivery{
		ReportID:      reportID,
		Status:        status,
		FailureReason: reason,
		CreatedAt:     now,
	}
	if status == "sent" {
		d.DeliveredAt = &now
	}
	if err := s.reportRepo.CreateDelivery(ctx, d); err != nil {
		logger.ErrorCtx(ctx, "Failed to record delivery", zap.Error(err))
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func nextRunAt(frequency string, from time.Time) time.Time {
	switch frequency {
	case "daily":
		return from.Add(24 * time.Hour)
	case "weekly":
		return from.AddDate(0, 0, 7)
	case "monthly":
		return from.AddDate(0, 1, 0)
	default:
		return from.AddDate(0, 0, 7)
	}
}

func periodFrom(frequency string, to time.Time) time.Time {
	switch frequency {
	case "daily":
		return to.AddDate(0, 0, -1)
	case "weekly":
		return to.AddDate(0, 0, -7)
	case "monthly":
		return to.AddDate(0, -1, 0)
	default:
		return to.AddDate(0, 0, -7)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}

// ── Email template ─────────────────────────────────────────────────────────────

type reportEmailData struct {
	AppName     string
	ReportName  string
	Frequency   string
	PeriodFrom  string
	PeriodTo    string
	GrandTotal  int64
	GrandUnique int64
	Links       []reportLinkRow
}

type reportLinkRow struct {
	Title        string
	Slug         string
	ShortURL     string
	TotalClicks  int64
	UniqueClicks int64
}

func renderReportEmail(
	name, frequency string,
	from, to time.Time,
	stats []linkStat,
	grandTotal, grandUnique int64,
	appName string,
) (string, error) {
	rows := make([]reportLinkRow, len(stats))
	for i, s := range stats {
		title := s.Title
		if title == "" {
			title = s.Slug
		}
		rows[i] = reportLinkRow{
			Title:        title,
			Slug:         s.Slug,
			ShortURL:     s.ShortURL,
			TotalClicks:  s.TotalClicks,
			UniqueClicks: s.UniqueClicks,
		}
	}

	data := reportEmailData{
		AppName:     appName,
		ReportName:  name,
		Frequency:   capitalize(frequency),
		PeriodFrom:  from.Format("Jan 2, 2006"),
		PeriodTo:    to.Format("Jan 2, 2006"),
		GrandTotal:  grandTotal,
		GrandUnique: grandUnique,
		Links:       rows,
	}

	tmpl, err := template.New("report").Parse(reportEmailTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const reportEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.ReportName}}</title>
</head>
<body style="margin:0;padding:0;font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;background:#f4f4f4;">
  <table role="presentation" style="width:100%;border-collapse:collapse;">
    <tr>
      <td align="center" style="padding:40px 0;">
        <table role="presentation" style="width:600px;max-width:100%;border-collapse:collapse;background:#fff;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.1);">
          <!-- Header -->
          <tr>
            <td style="padding:32px 40px 24px;background:#635bff;border-radius:8px 8px 0 0;">
              <h1 style="margin:0;color:#fff;font-size:24px;font-weight:700;">{{.AppName}}</h1>
              <p style="margin:8px 0 0;color:rgba(255,255,255,.8);font-size:14px;">{{.Frequency}} Analytics Report</p>
            </td>
          </tr>
          <!-- Period -->
          <tr>
            <td style="padding:24px 40px 0;">
              <p style="margin:0;color:#697386;font-size:14px;">
                Period: <strong style="color:#1a1f36;">{{.PeriodFrom}} – {{.PeriodTo}}</strong>
              </p>
              <h2 style="margin:8px 0 0;color:#1a1f36;font-size:20px;font-weight:700;">{{.ReportName}}</h2>
            </td>
          </tr>
          <!-- Summary cards -->
          <tr>
            <td style="padding:20px 40px;">
              <table role="presentation" style="width:100%;border-collapse:collapse;">
                <tr>
                  <td style="width:50%;padding-right:8px;">
                    <div style="background:#f7f9fc;border-radius:8px;padding:16px 20px;text-align:center;">
                      <div style="font-size:28px;font-weight:700;color:#635bff;">{{.GrandTotal}}</div>
                      <div style="font-size:13px;color:#697386;margin-top:4px;">Total Clicks</div>
                    </div>
                  </td>
                  <td style="width:50%;padding-left:8px;">
                    <div style="background:#f7f9fc;border-radius:8px;padding:16px 20px;text-align:center;">
                      <div style="font-size:28px;font-weight:700;color:#22c55e;">{{.GrandUnique}}</div>
                      <div style="font-size:13px;color:#697386;margin-top:4px;">Unique Clicks</div>
                    </div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <!-- Links table -->
          <tr>
            <td style="padding:0 40px 32px;">
              <h3 style="margin:0 0 12px;color:#1a1f36;font-size:15px;font-weight:600;">Link Breakdown</h3>
              <table role="presentation" style="width:100%;border-collapse:collapse;font-size:13px;">
                <thead>
                  <tr style="background:#f7f9fc;">
                    <th style="text-align:left;padding:8px 12px;color:#697386;font-weight:600;border-radius:6px 0 0 6px;">Link</th>
                    <th style="text-align:right;padding:8px 12px;color:#697386;font-weight:600;">Total</th>
                    <th style="text-align:right;padding:8px 12px;color:#697386;font-weight:600;border-radius:0 6px 6px 0;">Unique</th>
                  </tr>
                </thead>
                <tbody>
                  {{range .Links}}
                  <tr style="border-bottom:1px solid #f0f0f0;">
                    <td style="padding:10px 12px;">
                      <div style="font-weight:600;color:#1a1f36;">{{.Title}}</div>
                      <div style="color:#635bff;font-size:12px;margin-top:2px;">{{.ShortURL}}</div>
                    </td>
                    <td style="text-align:right;padding:10px 12px;font-weight:600;color:#1a1f36;">{{.TotalClicks}}</td>
                    <td style="text-align:right;padding:10px 12px;color:#697386;">{{.UniqueClicks}}</td>
                  </tr>
                  {{end}}
                </tbody>
              </table>
            </td>
          </tr>
          <!-- Footer -->
          <tr>
            <td style="padding:24px 40px;background:#f8f8f8;border-radius:0 0 8px 8px;text-align:center;">
              <p style="margin:0;color:#888;font-size:13px;">
                &copy; 2026 {{.AppName}}. This is an automated report. Do not reply.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`
