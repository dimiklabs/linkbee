package expiry

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/repository"
	emailSvc "github.com/shafikshaon/shortlink/service/email"
)

// ExpiryServiceI processes expiring link notifications.
type ExpiryServiceI interface {
	ProcessExpiringLinks(ctx context.Context)
}

type expiryService struct {
	linkRepo     repository.LinkRepositoryI
	userRepo     repository.UserRepositoryI
	emailSvc     emailSvc.EmailServiceI
	appCfg       *config.AppConfig
	emailCfg     *config.EmailConfig
	notifyBefore time.Duration
}

// NewExpiryService creates an expiry notification service that alerts users
// when their links are within notifyBefore (3 days) of expiring.
func NewExpiryService(
	linkRepo repository.LinkRepositoryI,
	userRepo repository.UserRepositoryI,
	emailSvc emailSvc.EmailServiceI,
	appCfg *config.AppConfig,
	emailCfg *config.EmailConfig,
) ExpiryServiceI {
	return &expiryService{
		linkRepo:     linkRepo,
		userRepo:     userRepo,
		emailSvc:     emailSvc,
		appCfg:       appCfg,
		emailCfg:     emailCfg,
		notifyBefore: 3 * 24 * time.Hour,
	}
}

func (s *expiryService) ProcessExpiringLinks(ctx context.Context) {
	links, err := s.linkRepo.GetLinksNearExpiry(ctx, s.notifyBefore, 100)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to fetch links near expiry", zap.Error(err))
		return
	}

	for _, link := range links {
		user, err := s.userRepo.GetByID(ctx, link.UserID)
		if err != nil {
			logger.WarnCtx(ctx, "Failed to get user for expiry notification",
				zap.String("link_id", link.ID.String()), zap.Error(err))
			continue
		}

		hours := time.Until(*link.ExpiresAt).Hours()
		daysLeft := int(hours/24) + 1
		if daysLeft < 1 {
			daysLeft = 1
		}

		userName := user.FirstName
		if userName == "" {
			userName = "there"
		}

		shortURL := fmt.Sprintf("%s/%s", s.appCfg.BaseDomain, link.Slug)
		title := link.Title
		if title == "" {
			title = link.Slug
		}

		body, err := renderExpiryEmail(expiryEmailData{
			AppName:      s.emailCfg.FromName,
			UserName:     userName,
			LinkTitle:    title,
			ShortURL:     shortURL,
			ExpiresAt:    link.ExpiresAt.Format("January 2, 2006 at 3:04 PM UTC"),
			DaysLeft:     daysLeft,
			DashboardURL: fmt.Sprintf("%s/dashboard/links", s.emailCfg.BaseURL),
		})
		if err != nil {
			logger.ErrorCtx(ctx, "Failed to render expiry email",
				zap.String("link_id", link.ID.String()), zap.Error(err))
			continue
		}

		dayWord := "days"
		if daysLeft == 1 {
			dayWord = "day"
		}
		subject := fmt.Sprintf("[%s] Your link \"%s\" expires in %d %s", s.emailCfg.FromName, title, daysLeft, dayWord)

		if err := s.emailSvc.SendHTML(ctx, user.Email, subject, body); err != nil {
			logger.WarnCtx(ctx, "Failed to send expiry notification email",
				zap.String("link_id", link.ID.String()),
				zap.String("email", user.Email),
				zap.Error(err))
			continue
		}

		if err := s.linkRepo.MarkExpiryNotified(ctx, link.ID); err != nil {
			logger.WarnCtx(ctx, "Failed to mark link expiry notified",
				zap.String("link_id", link.ID.String()), zap.Error(err))
		}

		logger.InfoCtx(ctx, "Link expiry notification sent",
			zap.String("link_id", link.ID.String()),
			zap.String("email", user.Email),
			zap.Int("days_left", daysLeft))
	}
}

