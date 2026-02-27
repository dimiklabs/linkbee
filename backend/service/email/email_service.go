package email

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shafikshaon/shortlink/config"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
	"github.com/shafikshaon/shortlink/repository"
)

type EmailServiceI interface {
	SendVerificationEmail(ctx context.Context, user *model.User) error
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error
	// SendHTML sends a generic HTML email. Used by other services (e.g., reporting).
	SendHTML(ctx context.Context, to, subject, body string) error
}

type EmailService struct {
	config           *config.EmailConfig
	verificationRepo repository.EmailVerificationRepositoryI
	userRepo         repository.UserRepositoryI
}

func NewEmailService(
	cfg *config.EmailConfig,
	verificationRepo repository.EmailVerificationRepositoryI,
	userRepo repository.UserRepositoryI,
) EmailServiceI {
	return &EmailService{
		config:           cfg,
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
	}
}

func (s *EmailService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, user *model.User) error {
	logger.InfoCtx(ctx, "Sending verification email",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email))

	// Generate verification token
	token, err := s.generateToken()
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to generate verification token",
			zap.Error(err))
		return err
	}

	// Delete any existing verification tokens for this user
	if err := s.verificationRepo.DeleteByUserID(ctx, user.ID); err != nil {
		logger.WarnCtx(ctx, "Failed to delete existing verification tokens",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
	}

	// Create verification record
	verification := &model.EmailVerification{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(s.config.VerifyTokenTTL) * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}

	if err := s.verificationRepo.Create(ctx, verification); err != nil {
		logger.ErrorCtx(ctx, "Failed to create verification record",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return err
	}

	// Build verification URL
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.BaseURL, token)

	// Render email template
	emailBody, err := s.renderVerificationEmail(user, verificationURL)
	if err != nil {
		logger.ErrorCtx(ctx, "Failed to render verification email template",
			zap.Error(err))
		return err
	}

	// Send email
	if err := s.sendEmail(ctx, user.Email, "Verify your email address", emailBody); err != nil {
		logger.ErrorCtx(ctx, "Failed to send verification email",
			zap.String("email", user.Email),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Verification email sent successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
		zap.String("verification_url", verificationURL))

	return nil
}

func (s *EmailService) VerifyEmail(ctx context.Context, token string) error {
	logger.InfoCtx(ctx, "Verifying email token")

	// Get verification record
	verification, err := s.verificationRepo.GetByToken(ctx, token)
	if err != nil {
		logger.WarnCtx(ctx, "Invalid or expired verification token",
			zap.Error(err))
		return fmt.Errorf("invalid or expired verification token")
	}

	// Mark token as used
	if err := s.verificationRepo.MarkAsUsed(ctx, verification.ID); err != nil {
		logger.ErrorCtx(ctx, "Failed to mark verification token as used",
			zap.String("id", verification.ID.String()),
			zap.Error(err))
		return err
	}

	// Update user's email verification status
	if err := s.userRepo.UpdateFields(ctx, verification.UserID, map[string]interface{}{
		"email_verified":    true,
		"email_verified_at": new(time.Now()),
	}); err != nil {
		logger.ErrorCtx(ctx, "Failed to update user email verification status",
			zap.String("user_id", verification.UserID.String()),
			zap.Error(err))
		return err
	}

	logger.InfoCtx(ctx, "Email verified successfully",
		zap.String("user_id", verification.UserID.String()))

	return nil
}

func (s *EmailService) ResendVerificationEmail(ctx context.Context, userID uuid.UUID, email string) error {
	logger.InfoCtx(ctx, "Resending verification email",
		zap.String("user_id", userID.String()),
		zap.String("email", email))

	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logger.ErrorCtx(ctx, "User not found",
			zap.String("user_id", userID.String()),
			zap.Error(err))
		return fmt.Errorf("user not found")
	}

	// Check if email is already verified
	if user.EmailVerified {
		logger.WarnCtx(ctx, "Email already verified",
			zap.String("user_id", userID.String()))
		return fmt.Errorf("email already verified")
	}

	// Send new verification email
	return s.SendVerificationEmail(ctx, user)
}

