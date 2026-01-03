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

	// Initialize handlers
	providerHandler := handlers.NewProviderHandler(providerService)
	keyHandler := handlers.NewKeyHandler(keyService)

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
}
