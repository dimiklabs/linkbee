package db

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/logger"
)

func NewPostgresDB(_ context.Context, cfg *config.DBConfig, isReplica bool) (*gorm.DB, error) {
	var dsn string
	var connType string

	if isReplica {
		connType = "replica"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.ReadHost, cfg.ReadPort, cfg.ReadUser, cfg.ReadPassword, cfg.Name)
	} else {
		connType = "master"
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	}

	logger.Debug("Attempting to connect to PostgreSQL",
		zap.String("type", connType),
		zap.String("database", cfg.Name))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL",
			zap.String("type", connType),
			zap.Error(err))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnection)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnection)

	logger.Debug("PostgreSQL connection pool configured",
		zap.String("type", connType),
		zap.Int("max_open_connections", cfg.MaxOpenConnection),
		zap.Int("max_idle_connections", cfg.MaxIdleConnection))

	return db, nil
}