// SendHTML sends a raw HTML email to the given address.
func (s *EmailService) SendHTML(ctx context.Context, to, subject, body string) error {
	return s.sendEmail(ctx, to, subject, body)
}

func (s *EmailService) renderVerificationEmail(user *model.User, verificationURL string) (string, error) {
	tmpl, err := template.New("verification").Parse(verificationEmailTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		UserName        string
		VerificationURL string
		ExpiryHours     int
		AppName         string
	}{
		UserName:        s.getUserDisplayName(user),
		VerificationURL: verificationURL,
		ExpiryHours:     s.config.VerifyTokenTTL,
		AppName:         s.config.FromName,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *EmailService) getUserDisplayName(user *model.User) string {
	if user.FirstName != "" {
		return user.FirstName
	}
	return "there"
}

func (s *EmailService) sendEmail(ctx context.Context, to, subject, body string) error {
	from := s.config.FromEmail
	password := s.config.SMTPPassword
	smtpHost := s.config.SMTPHost
	smtpPort := s.config.SMTPPort

	// Build email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message bytes.Buffer
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	// Send email
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	var auth smtp.Auth
	if s.config.SMTPUser != "" && password != "" {
		auth = smtp.PlainAuth("", s.config.SMTPUser, password, smtpHost)
	}

	if err := smtp.SendMail(addr, auth, from, []string{to}, message.Bytes()); err != nil {
		logger.ErrorCtx(ctx, "SMTP send failed",
			zap.String("to", to),
			zap.String("host", smtpHost),
			zap.Error(err))
		return err
	}

	return nil
}

const verificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; max-width: 100%; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background-color: #4F46E5; border-radius: 8px 8px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">{{.AppName}}</h1>
                        </td>
                    </tr>

                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="margin: 0 0 20px; color: #1a1a1a; font-size: 24px; font-weight: 600;">Verify Your Email Address</h2>

                            <p style="margin: 0 0 20px; color: #4a4a4a; font-size: 16px; line-height: 1.6;">
                                Hi {{.UserName}},
                            </p>

                            <p style="margin: 0 0 20px; color: #4a4a4a; font-size: 16px; line-height: 1.6;">
                                Thank you for signing up! Please verify your email address by clicking the button below:
                            </p>

                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 20px 0 30px;">
                                        <a href="{{.VerificationURL}}" style="display: inline-block; padding: 16px 40px; background-color: #4F46E5; color: #ffffff; text-decoration: none; font-size: 16px; font-weight: 600; border-radius: 6px; transition: background-color 0.2s;">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>

                            <p style="margin: 0 0 20px; color: #4a4a4a; font-size: 16px; line-height: 1.6;">
                                This link will expire in <strong>{{.ExpiryHours}} hours</strong>.
                            </p>

                            <p style="margin: 0 0 20px; color: #4a4a4a; font-size: 16px; line-height: 1.6;">
                                If you didn't create an account with us, you can safely ignore this email.
                            </p>

                            <!-- Alternative Link -->
                            <div style="margin-top: 30px; padding: 20px; background-color: #f8f8f8; border-radius: 6px;">
                                <p style="margin: 0 0 10px; color: #666666; font-size: 14px;">
                                    If the button doesn't work, copy and paste this link into your browser:
                                </p>
                                <p style="margin: 0; color: #4F46E5; font-size: 14px; word-break: break-all;">
                                    {{.VerificationURL}}
                                </p>
                            </div>
                        </td>
                    </tr>

                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f8f8f8; border-radius: 0 0 8px 8px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #888888; font-size: 14px;">
                                &copy; 2024 {{.AppName}}. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #888888; font-size: 12px;">
                                This is an automated message. Please do not reply to this email.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
