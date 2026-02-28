package cmd

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/shafikshaon/linkbee/config"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/server"
	"github.com/shafikshaon/linkbee/util"
)

func serveRest(ctx context.Context) {
	cfg := config.NewConfig(ctx)

	defer logger.Sync()

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Running in production mode")
	} else {
		logger.Info("Running in development mode")
	}

	srv := server.NewServer(ctx, cfg)

	app := gin.New()

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("password_complexity", util.PasswordComplexityValidator); err != nil {
			logger.Fatal("Failed to register password complexity validator", zap.Error(err))
		}
		logger.Debug("Custom validators registered successfully")
	}

	srv.ConfigureRoutes(ctx, app)

	logger.Info("Starting HTTP server",
		zap.String("address", cfg.App.ServerAddress),
		zap.String("env", cfg.App.Env))

	if err := app.Run(cfg.App.ServerAddress); err != nil {
		logger.Fatal("Failed to start server",
			zap.Error(err),
			zap.String("address", cfg.App.ServerAddress))
	}
}
