package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smoothweb/backend/internal/auth"
	"github.com/smoothweb/backend/internal/services"
)

type CompletionHandler struct {
	completionService *services.CompletionService
}

func NewCompletionHandler(completionService *services.CompletionService) *CompletionHandler {
	return &CompletionHandler{
		completionService: completionService,
	}
}

// GetCompletionScore retrieves profile completion score
// GET /api/v1/completion
func (h *CompletionHandler) GetCompletionScore(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	score, err := h.completionService.GetCompletionScore(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get completion score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    score,
	})
}

// RecalculateCompletionScore forces recalculation of completion score
// POST /api/v1/completion/recalculate
func (h *CompletionHandler) RecalculateCompletionScore(c *gin.Context) {
	userID := auth.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	score, err := h.completionService.RecalculateCompletionScore(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to recalculate completion score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    score,
		"message": "completion score recalculated successfully",
	})
}

// GetMilestones retrieves all active milestones
// GET /api/v1/completion/milestones
func (h *CompletionHandler) GetMilestones(c *gin.Context) {
	milestones, err := h.completionService.GetMilestones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get milestones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    milestones,
	})
}

// GetLeaderboard retrieves top users by completion score
// GET /api/v1/completion/leaderboard
func (h *CompletionHandler) GetLeaderboard(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	leaderboard, err := h.completionService.GetLeaderboard(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    leaderboard,
	})
}