type expiryEmailData struct {
	AppName      string
	UserName     string
	LinkTitle    string
	ShortURL     string
	ExpiresAt    string
	DaysLeft     int
	DashboardURL string
}

func renderExpiryEmail(data expiryEmailData) (string, error) {
	tmpl, err := template.New("expiry").Parse(expiryEmailTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

const expiryEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Link Expiring Soon</title>
</head>
<body style="margin:0;padding:0;font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;background-color:#f4f4f4;">
    <table role="presentation" style="width:100%;border-collapse:collapse;">
        <tr>
            <td align="center" style="padding:40px 0;">
                <table role="presentation" style="width:600px;max-width:100%;border-collapse:collapse;background-color:#ffffff;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding:40px 40px 30px;text-align:center;background-color:#f59e0b;border-radius:8px 8px 0 0;">
                            <h1 style="margin:0;color:#ffffff;font-size:28px;font-weight:600;">{{.AppName}}</h1>
                            <p style="margin:10px 0 0;color:#fffbeb;font-size:14px;">Link Expiry Notice</p>
                        </td>
                    </tr>

                    <!-- Content -->
                    <tr>
                        <td style="padding:40px;">
                            <h2 style="margin:0 0 20px;color:#1a1a1a;font-size:22px;font-weight:600;">
                                ⏰ Your link expires in {{.DaysLeft}} day{{if gt .DaysLeft 1}}s{{end}}
                            </h2>

                            <p style="margin:0 0 16px;color:#4a4a4a;font-size:16px;line-height:1.6;">
                                Hi {{.UserName}},
                            </p>

                            <p style="margin:0 0 24px;color:#4a4a4a;font-size:16px;line-height:1.6;">
                                Just a heads-up — one of your shortened links is about to expire.
                            </p>

                            <!-- Link card -->
                            <table role="presentation" style="width:100%;border-collapse:collapse;margin-bottom:24px;">
                                <tr>
                                    <td style="padding:20px;background-color:#fffbeb;border:1px solid #fde68a;border-radius:8px;">
                                        <p style="margin:0 0 6px;color:#92400e;font-size:12px;font-weight:600;text-transform:uppercase;letter-spacing:0.05em;">Link</p>
                                        <p style="margin:0 0 12px;color:#1a1a1a;font-size:16px;font-weight:600;">{{.LinkTitle}}</p>
                                        <p style="margin:0 0 8px;font-size:14px;">
                                            <a href="{{.ShortURL}}" style="color:#635bff;text-decoration:none;">{{.ShortURL}}</a>
                                        </p>
                                        <p style="margin:0;color:#6b7280;font-size:13px;">
                                            Expires: <strong style="color:#dc2626;">{{.ExpiresAt}}</strong>
                                        </p>
                                    </td>
                                </tr>
                            </table>

                            <p style="margin:0 0 24px;color:#4a4a4a;font-size:16px;line-height:1.6;">
                                After this date, the link will stop redirecting visitors. You can extend the expiry date from your dashboard.
                            </p>

                            <!-- CTA -->
                            <table role="presentation" style="width:100%;border-collapse:collapse;">
                                <tr>
                                    <td align="center" style="padding:8px 0 32px;">
                                        <a href="{{.DashboardURL}}" style="display:inline-block;padding:14px 36px;background-color:#635bff;color:#ffffff;text-decoration:none;font-size:15px;font-weight:600;border-radius:6px;">
                                            Go to Dashboard
                                        </a>
                                    </td>
                                </tr>
                            </table>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="padding:24px 40px;background-color:#f8f8f8;border-radius:0 0 8px 8px;text-align:center;">
                            <p style="margin:0 0 8px;color:#888888;font-size:13px;">
                                &copy; 2024 {{.AppName}}. All rights reserved.
                            </p>
                            <p style="margin:0;color:#aaaaaa;font-size:12px;">
                                This is an automated reminder. You can manage your links at any time from your dashboard.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
