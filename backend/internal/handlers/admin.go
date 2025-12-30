package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/models"
	"github.com/smoothweb/backend/internal/utils"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

type StatisticsResponse struct {
	TotalUsers        int64     `json:"total_users"`
	ActiveUsers       int64     `json:"active_users"`
	AdminUsers        int64     `json:"admin_users"`
	ProfilesCompleted int64     `json:"profiles_completed"`
	NewUsersToday     int64     `json:"new_users_today"`
	NewUsersWeek      int64     `json:"new_users_week"`
	AvgCompletion     float64   `json:"avg_completion"`
	LastUpdated       time.Time `json:"last_updated"`
}

// GetStatistics retrieves platform statistics
// GET /api/v1/admin/stats
func (h *AdminHandler) GetStatistics(c *gin.Context) {
	var stats StatisticsResponse

	// Count total users
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)

	// Count active users (logged in within last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	h.db.Model(&models.User{}).Where("last_active_at >= ?", sevenDaysAgo).Count(&stats.ActiveUsers)

	// Count admin users
	h.db.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.AdminUsers)

	// Count users with complete profiles (100% completion)
	h.db.Model(&models.ProfileCompletion{}).Where("score = 100").Count(&stats.ProfilesCompleted)

	// Count new users today
	today := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.NewUsersToday)

	// Count new users this week
	weekAgo := time.Now().AddDate(0, 0, -7)
	h.db.Model(&models.User{}).Where("created_at >= ?", weekAgo).Count(&stats.NewUsersWeek)

	// Calculate average completion score
	var avgCompletion sql.NullFloat64
	h.db.Model(&models.ProfileCompletion{}).Select("AVG(score)").Scan(&avgCompletion)
	if avgCompletion.Valid {
		stats.AvgCompletion = avgCompletion.Float64
	}

	stats.LastUpdated = time.Now()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ListUsers retrieves all users with pagination and filters
// GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	role := c.Query("role")
	search := c.Query("search")

	offset := (page - 1) * limit

	query := h.db.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	var users []models.User
	var total int64

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users)

	// Remove passwords from response
	for i := range users {
		users[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
				"pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// CreateUser creates a user account (admin only).
// POST /api/v1/admin/users
func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role" binding:"omitempty,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	result := h.db.Where("email = ?", req.Email).Or("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         role,
		Status:       "active",
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	enforcer, exists := c.Get("enforcer")
	if exists {
		if enf, ok := enforcer.(interface {
			AddRoleForUser(string, string) (bool, error)
		}); ok {
			userIDStr := strconv.FormatUint(uint64(user.ID), 10)
			enf.AddRoleForUser(userIDStr, user.Role)
		}
	}

	user.PasswordHash = ""

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

// DeleteUser deletes a user account
// DELETE /api/v1/admin/users/:id
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Check if user exists
	var user models.User
	result := h.db.First(&user, uint(userID))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Prevent deleting self
	currentUserID := auth.GetUserID(c)
	if uint(userID) == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	// Delete user (cascade will handle related records)
	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user deleted successfully",
	})
}

// ChangeUserRole changes a user's role
// PATCH /api/v1/admin/users/:id/role
func (h *AdminHandler) ChangeUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var user models.User
	result := h.db.First(&user, uint(userID))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Prevent changing own role
	currentUserID := auth.GetUserID(c)
	if uint(userID) == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change your own role"})
		return
	}

	// Ensure at least one admin remains
	if req.Role == "user" && user.Role == "admin" {
		var adminCount int64
		h.db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot remove last admin"})
			return
		}
	}

	// Update user role
	if err := h.db.Model(&user).Update("role", req.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user role"})
		return
	}

	enforcer, exists := c.Get("enforcer")
	if exists {
		if enf, ok := enforcer.(interface {
			AddRoleForUser(string, string) (bool, error)
			DeleteRoleForUser(string, string) (bool, error)
		}); ok {
			userIDStr := strconv.FormatUint(uint64(userID), 10)
			enf.DeleteRoleForUser(userIDStr, user.Role)
			enf.AddRoleForUser(userIDStr, req.Role)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user role updated successfully",
		"data": gin.H{
			"user_id": user.ID,
			"role":    req.Role,
		},
	})
}

// ApproveUser sets a pending user to active.
// PATCH /api/v1/admin/users/:id/approve
func (h *AdminHandler) ApproveUser(c *gin.Context) {
	userIDStr := c.Param("id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var user models.User
	result := h.db.First(&user, uint(userID))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if user.Status == "active" {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "user already active",
		})
		return
	}

	if err := h.db.Model(&user).Update("status", "active").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user approved successfully",
		"data": gin.H{
			"user_id": user.ID,
			"status":  "active",
		},
	})
}
