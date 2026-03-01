package server

import (
	"context"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/shafikshaon/linkbee/docs" // swagger generated docs
	"github.com/shafikshaon/linkbee/handler"
	"github.com/shafikshaon/linkbee/logger"
	"github.com/shafikshaon/linkbee/middlewares"
	"github.com/shafikshaon/linkbee/repository"
	"github.com/shafikshaon/linkbee/util"
	analyticsSvc   "github.com/shafikshaon/linkbee/service/analytics"
	dashboardSvc   "github.com/shafikshaon/linkbee/service/dashboard"
	expirySvc      "github.com/shafikshaon/linkbee/service/expiry"
	reportingSvc   "github.com/shafikshaon/linkbee/service/reporting"
	apiKeySvc "github.com/shafikshaon/linkbee/service/apikey"
	authSrv "github.com/shafikshaon/linkbee/service/auth"
	clickSvc "github.com/shafikshaon/linkbee/service/click"
	demoSvc "github.com/shafikshaon/linkbee/service/demo"
	emailSrv "github.com/shafikshaon/linkbee/service/email"
	folderSvc "github.com/shafikshaon/linkbee/service/folder"
	geoSvc "github.com/shafikshaon/linkbee/service/geo"
	georoutingSvc "github.com/shafikshaon/linkbee/service/georouting"
	githubSrv "github.com/shafikshaon/linkbee/service/github"
	googleSrv "github.com/shafikshaon/linkbee/service/google"
	healthSrv "github.com/shafikshaon/linkbee/service/health"
	linkSvc "github.com/shafikshaon/linkbee/service/link"
	qrSvc "github.com/shafikshaon/linkbee/service/qr"
	rateLimitSrv "github.com/shafikshaon/linkbee/service/ratelimit"
	redirectSvc "github.com/shafikshaon/linkbee/service/redirect"
	splitSvc "github.com/shafikshaon/linkbee/service/split"
	userSrv "github.com/shafikshaon/linkbee/service/user"
	adminSvc   "github.com/shafikshaon/linkbee/service/admin"
	auditSvc   "github.com/shafikshaon/linkbee/service/audit"
	billingSvc "github.com/shafikshaon/linkbee/service/billing"
	bioSvc     "github.com/shafikshaon/linkbee/service/bio"
	domainSvc  "github.com/shafikshaon/linkbee/service/domain"
	pixelSvc   "github.com/shafikshaon/linkbee/service/pixel"
	previewSvc "github.com/shafikshaon/linkbee/service/preview"
	teamSvc    "github.com/shafikshaon/linkbee/service/team"
	webhookSvc "github.com/shafikshaon/linkbee/service/webhook"
	"github.com/shafikshaon/linkbee/worker"
)

