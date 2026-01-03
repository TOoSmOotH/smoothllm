package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/config"
	"github.com/smoothweb/backend/internal/custom"
	"github.com/smoothweb/backend/internal/database"
	"github.com/smoothweb/backend/internal/handlers"
	"github.com/smoothweb/backend/internal/middleware"
	"github.com/smoothweb/backend/internal/rbac"
	"github.com/smoothweb/backend/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	gin.SetMode(cfg.GinMode)

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	db, err := database.NewDatabase(cfg.DBPath, cfg.DBEncryptionKey)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to run auto-migration: %v", err)
	}

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTAccessTokenExpiry, cfg.JWTRefreshTokenExpiry)

	policyPath := "./configs/casbin_policy.conf"
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		log.Fatalf("Policy file not found at: %s", policyPath)
	}

	enforcer, err := rbac.NewEnforcer(db.GetDB(), policyPath)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin enforcer: %v", err)
	}

	authService := auth.NewService(db.GetDB(), jwtService, enforcer)
	profileService := services.NewProfileService(db.GetDB())
	privacyService := services.NewPrivacyService(db.GetDB())
	mediaService := services.NewMediaService(db.GetDB())
	socialLinkService := services.NewSocialLinkService(db.GetDB())
	completionService := services.NewCompletionService(db.GetDB())

	// Create upload directories
	if err := mediaService.EnsureUploadDirectories(); err != nil {
		log.Fatalf("Failed to create upload directories: %v", err)
	}

	router := gin.Default()

	router.Use(middleware.CORS(cfg.AllowedOrigins))
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	rbacMiddleware := rbac.NewMiddleware(enforcer)

	v1 := router.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler(authService)
		profileHandler := handlers.NewProfileHandler(profileService)
		privacyHandler := handlers.NewPrivacyHandler(privacyService)
		mediaHandler := handlers.NewMediaHandler(mediaService, db.GetDB())
		socialLinkHandler := handlers.NewSocialLinkHandler(socialLinkService)
		completionHandler := handlers.NewCompletionHandler(completionService)
		adminHandler := handlers.NewAdminHandler(db.GetDB())
		settingsHandler := handlers.NewSettingsHandler(db.GetDB())

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)

			protected := authGroup.Use(auth.AuthMiddleware(jwtService))
			{
				protected.GET("/me", authHandler.GetCurrentUser)
			}
		}

		// Public profile routes (no authentication required)
		profilePublic := v1.Group("/profile")
		{
			profilePublic.GET("/:username", profileHandler.GetProfile)
			profilePublic.GET("/id/:id", profileHandler.GetProfileByID)
			profilePublic.GET("/check-username", profileHandler.CheckUsernameAvailability)
		}

		// Public settings routes
		v1.GET("/settings/theme", settingsHandler.GetTheme)

		// Protected profile routes (authentication required)
		protected := v1.Group("/")
		protected.Use(auth.AuthMiddleware(jwtService))
		{
			// Privacy routes
			privacyGroup := protected.Group("/privacy")
			privacyGroup.GET("", rbacMiddleware.Authorize("/api/v1/privacy", "GET"), privacyHandler.GetPrivacySettings)
			privacyGroup.PUT("", rbacMiddleware.Authorize("/api/v1/privacy", "PUT"), privacyHandler.UpdatePrivacySettings)
			privacyGroup.POST("/preset/:name", rbacMiddleware.Authorize("/api/v1/privacy/preset/:name", "POST"), privacyHandler.ApplyPrivacyPreset)
			privacyGroup.GET("/presets", rbacMiddleware.Authorize("/api/v1/privacy/presets", "GET"), privacyHandler.GetPrivacyPresets)

			// Profile routes
			protected.GET("/profile", rbacMiddleware.Authorize("/api/v1/profile", "GET"), profileHandler.GetMyProfile)
			protected.PUT("/profile", rbacMiddleware.Authorize("/api/v1/profile", "PUT"), profileHandler.UpdateProfile)
			protected.POST("/profile", rbacMiddleware.Authorize("/api/v1/profile", "POST"), profileHandler.CreateProfile)
			protected.DELETE("/profile", rbacMiddleware.Authorize("/api/v1/profile", "DELETE"), profileHandler.DeleteProfile)

			admin := protected.Group("/admin")
			admin.Use(rbacMiddleware.RequireRole("admin"))
			{
				// Admin statistics
				admin.GET("/stats", rbacMiddleware.Authorize("/api/v1/admin/stats", "GET"), adminHandler.GetStatistics)

				// User management
				admin.GET("/users", rbacMiddleware.Authorize("/api/v1/admin/users", "GET"), adminHandler.ListUsers)
				admin.POST("/users", rbacMiddleware.Authorize("/api/v1/admin/users", "POST"), adminHandler.CreateUser)
				admin.DELETE("/users/:id", rbacMiddleware.Authorize("/api/v1/admin/users/:id", "DELETE"), adminHandler.DeleteUser)
				admin.PATCH("/users/:id/role", rbacMiddleware.Authorize("/api/v1/admin/users/:id/role", "PATCH"), adminHandler.ChangeUserRole)
				admin.PATCH("/users/:id/approve", rbacMiddleware.Authorize("/api/v1/admin/users/:id/approve", "PATCH"), adminHandler.ApproveUser)

				// Admin settings
				admin.PUT("/settings/theme", rbacMiddleware.Authorize("/api/v1/admin/settings/theme", "PUT"), settingsHandler.UpdateTheme)
				admin.GET("/settings/registration", rbacMiddleware.Authorize("/api/v1/admin/settings/registration", "GET"), settingsHandler.GetRegistrationSettings)
				admin.PUT("/settings/registration", rbacMiddleware.Authorize("/api/v1/admin/settings/registration", "PUT"), settingsHandler.UpdateRegistrationSettings)

				// Admin privacy routes
				admin.GET("/users/:userId/privacy", rbacMiddleware.Authorize("/api/v1/admin/users/:userId/privacy", "GET"), privacyHandler.GetPrivacySettingsByAdmin)
				admin.PUT("/users/:userId/privacy", rbacMiddleware.Authorize("/api/v1/admin/users/:userId/privacy", "PUT"), privacyHandler.UpdatePrivacySettingsByAdmin)
			}
		}

		// Media routes (authentication required)
		mediaGroup := v1.Group("/media")
		mediaGroup.Use(auth.AuthMiddleware(jwtService))
		{
			mediaGroup.POST("/avatar", rbacMiddleware.Authorize("/api/v1/media/avatar", "POST"), mediaHandler.UploadAvatar)
			mediaGroup.POST("/cover", rbacMiddleware.Authorize("/api/v1/media/cover", "POST"), mediaHandler.UploadCoverPhoto)
			mediaGroup.POST("/:id/crop", rbacMiddleware.Authorize("/api/v1/media/:id/crop", "POST"), mediaHandler.CropMedia)
			mediaGroup.DELETE("/:id", rbacMiddleware.Authorize("/api/v1/media/:id", "DELETE"), mediaHandler.DeleteMedia)
			mediaGroup.GET("/:id", rbacMiddleware.Authorize("/api/v1/media/:id", "GET"), mediaHandler.GetMedia)
			mediaGroup.GET("/user/:userId", rbacMiddleware.Authorize("/api/v1/media/user/:userId", "GET"), mediaHandler.GetUserMedia)
		}

		// Social links routes (authentication required)
		socialGroup := v1.Group("/social")
		socialGroup.Use(auth.AuthMiddleware(jwtService))
		{
			socialGroup.GET("", rbacMiddleware.Authorize("/api/v1/social", "GET"), socialLinkHandler.GetSocialLinks)
			socialGroup.POST("", rbacMiddleware.Authorize("/api/v1/social", "POST"), socialLinkHandler.AddSocialLink)
			socialGroup.GET("/:id", rbacMiddleware.Authorize("/api/v1/social/:id", "GET"), socialLinkHandler.GetSocialLink)
			socialGroup.PUT("/:id", rbacMiddleware.Authorize("/api/v1/social/:id", "PUT"), socialLinkHandler.UpdateSocialLink)
			socialGroup.DELETE("/:id", rbacMiddleware.Authorize("/api/v1/social/:id", "DELETE"), socialLinkHandler.DeleteSocialLink)
			socialGroup.PUT("/reorder", rbacMiddleware.Authorize("/api/v1/social/reorder", "PUT"), socialLinkHandler.ReorderSocialLinks)
		}

		// Public social links route (no authentication required)
		v1.GET("/social/user/:userId", socialLinkHandler.GetUserSocialLinks)

		// Completion routes (authentication required)
		completionGroup := v1.Group("/completion")
		completionGroup.Use(auth.AuthMiddleware(jwtService))
		{
			completionGroup.GET("", rbacMiddleware.Authorize("/api/v1/completion", "GET"), completionHandler.GetCompletionScore)
			completionGroup.POST("/recalculate", rbacMiddleware.Authorize("/api/v1/completion/recalculate", "POST"), completionHandler.RecalculateCompletionScore)
			completionGroup.GET("/milestones", rbacMiddleware.Authorize("/api/v1/completion/milestones", "GET"), completionHandler.GetMilestones)
			completionGroup.GET("/leaderboard", rbacMiddleware.Authorize("/api/v1/completion/leaderboard", "GET"), completionHandler.GetLeaderboard)
		}

		custom.RegisterRoutes(v1, custom.Dependencies{
			DB:     db.GetDB(),
			Config: cfg,
			JWT:    jwtService,
			RBAC:   rbacMiddleware,
		})
	}

	// Register LLM proxy routes at /v1 (outside /api/v1 for OpenAI compatibility)
	custom.RegisterProxyRoutes(router, custom.Dependencies{
		DB:     db.GetDB(),
		Config: cfg,
		JWT:    jwtService,
		RBAC:   rbacMiddleware,
	})

	// Static file serving for uploaded media (public, no auth required)
	router.Static("/uploads", "./uploads")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func autoMigrate(db *database.Database) error {
	return db.AutoMigrate()
}
