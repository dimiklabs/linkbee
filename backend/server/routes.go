package server

import (
	"context"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/shafikshaon/shortlink/handler"
	"github.com/shafikshaon/shortlink/logger"
	"github.com/shafikshaon/shortlink/middlewares"
	"github.com/shafikshaon/shortlink/repository"
	analyticsSvc "github.com/shafikshaon/shortlink/service/analytics"
	authSrv "github.com/shafikshaon/shortlink/service/auth"
	folderSvc "github.com/shafikshaon/shortlink/service/folder"
	clickSvc "github.com/shafikshaon/shortlink/service/click"
	demoSvc "github.com/shafikshaon/shortlink/service/demo"
	emailSrv "github.com/shafikshaon/shortlink/service/email"
	facebookSrv "github.com/shafikshaon/shortlink/service/facebook"
	githubSrv "github.com/shafikshaon/shortlink/service/github"
	googleSrv "github.com/shafikshaon/shortlink/service/google"
	healthSrv "github.com/shafikshaon/shortlink/service/health"
	linkSvc "github.com/shafikshaon/shortlink/service/link"
	qrSvc "github.com/shafikshaon/shortlink/service/qr"
	rateLimitSrv "github.com/shafikshaon/shortlink/service/ratelimit"
	redirectSvc "github.com/shafikshaon/shortlink/service/redirect"
	userSrv "github.com/shafikshaon/shortlink/service/user"
	"github.com/shafikshaon/shortlink/worker"
)

func (s *Server) ConfigureRoutes(ctx context.Context, router *gin.Engine) {
	logger.Debug("Configuring routes")

	s.setupMiddleware(router)

	// ── Repositories ─────────────────────────────────────────────────────────
	folderRepo            := repository.NewFolderRepository(s.MasterDB, s.ReplicaDB)
	userRepo              := repository.NewUserRepository(s.MasterDB, s.ReplicaDB)
	passwordResetRepo     := repository.NewPasswordResetRepository(s.MasterDB, s.ReplicaDB)
	tokenBlacklistRepo   := repository.NewTokenBlacklistRepository(s.Cache)
	sessionRepo          := repository.NewSessionRepository(s.MasterDB, s.ReplicaDB)
	emailVerificationRepo := repository.NewEmailVerificationRepository(s.MasterDB, s.ReplicaDB)
	oauthStateRepo       := repository.NewOAuthStateRepository(s.Cache)
	rateLimitRepo        := repository.NewRateLimitRepository(s.Cache)
	linkRepo             := repository.NewLinkRepository(s.MasterDB, s.ReplicaDB)
	clickEventRepo       := repository.NewClickEventRepository(s.MasterDB, s.ReplicaDB)

	// ── Services ──────────────────────────────────────────────────────────────
	healthService       := healthSrv.NewHealthService(s.MasterDB, s.ReplicaDB, s.Cache, s.Cfg.App.Env)
	userService         := userSrv.NewUserService(userRepo)
	emailService        := emailSrv.NewEmailService(s.Cfg.Email, emailVerificationRepo, userRepo)
	authService         := authSrv.NewAuthService(userService, passwordResetRepo, tokenBlacklistRepo, sessionRepo, emailService, s.Cfg.App, s.Cfg.Session)
	rateLimitService    := rateLimitSrv.NewRateLimitService(rateLimitRepo, s.Cfg.RateLimit)
	googleOAuthService  := googleSrv.NewGoogleOAuthService(s.Cfg.Google, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)
	githubOAuthService  := githubSrv.NewGitHubOAuthService(s.Cfg.GitHub, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)
	facebookOAuthService := facebookSrv.NewFacebookOAuthService(s.Cfg.Facebook, s.Cfg.App, userRepo, oauthStateRepo, sessionRepo, tokenBlacklistRepo)

	folderService    := folderSvc.NewFolderService(folderRepo)
	linkService      := linkSvc.NewLinkService(linkRepo, s.Cfg.App, s.Cfg.Link)
	redirectService  := redirectSvc.NewRedirectService(linkRepo, s.Cache, s.Cfg.Link.CacheTTLSeconds)
	clickService     := clickSvc.NewClickService(s.Cache)
	qrService        := qrSvc.NewQRService()
	analyticsService := analyticsSvc.NewAnalyticsService(linkRepo, clickEventRepo)
	demoService      := demoSvc.NewDemoService(linkRepo, s.Cache, s.Cfg.App, s.Cfg.Link)

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
	folderHandler    := handler.NewFolderHandler(folderService)
	linkHandler      := handler.NewLinkHandler(linkService)
	redirectHandler  := handler.NewRedirectHandler(redirectService, clickService, linkRepo)
	qrHandler        := handler.NewQRHandler(qrService, linkService, s.Cfg.App)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService, linkRepo, clickEventRepo, s.Cfg.App)
	demoHandler      := handler.NewDemoHandler(demoService)

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
		}

		// Protected routes (JWT required)
		v1Auth := v1.Group("")
		v1Auth.Use(middlewares.AuthMiddleware(s.Cfg.App, tokenBlacklistRepo))
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
			v1Auth.POST("/links", linkHandler.CreateLink)
			v1Auth.POST("/links/import", linkHandler.ImportLinks)
			v1Auth.GET("/links/:id", linkHandler.GetLink)
			v1Auth.PUT("/links/:id", linkHandler.UpdateLink)
			v1Auth.DELETE("/links/:id", linkHandler.DeleteLink)
			v1Auth.PATCH("/links/:id/star", linkHandler.ToggleStar)
			v1Auth.POST("/links/:id/health-check", linkHandler.CheckLinkHealth)

			// Link extras
			v1Auth.GET("/links/:id/qr", qrHandler.GetQRCode)
			v1Auth.GET("/links/:id/analytics", analyticsHandler.GetLinkAnalytics)
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