func (s *Server) ConfigureRoutes(ctx context.Context, router *gin.Engine) {
	logger.Debug("Configuring routes")

	router.MaxMultipartMemory = 5 << 20 // 5 MB

	s.setupMiddleware(router)

	// ── Repositories ─────────────────────────────────────────────────────────
	reportRepo            := repository.NewAnalyticsReportRepository(s.MasterDB, s.ReplicaDB)
	subRepo               := repository.NewSubscriptionRepository(s.MasterDB, s.ReplicaDB)
	bioRepo                    := repository.NewBioRepository(s.MasterDB, s.ReplicaDB)
	bioLinkClickEventRepo      := repository.NewBioLinkClickEventRepository(s.MasterDB, s.ReplicaDB)
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
	customDomainRepo      := repository.NewCustomDomainRepository(s.MasterDB, s.ReplicaDB)
	auditLogRepo          := repository.NewAuditLogRepository(s.MasterDB, s.ReplicaDB)
	totpBackupCodeRepo    := repository.NewTotpBackupCodeRepository(s.MasterDB, s.ReplicaDB)
	webhookDeliveryRepo   := repository.NewWebhookDeliveryRepository(s.MasterDB, s.ReplicaDB)
	teamRepo              := repository.NewTeamRepository(s.MasterDB, s.ReplicaDB)

	// ── Services ──────────────────────────────────────────────────────────────
	healthService        := healthSrv.NewHealthService(s.MasterDB, s.ReplicaDB, s.Cache, s.Cfg.App.Env)
	userService          := userSrv.NewUserService(userRepo)
	emailService         := emailSrv.NewEmailService(s.Cfg.Email, emailVerificationRepo, userRepo)
	authService          := authSrv.NewAuthService(userService, passwordResetRepo, tokenBlacklistRepo, sessionRepo, totpBackupCodeRepo, emailService, s.Cfg.App, s.Cfg.Session)
	rateLimitService     := rateLimitSrv.NewRateLimitService(rateLimitRepo, s.Cfg.RateLimit)
	googleOAuthService   := googleSrv.NewGoogleOAuthService(s.Cfg.Google, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)
	githubOAuthService   := githubSrv.NewGitHubOAuthService(s.Cfg.GitHub, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)
	billingService     := billingSvc.NewBillingService(subRepo, s.Cfg.Billing)
	planEnforcer       := billingSvc.NewPlanEnforcer(subRepo, linkRepo, apiKeyRepo, webhookRepo)
	adminService       := adminSvc.NewAdminService(userRepo, linkRepo, clickEventRepo)
	bioService         := bioSvc.NewBioService(bioRepo)
	bioClickService    := clickSvc.NewBioClickService(s.Cache, bioRepo)
	previewService     := previewSvc.NewPreviewService(s.Cache)
	folderService      := folderSvc.NewFolderService(folderRepo, clickEventRepo)
	apiKeyService      := apiKeySvc.NewAPIKeyService(apiKeyRepo)
	webhookService     := webhookSvc.NewWebhookService(webhookRepo, webhookDeliveryRepo)
	pixelService       := pixelSvc.NewPixelService(pixelRepo, linkRepo)
	splitService       := splitSvc.NewSplitService(variantRepo, linkRepo)
	slugGen            := util.NewSlugGenerator(s.Cfg.Link.SlugSecret, s.Cfg.Link.SlugLength)
	linkService        := linkSvc.NewLinkService(linkRepo, s.Cache, slugGen, s.Cfg.App, s.Cfg.Link)
	geoService         := geoSvc.NewGeoService(s.Cfg.App.GeoDBPath)
	geoRoutingService  := georoutingSvc.NewGeoRoutingService(geoRuleRepo, linkRepo)
	redirectService    := redirectSvc.NewRedirectService(linkRepo, variantRepo, geoRuleRepo, s.Cache, s.Cfg.Link.CacheTTLSeconds)
	clickService       := clickSvc.NewClickService(s.Cache)
	qrService          := qrSvc.NewQRService()
	analyticsService   := analyticsSvc.NewAnalyticsService(linkRepo, clickEventRepo)
	dashboardService   := dashboardSvc.NewDashboardService(linkRepo, clickEventRepo, s.Cfg.App)
	reportingService   := reportingSvc.NewReportingService(reportRepo, userRepo, clickEventRepo, linkRepo, emailService, s.Cfg.App, s.Cfg.Email)
	expiryService      := expirySvc.NewExpiryService(linkRepo, userRepo, emailService, s.Cfg.App, s.Cfg.Email)
	demoService        := demoSvc.NewDemoService(linkRepo, s.Cache, slugGen, s.Cfg.App, s.Cfg.Link)
	domainService      := domainSvc.NewDomainService(customDomainRepo)
	auditService       := auditSvc.NewAuditService(auditLogRepo)
	teamService        := teamSvc.NewTeamService(teamRepo, userRepo)

	// ── Report worker (background goroutine) ──────────────────────────────────
	reportWorker := worker.NewReportWorker(reportingService)
	go reportWorker.Start(ctx)

	// ── Expiry notification worker (background goroutine) ─────────────────────
	expiryWorker := worker.NewExpiryWorker(expiryService)
	go expiryWorker.Start(ctx)

	// ── Click worker (background goroutine) ───────────────────────────────────
	clickWorker := worker.NewClickWorker(
		s.Cache,
		clickEventRepo,
		linkRepo,
		s.Cfg.Link.ClickQueueBatchSize,
		s.Cfg.Link.ClickQueueFlushInterval,
	)
	go clickWorker.Start(ctx)

	// ── Bio click worker (background goroutine) ────────────────────────────────
	bioClickWorker := worker.NewBioClickWorker(
		s.Cache,
		bioLinkClickEventRepo,
		bioRepo,
		s.Cfg.Link.ClickQueueBatchSize,
		s.Cfg.Link.ClickQueueFlushInterval,
	)
	go bioClickWorker.Start(ctx)

	// ── Health worker (background goroutine) ──────────────────────────────────
	healthWorker := worker.NewHealthWorker(linkRepo)
	go healthWorker.Start(ctx)

	// ── Handlers ──────────────────────────────────────────────────────────────
	healthHandler    := handler.NewHealthHandler(healthService)
	authHandler      := handler.NewAuthHandler(authService, rateLimitService, auditService)
	oauthHandler     := handler.NewOAuthHandler(googleOAuthService, githubOAuthService, s.Cfg.App)
	billingHandler   := handler.NewBillingHandler(billingService, linkRepo, apiKeyRepo, webhookRepo)
	adminHandler     := handler.NewAdminHandler(adminService, s.Cfg.App)
	exportHandler    := handler.NewExportHandler(userRepo, linkRepo, s.Cfg.App)
	uploadsDir := "./uploads"
	bioHandler       := handler.NewBioHandler(bioService, bioClickService, planEnforcer, uploadsDir, "https://"+s.Cfg.App.BaseDomain)
	previewHandler   := handler.NewPreviewHandler(previewService, linkService)
	folderHandler    := handler.NewFolderHandler(folderService)
	apiKeyHandler    := handler.NewAPIKeyHandler(apiKeyService, planEnforcer, auditService)
	webhookHandler   := handler.NewWebhookHandler(webhookService, planEnforcer)
	pixelHandler     := handler.NewPixelHandler(pixelService)
	splitHandler     := handler.NewSplitHandler(splitService)
	geoHandler       := handler.NewGeoHandler(geoRoutingService)
	linkHandler      := handler.NewLinkHandler(linkService, webhookService, planEnforcer, auditService)
	auditLogHandler  := handler.NewAuditLogHandler(auditLogRepo)
	domainHandler    := handler.NewDomainHandler(domainService, auditService)
	redirectHandler  := handler.NewRedirectHandler(redirectService, clickService, geoService, pixelService, webhookService, linkRepo, customDomainRepo, s.Cfg.App.BaseDomain)
	qrHandler        := handler.NewQRHandler(qrService, linkService, s.Cfg.App)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService, linkRepo, clickEventRepo, s.Cfg.App, planEnforcer)
	dashboardHandler := handler.NewDashboardHandler(dashboardService, planEnforcer)
	reportHandler    := handler.NewReportHandler(reportingService, planEnforcer)
	demoHandler      := handler.NewDemoHandler(demoService)
	teamHandler      := handler.NewTeamHandler(teamService)

	// ── Static file serving (avatars, etc.) ──────────────────────────────────
	router.Static("/uploads", "./uploads")

	// ── Swagger UI ───────────────────────────────────────────────────────────
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Health ────────────────────────────────────────────────────────────────
	router.GET("/health", healthHandler.Check)

	// ── Redirect (:slug must be registered before API groups) ─────────────────
	router.GET("/:slug", redirectHandler.Redirect)

	// ── API v1 ────────────────────────────────────────────────────────────────
	v1 := router.Group("/api/v1")
	{
		// Billing webhook (public — verified via HMAC signature)
		v1.POST("/billing/webhook", billingHandler.PaddleWebhook)

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
			v1Public.POST("/auth/totp/verify-login", authHandler.VerifyTOTPLogin)

			// OAuth
			v1Public.GET("/auth/google", oauthHandler.GoogleLogin)
			v1Public.GET("/auth/google/callback", oauthHandler.GoogleCallback)
			v1Public.GET("/auth/github", oauthHandler.GitHubLogin)
			v1Public.GET("/auth/github/callback", oauthHandler.GitHubCallback)
			// Demo shorten (rate-limited per IP)
			v1Public.POST("/demo/shorten", demoHandler.ShortenURL)

			// SSE live click counter (inline JWT via ?token= because EventSource can't set headers)
			v1Public.GET("/links/:id/analytics/live", analyticsHandler.StreamLiveCount)

			// Public bio page
			v1Public.GET("/bio/public/:username", bioHandler.GetPublicBioPage)
			// Bio link click tracking (fire-and-forget, no auth)
			v1Public.POST("/bio/public/:username/links/:id/click", bioHandler.RecordBioLinkClick)
		}

		// Protected routes (JWT or API key required)
		v1Auth := v1.Group("")
		v1Auth.Use(middlewares.AuthOrAPIKeyMiddleware(s.Cfg.App, tokenBlacklistRepo, apiKeyService))
		v1Auth.Use(middlewares.SessionActivityMiddleware(sessionRepo, s.Cfg.Session))
		{
			// Dashboard
			v1Auth.GET("/dashboard/overview", dashboardHandler.GetOverview)
			v1Auth.GET("/dashboard/analytics", dashboardHandler.GetGlobalAnalytics)
			v1Auth.GET("/dashboard/analytics/comparison", dashboardHandler.GetGlobalAnalyticsComparison)

			// TOTP / 2FA
			v1Auth.GET("/auth/totp/status", authHandler.GetTOTPStatus)
			v1Auth.GET("/auth/totp/setup", authHandler.SetupTOTP)
			v1Auth.POST("/auth/totp/confirm", authHandler.ConfirmTOTP)
			v1Auth.DELETE("/auth/totp/disable", authHandler.DisableTOTP)

			// Auth profile & session management
			v1Auth.GET("/auth/profile", authHandler.GetProfile)
			v1Auth.PUT("/auth/profile", authHandler.UpdateProfile)
			v1Auth.PUT("/auth/change-password", authHandler.ChangePassword)
			v1Auth.POST("/auth/logout", authHandler.Logout)
			v1Auth.DELETE("/auth/account", authHandler.DeleteAccount)
			v1Auth.GET("/auth/data-export", exportHandler.ExportData)
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
			v1Auth.GET("/links/tags", linkHandler.GetUserTags)
			v1Auth.GET("/links/export", exportHandler.ExportLinksCSV)
			v1Auth.GET("/links/duplicate", linkHandler.CheckDuplicate)
			v1Auth.GET("/links/comparison", analyticsHandler.GetMultiLinkComparison)
			v1Auth.POST("/links", linkHandler.CreateLink)
			v1Auth.POST("/links/bulk", linkHandler.BulkAction)
			v1Auth.POST("/links/import", linkHandler.ImportLinks)
			v1Auth.GET("/links/:id", linkHandler.GetLink)
			v1Auth.PUT("/links/:id", linkHandler.UpdateLink)
			v1Auth.DELETE("/links/:id", linkHandler.DeleteLink)
			v1Auth.PATCH("/links/:id/star", linkHandler.ToggleStar)
			v1Auth.POST("/links/:id/clone", linkHandler.CloneLink)
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
			v1Auth.GET("/links/:id/analytics/comparison", analyticsHandler.GetPeriodComparison)
			v1Auth.GET("/links/:id/preview", previewHandler.GetLinkPreview)

			// Scheduled analytics reports
			v1Auth.GET("/reports", reportHandler.ListReports)
			v1Auth.POST("/reports", reportHandler.CreateReport)
			v1Auth.GET("/reports/:id", reportHandler.GetReport)
			v1Auth.PUT("/reports/:id", reportHandler.UpdateReport)
			v1Auth.DELETE("/reports/:id", reportHandler.DeleteReport)
			v1Auth.POST("/reports/:id/send", reportHandler.SendReportNow)
			v1Auth.GET("/reports/:id/deliveries", reportHandler.GetReportDeliveries)

			// API key management (JWT only — can't bootstrap with API key)
			v1Auth.GET("/api-keys", apiKeyHandler.ListAPIKeys)
			v1Auth.POST("/api-keys", apiKeyHandler.CreateAPIKey)
			v1Auth.DELETE("/api-keys/:id", apiKeyHandler.RevokeAPIKey)

			// Webhooks
			v1Auth.GET("/webhooks", webhookHandler.ListWebhooks)
			v1Auth.POST("/webhooks", webhookHandler.CreateWebhook)
			v1Auth.PUT("/webhooks/:id", webhookHandler.UpdateWebhook)
			v1Auth.DELETE("/webhooks/:id", webhookHandler.DeleteWebhook)
			v1Auth.GET("/webhooks/:id/secret", webhookHandler.GetWebhookSecret)
			v1Auth.POST("/webhooks/:id/test", webhookHandler.TestWebhook)
			v1Auth.GET("/webhooks/:id/deliveries", webhookHandler.GetDeliveries)
			v1Auth.POST("/webhooks/:id/deliveries/:deliveryId/resend", webhookHandler.ResendDelivery)

			// Billing (subscription + usage + checkout)
			v1Auth.GET("/billing/subscription", billingHandler.GetSubscription)
			v1Auth.GET("/billing/usage", billingHandler.GetUsage)
			v1Auth.GET("/billing/checkout/:plan", billingHandler.GetCheckoutURL)

			// Admin (role=admin only)
			v1Admin := v1Auth.Group("/admin")
			v1Admin.Use(middlewares.AdminMiddleware())
			{
				v1Admin.GET("/stats", adminHandler.GetStats)
				v1Admin.GET("/growth", adminHandler.GetGrowthTimeSeries)
				v1Admin.GET("/users", adminHandler.ListUsers)
				v1Admin.PATCH("/users/:id/status", adminHandler.UpdateUserStatus)
				v1Admin.PATCH("/users/:id/role", adminHandler.UpdateUserRole)
				v1Admin.POST("/users/:id/impersonate", adminHandler.ImpersonateUser)
			}

			// Audit logs
			v1Auth.GET("/audit-logs/export", auditLogHandler.ExportAuditLogs)
			v1Auth.GET("/audit-logs", auditLogHandler.ListAuditLogs)

			// Custom Domains
			v1Auth.GET("/domains", domainHandler.ListDomains)
			v1Auth.POST("/domains", domainHandler.AddDomain)
			v1Auth.POST("/domains/:id/verify", domainHandler.VerifyDomain)
			v1Auth.DELETE("/domains/:id", domainHandler.DeleteDomain)

			// Bio page (link-in-bio)
			v1Auth.GET("/bio", bioHandler.GetBioPage)
			v1Auth.PUT("/bio", bioHandler.UpdateBioPage)
			v1Auth.POST("/bio/avatar", bioHandler.UploadBioAvatar)
			v1Auth.GET("/bio/links", bioHandler.ListBioLinks)
			v1Auth.POST("/bio/links", bioHandler.CreateBioLink)
			v1Auth.PATCH("/bio/links/reorder", bioHandler.ReorderBioLinks)
			v1Auth.PUT("/bio/links/:id", bioHandler.UpdateBioLink)
			v1Auth.DELETE("/bio/links/:id", bioHandler.DeleteBioLink)

			// Teams (static routes must come before parameterized ones)
			v1Auth.POST("/teams/invite/accept", teamHandler.AcceptInvite)
			v1Auth.GET("/teams", teamHandler.ListMyTeams)
			v1Auth.POST("/teams", teamHandler.CreateTeam)
			v1Auth.GET("/teams/:id", teamHandler.GetTeam)
			v1Auth.PUT("/teams/:id", teamHandler.UpdateTeam)
			v1Auth.DELETE("/teams/:id", teamHandler.DeleteTeam)
			v1Auth.POST("/teams/:id/members", teamHandler.InviteMember)
			v1Auth.GET("/teams/:id/members", teamHandler.ListMembers)
			v1Auth.PATCH("/teams/:id/members/:userID/role", teamHandler.UpdateMemberRole)
			v1Auth.DELETE("/teams/:id/members/:userID", teamHandler.RemoveMember)
			v1Auth.POST("/teams/:id/leave", teamHandler.LeaveTeam)
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
