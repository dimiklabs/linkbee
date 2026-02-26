package server

import (
	"context"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/shafikshaon/shortlink/docs" // swagger generated docs
	"github.com/shafikshaon/shortlink/handler"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/repository"
	analyticsSvc "github.com/shafikshaon/shortlink/service/analytics"
	apiKeySvc "github.com/shafikshaon/shortlink/service/apikey"
	authSrv "github.com/shafikshaon/shortlink/service/auth"
	clickSvc "github.com/shafikshaon/shortlink/service/click"
	demoSvc "github.com/shafikshaon/shortlink/service/demo"
	emailSrv "github.com/shafikshaon/shortlink/service/email"
	facebookSrv "github.com/shafikshaon/shortlink/service/facebook"
	folderSvc "github.com/shafikshaon/shortlink/service/folder"
	geoSvc "github.com/shafikshaon/shortlink/service/geo"
	georoutingSvc "github.com/shafikshaon/shortlink/service/georouting"
	githubSrv "github.com/shafikshaon/shortlink/service/github"
	googleSrv "github.com/shafikshaon/shortlink/service/google"
	healthSrv "github.com/shafikshaon/shortlink/service/health"
	linkSvc "github.com/shafikshaon/shortlink/service/link"
	qrSvc "github.com/shafikshaon/shortlink/service/qr"
	rateLimitSrv "github.com/shafikshaon/shortlink/service/ratelimit"
	redirectSvc "github.com/shafikshaon/shortlink/service/redirect"
	splitSvc "github.com/shafikshaon/shortlink/service/split"
	userSrv "github.com/shafikshaon/shortlink/service/user"
	bioSvc     "github.com/shafikshaon/shortlink/service/bio"
	pixelSvc   "github.com/shafikshaon/shortlink/service/pixel"
	previewSvc "github.com/shafikshaon/shortlink/service/preview"
	webhookSvc "github.com/shafikshaon/shortlink/service/webhook"
	"github.com/shafikshaon/shortlink/worker"
)

