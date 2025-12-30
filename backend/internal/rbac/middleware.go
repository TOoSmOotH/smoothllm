package rbac

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
)

type Middleware struct {
	enforcer *Enforcer
}

func NewMiddleware(enforcer *Enforcer) *Middleware {
	return &Middleware{
		enforcer: enforcer,
	}
}

func (m *Middleware) Authorize(obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := auth.GetUserID(c)
		role := auth.GetUserRole(c)

		fmt.Printf("RBAC Debug: userID=%d, role=%s, obj=%s, act=%s\n", userID, role, obj, act)

		if userID == 0 || role == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Use role instead of userID for authorization
		allowed, err := m.enforcer.Enforce(role, obj, act)
		fmt.Printf("RBAC Debug: allowed=%v, err=%v\n", allowed, err)
		if err != nil {
			c.JSON(500, gin.H{"error": "error checking permissions"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(403, gin.H{"error": fmt.Sprintf("access denied to %s %s", act, obj)})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *Middleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := auth.GetUserRole(c)

		if userRole == "" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
