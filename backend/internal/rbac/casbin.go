package rbac

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Enforcer struct {
	*casbin.Enforcer
}

func NewEnforcer(db *gorm.DB, policyPath string) (*Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer(policyPath, adapter)
	if err != nil {
		return nil, err
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	policies := [][]string{
		{"admin", "/api/v1/*", "*"},
		{"admin", "/api/v1/admin/*", "*"},
		{"admin", "/api/v1/admin/stats", "GET"},
		{"admin", "/api/v1/admin/users", "GET"},
		{"admin", "/api/v1/admin/users", "POST"},
		{"admin", "/api/v1/admin/users/:id", "DELETE"},
		{"admin", "/api/v1/admin/users/:id/role", "PATCH"},
		{"admin", "/api/v1/admin/users/:id/approve", "PATCH"},
		{"admin", "/api/v1/admin/settings/theme", "PUT"},
		{"admin", "/api/v1/admin/settings/registration", "GET"},
		{"admin", "/api/v1/admin/settings/registration", "PUT"},
		{"user", "/api/v1/auth/*", "*"},
		{"user", "/api/v1/profile", "*"},
		{"user", "/api/v1/profile", "GET"},
		{"user", "/api/v1/profile", "PUT"},
		{"user", "/api/v1/privacy", "*"},
		{"user", "/api/v1/media/avatar", "POST"},
		{"user", "/api/v1/media/cover", "POST"},
		{"user", "/api/v1/media/:id", "GET"},
		{"user", "/api/v1/media/:id", "DELETE"},
		{"user", "/api/v1/media/:id/crop", "POST"},
		{"user", "/api/v1/media/user/:userId", "GET"},
		{"user", "/api/v1/social", "GET"},
		{"user", "/api/v1/social", "POST"},
		{"user", "/api/v1/social/:id", "*"},
		{"user", "/api/v1/social/reorder", "PUT"},
		{"anonymous", "/api/v1/social/user/:userId", "GET"},
		{"user", "/api/v1/completion", "GET"},
		{"user", "/api/v1/completion/recalculate", "POST"},
		{"user", "/api/v1/completion/milestones", "GET"},
		{"user", "/api/v1/completion/leaderboard", "GET"},
	}

	for _, policy := range policies {
		added, err := enforcer.AddPolicy(policy)
		if err != nil {
			log.Printf("Failed to add policy %v: %v", policy, err)
		} else if added {
			log.Printf("Added policy: %v", policy)
		}
	}

	// Note: Role assignments should be done when users are created/updated
	// This is handled in the auth service

	return &Enforcer{Enforcer: enforcer}, nil
}

func (e *Enforcer) AddRoleForUser(userID, role string) (bool, error) {
	return e.Enforcer.AddRoleForUser(userID, role)
}

func (e *Enforcer) DeleteRoleForUser(userID, role string) (bool, error) {
	return e.Enforcer.DeleteRoleForUser(userID, role)
}

func (e *Enforcer) DeleteRolesForUser(userID string) (bool, error) {
	return e.Enforcer.DeleteRolesForUser(userID)
}

func (e *Enforcer) GetRolesForUser(userID string) ([]string, error) {
	return e.Enforcer.GetRolesForUser(userID)
}

func (e *Enforcer) HasRoleForUser(userID, role string) (bool, error) {
	return e.Enforcer.HasRoleForUser(userID, role)
}