func (s *Server) ConfigureRoutes(ctx context.Context, router *gin.Engine) {
	logger.Debug("Configuring routes")

	s.setupMiddleware(router)

	// ── Repositories ─────────────────────────────────────────────────────────
	bioRepo               := repository.NewBioRepository(s.MasterDB, s.ReplicaDB)
	variantRepo           := repository.NewLinkVariantRepository(s.MasterDB, s.ReplicaDB)
	geoRuleRepo           := repository.NewLinkGeoRuleRepository(s.MasterDB, s.ReplicaDB)
	apiKeyRepo            := repository.NewAPIKeyRepository(s.MasterDB, s.ReplicaDB)
	webhookRepo           := repository.NewWebhookRepository(s.MasterDB, s.ReplicaDB)
	pixelRepo             := repository.NewRetargetingPixelRepository(s.MasterDB, s.ReplicaDB)
	folderRepo            := repository.NewFolderRepository(s.MasterDB, s.ReplicaDB)
	userRepo              := repository.NewUserRepository(s.MasterDB, s.ReplicaDB)
	passwordResetRepo     := repository.NewPasswordResetRepository(s.MasterDB, s.ReplicaDB)
	tokenBlacklistRepo    := repository.NewTokenBlacklistRepository(s.Cache)
	sessionRepo           := repository.NewSessionRepository(s.MasterDB, s.ReplicaDB)
	emailVerificationRepo := repository.NewEmailVerificationRepository(s.MasterDB, s.ReplicaDB)
	oauthStateRepo        := repository.NewOAuthStateRepository(s.Cache)
	rateLimitRepo         := repository.NewRateLimitRepository(s.Cache)
	linkRepo              := repository.NewLinkRepository(s.MasterDB, s.ReplicaDB)
	clickEventRepo        := repository.NewClickEventRepository(s.MasterDB, s.ReplicaDB)

	// ── Services ──────────────────────────────────────────────────────────────
	healthService        := healthSrv.NewHealthService(s.MasterDB, s.ReplicaDB, s.Cache, s.Cfg.App.Env)
	userService          := userSrv.NewUserService(userRepo)
	emailService         := emailSrv.NewEmailService(s.Cfg.Email, emailVerificationRepo, userRepo)
	authService          := authSrv.NewAuthService(userService, passwordResetRepo, tokenBlacklistRepo, sessionRepo, emailService, s.Cfg.App, s.Cfg.Session)
	rateLimitService     := rateLimitSrv.NewRateLimitService(rateLimitRepo, s.Cfg.RateLimit)
	googleOAuthService   := googleSrv.NewGoogleOAuthService(s.Cfg.Google, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)
	githubOAuthService   := githubSrv.NewGitHubOAuthService(s.Cfg.GitHub, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)
	facebookOAuthService := facebookSrv.NewFacebookOAuthService(s.Cfg.Facebook, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)

	bioService         := bioSvc.NewBioService(bioRepo)
	previewService     := previewSvc.NewPreviewService(s.Cache)
	folderService      := folderSvc.NewFolderService(folderRepo)
	apiKeyService      := apiKeySvc.NewAPIKeyService(apiKeyRepo)
	webhookService     := webhookSvc.NewWebhookService(webhookRepo)
	pixelService       := pixelSvc.NewPixelService(pixelRepo, linkRepo)
	splitService       := splitSvc.NewSplitService(variantRepo, linkRepo)
	linkService        := linkSvc.NewLinkService(linkRepo, s.Cfg.App, s.Cfg.Link)
	geoService         := geoSvc.NewGeoService(s.Cfg.App.GeoDBPath)
	geoRoutingService  := georoutingSvc.NewGeoRoutingService(geoRuleRepo, linkRepo)
	redirectService    := redirectSvc.NewRedirectService(linkRepo, variantRepo, geoRuleRepo, s.Cache, s.Cfg.Link.CacheTTLSeconds)
	clickService       := clickSvc.NewClickService(s.Cache)
	qrService          := qrSvc.NewQRService()
	analyticsService   := analyticsSvc.NewAnalyticsService(linkRepo, clickEventRepo)
	demoService        := demoSvc.NewDemoService(linkRepo, s.Cache, s.Cfg.App, s.Cfg.Link)

	// ── Click worker (background goroutine) ───────────────────────────────────
	clickWorker := worker.NewClickWorker(
		s.Cache,
		clickEventRepo,
		s.Cfg.Link.ClickQueueBatchSize,
		s.Cfg.Link.ClickQueueFlushInterval,
	)
	go clickWorker.Start(ctx)

	// ── Health worker (background goroutine) ──────────────────────────────────
	healthWorker := worker.NewHealthWorker(linkRepo)
	go healthWorker.Start(ctx)

	// ── Handlers ──────────────────────────────────────────────────────────────
	healthHandler    := handler.NewHealthHandler(healthService)
	authHandler      := handler.NewAuthHandler(authService, rateLimitService)
	oauthHandler     := handler.NewOAuthHandler(googleOAuthService, githubOAuthService, facebookOAuthService)
	bioHandler       := handler.NewBioHandler(bioService)
	previewHandler   := handler.NewPreviewHandler(previewService, linkService)
	folderHandler    := handler.NewFolderHandler(folderService)
	apiKeyHandler    := handler.NewAPIKeyHandler(apiKeyService)
	webhookHandler   := handler.NewWebhookHandler(webhookService)
	pixelHandler     := handler.NewPixelHandler(pixelService)
	splitHandler     := handler.NewSplitHandler(splitService)
	geoHandler       := handler.NewGeoHandler(geoRoutingService)
	linkHandler      := handler.NewLinkHandler(linkService, webhookService)
	redirectHandler  := handler.NewRedirectHandler(redirectService, clickService, geoService, pixelService, webhookService, linkRepo)
	qrHandler        := handler.NewQRHandler(qrService, linkService, s.Cfg.App)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService, linkRepo, clickEventRepo, s.Cfg.App)
	demoHandler      := handler.NewDemoHandler(demoService)

	// ── Swagger UI ───────────────────────────────────────────────────────────
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Health ────────────────────────────────────────────────────────────────
	router.GET("/health", healthHandler.Check)

	// ── Redirect (:slug must be registered before API groups) ─────────────────
	router.GET("/:slug", redirectHandler.Redirect)

	// ── API v1 ────────────────────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		// Public routes (no auth required)
		v1Public := v1.Group("")
		{
			// Auth
			v1Public.POST("/auth/signup", authHandler.Signup)
			v1Public.POST("/auth/login", authHandler.Login)
			v1Public.GET("/auth/rate-limit-status", authHandler.GetRateLimitStatus)
			v1Public.POST("/auth/refresh", authHandler.RefreshToken)
			v1Public.POST("/auth/forgot-password", authHandler.ForgotPassword)
			v1Public.POST("/auth/reset-password", authHandler.ResetPassword)
			v1Public.POST("/auth/reactivate", authHandler.ReactivateAccount)
			v1Public.POST("/auth/session/validate", authHandler.ValidateSession)
			v1Public.POST("/auth/verify-email", authHandler.VerifyEmail)
			v1Public.POST("/auth/resend-verification", authHandler.ResendVerificationEmail)

			// OAuth
			v1Public.GET("/auth/google", oauthHandler.GoogleLogin)
			v1Public.GET("/auth/google/callback", oauthHandler.GoogleCallback)
			v1Public.GET("/auth/github", oauthHandler.GitHubLogin)
			v1Public.GET("/auth/github/callback", oauthHandler.GitHubCallback)
			v1Public.GET("/auth/facebook", oauthHandler.FacebookLogin)
			v1Public.GET("/auth/facebook/callback", oauthHandler.FacebookCallback)

			// Demo shorten (rate-limited per IP)
			v1Public.POST("/demo/shorten", demoHandler.ShortenURL)

			// SSE live click counter (inline JWT via ?token= because EventSource can't set headers)
			v1Public.GET("/links/:id/analytics/live", analyticsHandler.StreamLiveCount)

			// Public bio page
			v1Public.GET("/bio/public/:username", bioHandler.GetPublicBioPage)
		}

		// Protected routes (JWT or API key required)
		v1Auth := v1.Group("")
		v1Auth.Use(middlewares.AuthOrAPIKeyMiddleware(s.Cfg.App, tokenBlacklistRepo, apiKeyService))
		v1Auth.Use(middlewares.SessionActivityMiddleware(sessionRepo, s.Cfg.Session))
		{
			// Auth profile & session management
			v1Auth.GET("/auth/profile", authHandler.GetProfile)
			v1Auth.PUT("/auth/profile", authHandler.UpdateProfile)
			v1Auth.PUT("/auth/change-password", authHandler.ChangePassword)
			v1Auth.POST("/auth/logout", authHandler.Logout)
			v1Auth.DELETE("/auth/account", authHandler.DeleteAccount)
			v1Auth.GET("/auth/sessions", authHandler.GetSessions)
			v1Auth.DELETE("/auth/sessions/:id", authHandler.DeleteSession)
			v1Auth.POST("/auth/logout-all", authHandler.LogoutAll)

			// Folders
			v1Auth.GET("/folders", folderHandler.ListFolders)
			v1Auth.POST("/folders", folderHandler.CreateFolder)
			v1Auth.PUT("/folders/:id", folderHandler.UpdateFolder)
			v1Auth.DELETE("/folders/:id", folderHandler.DeleteFolder)

			// Link CRUD
			v1Auth.GET("/links", linkHandler.ListLinks)
			v1Auth.GET("/links/duplicate", linkHandler.CheckDuplicate)
			v1Auth.POST("/links", linkHandler.CreateLink)
			v1Auth.POST("/links/import", linkHandler.ImportLinks)
			v1Auth.GET("/links/:id", linkHandler.GetLink)
			v1Auth.PUT("/links/:id", linkHandler.UpdateLink)
			v1Auth.DELETE("/links/:id", linkHandler.DeleteLink)
			v1Auth.PATCH("/links/:id/star", linkHandler.ToggleStar)
			v1Auth.POST("/links/:id/health-check", linkHandler.CheckLinkHealth)

			// A/B split testing
			v1Auth.PATCH("/links/:id/split-test", splitHandler.ToggleSplitTest)
			v1Auth.GET("/links/:id/variants", splitHandler.ListVariants)
			v1Auth.POST("/links/:id/variants", splitHandler.CreateVariant)
			v1Auth.PUT("/links/:id/variants/:variantId", splitHandler.UpdateVariant)
			v1Auth.DELETE("/links/:id/variants/:variantId", splitHandler.DeleteVariant)

			// Retargeting pixels
			v1Auth.PATCH("/links/:id/pixel-tracking", pixelHandler.TogglePixelTracking)
			v1Auth.GET("/links/:id/pixels", pixelHandler.ListPixels)
			v1Auth.POST("/links/:id/pixels", pixelHandler.CreatePixel)
			v1Auth.DELETE("/links/:id/pixels/:pixelId", pixelHandler.DeletePixel)

			// Geo routing
			v1Auth.PATCH("/links/:id/geo-routing", geoHandler.ToggleGeoRouting)
			v1Auth.GET("/links/:id/geo-rules", geoHandler.ListGeoRules)
			v1Auth.POST("/links/:id/geo-rules", geoHandler.CreateGeoRule)
			v1Auth.PUT("/links/:id/geo-rules/:ruleId", geoHandler.UpdateGeoRule)
			v1Auth.DELETE("/links/:id/geo-rules/:ruleId", geoHandler.DeleteGeoRule)

			// Link extras
			v1Auth.GET("/links/:id/qr", qrHandler.GetQRCode)
			v1Auth.GET("/links/:id/analytics", analyticsHandler.GetLinkAnalytics)
			v1Auth.GET("/links/:id/preview", previewHandler.GetLinkPreview)

			// API key management (JWT only — can't bootstrap with API key)
			v1Auth.GET("/api-keys", apiKeyHandler.ListAPIKeys)
			v1Auth.POST("/api-keys", apiKeyHandler.CreateAPIKey)
			v1Auth.DELETE("/api-keys/:id", apiKeyHandler.RevokeAPIKey)

			// Webhooks
			v1Auth.GET("/webhooks", webhookHandler.ListWebhooks)
			v1Auth.POST("/webhooks", webhookHandler.CreateWebhook)
			v1Auth.PUT("/webhooks/:id", webhookHandler.UpdateWebhook)
			v1Auth.DELETE("/webhooks/:id", webhookHandler.DeleteWebhook)

			// Bio page (link-in-bio)
			v1Auth.GET("/bio", bioHandler.GetBioPage)
			v1Auth.PUT("/bio", bioHandler.UpdateBioPage)
			v1Auth.GET("/bio/links", bioHandler.ListBioLinks)
			v1Auth.POST("/bio/links", bioHandler.CreateBioLink)
			v1Auth.PATCH("/bio/links/reorder", bioHandler.ReorderBioLinks)
			v1Auth.PUT("/bio/links/:id", bioHandler.UpdateBioLink)
			v1Auth.DELETE("/bio/links/:id", bioHandler.DeleteBioLink)
		}
	}

	s.Gin = router
	logger.Info("Routes configured successfully")
}

func (s *Server) setupMiddleware(router *gin.Engine) {
	logger.Debug("Setting up middlewares")

	router.Use(gin.Recovery())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.LoggerMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middlewares.CorsMiddleware(s.Cfg.App.CorsAllowedOrigins))
	router.Use(middlewares.SecurityHeaders())

	logger.Debug("Middlewares configured")
}
