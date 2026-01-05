package custom

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/config"
	"github.com/smoothweb/backend/internal/custom/handlers"
	"github.com/smoothweb/backend/internal/custom/services"
	"github.com/smoothweb/backend/internal/rbac"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB     *gorm.DB
	Config *config.Config
	JWT    *auth.JWTService
	RBAC   *rbac.Middleware
}

// RegisterRoutes lets downstream projects add routes without touching core wiring.
func RegisterRoutes(v1 *gin.RouterGroup, deps Dependencies) {
	// Run custom migrations
	if err := AutoMigrate(deps.DB); err != nil {
		log.Fatalf("Failed to run custom migrations: %v", err)
	}

	// Initialize services
	providerService := services.NewProviderService(deps.DB)
	keyService := services.NewKeyService(deps.DB)
	usageService := services.NewUsageService(deps.DB)
	oauthService := services.NewOAuthService(deps.DB, providerService, deps.Config.FrontendURL)

	// Initialize handlers
	providerHandler := handlers.NewProviderHandler(providerService)
	keyHandler := handlers.NewKeyHandler(keyService)
	usageHandler := handlers.NewUsageHandler(usageService)
	oauthHandler := handlers.NewOAuthHandler(oauthService, deps.Config.FrontendURL)

	// Provider routes (protected with JWT)
	providers := v1.Group("/providers")
	providers.Use(auth.AuthMiddleware(deps.JWT))
	{
		providers.GET("", providerHandler.ListProviders)
		providers.POST("", providerHandler.CreateProvider)
		providers.POST("/test", providerHandler.TestConnectionWithCredentials)
		providers.GET("/:id", providerHandler.GetProvider)
		providers.PUT("/:id", providerHandler.UpdateProvider)
		providers.DELETE("/:id", providerHandler.DeleteProvider)
		providers.POST("/:id/test", providerHandler.TestConnection)
	}

	// OAuth routes for Claude Max
	oauth := v1.Group("/oauth/anthropic")
	{
		// Authorize endpoint requires JWT auth (user must be logged in to start OAuth)
		oauth.GET("/authorize", auth.AuthMiddleware(deps.JWT), oauthHandler.Authorize)
		// Callback doesn't require JWT auth (comes from Anthropic redirect)
		oauth.GET("/callback", oauthHandler.Callback)
		// Disconnect and test require JWT auth
		oauth.POST("/disconnect/:id", auth.AuthMiddleware(deps.JWT), oauthHandler.Disconnect)
		oauth.POST("/test/:id", auth.AuthMiddleware(deps.JWT), oauthHandler.TestOAuthConnection)
	}

	// API Key routes (protected with JWT)
	keys := v1.Group("/keys")
	keys.Use(auth.AuthMiddleware(deps.JWT))
	{
		keys.GET("", keyHandler.ListKeys)
		keys.POST("", keyHandler.CreateKey)
		keys.GET("/:id", keyHandler.GetKey)
		keys.PUT("/:id", keyHandler.UpdateKey)
		keys.DELETE("/:id", keyHandler.DeleteKey)
		keys.POST("/:id/revoke", keyHandler.RevokeKey)
	}

	// Usage routes (protected with JWT)
	usage := v1.Group("/usage")
	usage.Use(auth.AuthMiddleware(deps.JWT))
	{
		usage.GET("", usageHandler.GetUsageSummary)
		usage.GET("/daily", usageHandler.GetDailyUsage)
		usage.GET("/by-key", usageHandler.GetUsageByKey)
		usage.GET("/by-provider", usageHandler.GetUsageByProvider)
		usage.GET("/by-model", usageHandler.GetUsageByModel)
		usage.GET("/recent", usageHandler.GetRecentUsage)
	}
}

// RegisterProxyRoutes registers the LLM proxy routes at /v1 (outside /api/v1 group)
// This is necessary for OpenAI API compatibility - clients expect /v1/chat/completions
// These routes use proxy API key authentication, not JWT
func RegisterProxyRoutes(router *gin.Engine, deps Dependencies) {
	// Initialize services
	providerService := services.NewProviderService(deps.DB)
	keyService := services.NewKeyService(deps.DB)
	usageService := services.NewUsageService(deps.DB)
	oauthService := services.NewOAuthService(deps.DB, providerService, deps.Config.FrontendURL)
	proxyService := services.NewProxyService(keyService, providerService, usageService, oauthService)

	// Initialize proxy handler
	proxyHandler := handlers.NewProxyHandler(proxyService)

	// Proxy routes at /v1 (OpenAI-compatible endpoints)
	// These use proxy API key authentication (Bearer sk-smoothllm-xxx), not JWT
	v1Proxy := router.Group("/v1")
	{
		// OpenAI-compatible chat completions endpoint
		v1Proxy.POST("/chat/completions", proxyHandler.ChatCompletions)

		// OpenAI-compatible models list endpoint
		v1Proxy.GET("/models", proxyHandler.ListModels)
	}
}
