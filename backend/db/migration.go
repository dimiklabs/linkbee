package db

import (
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/model"
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
	)
	if err != nil {
		logger.Error("Auto migration failed", zap.Error(err))
		return err
	}

	logger.Info("Auto migration completed successfully")
	return nil
}
