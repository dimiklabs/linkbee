package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/model"
)

func RunAutoMigration(db *gorm.DB) error {
	logger.Info("Running auto migration")

	err := db.AutoMigrate(
		&model.User{},
		&model.PasswordReset{},
		&model.Session{},
		&model.EmailVerification{},
		&model.APIKey{},
		&model.Folder{},
		&model.Link{},
		&model.LinkVariant{},
		&model.LinkGeoRule{},
		&model.ClickEvent{},
		&model.Webhook{},
		&model.RetargetingPixel{},
		&model.BioPage{},
		&model.BioLink{},
		&model.Subscription{},
		&model.CustomDomain{},
		&model.AuditLog{},
		&model.TotpBackupCode{},
		&model.WebhookDelivery{},
		&model.AnalyticsReport{},
		&model.ReportDelivery{},
		&model.Team{},
		&model.TeamMember{},
	)
	if err != nil {
		logger.Error("Auto migration failed", zap.Error(err))
		return err
	}

	logger.Info("Auto migration completed successfully")
	return nil
}
