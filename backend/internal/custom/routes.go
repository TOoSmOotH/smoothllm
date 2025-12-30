package custom

import (
	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/config"
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
	_ = v1
	_ = deps
}
