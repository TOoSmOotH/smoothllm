package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ProfileCompletion tracks the completion status of a user's profile
type ProfileCompletion struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null" json:"user_id"`

	// Completion score (0-100)
	Score int `gorm:"default:0;index" json:"score"`

	// Track which fields are completed
	CompletedFields datatypes.JSON `gorm:"type:json" json:"completed_fields"`

	// Required vs optional fields
	RequiredCompleted int `gorm:"default:0" json:"required_completed"`
	RequiredTotal     int `gorm:"default:0" json:"required_total"`
	OptionalCompleted int `gorm:"default:0" json:"optional_completed"`
	OptionalTotal     int `gorm:"default:0" json:"optional_total"`

	// Milestones
	FirstCompletionAt *time.Time `json:"first_completion_at,omitempty"`
	LastUpdatedAt     time.Time  `json:"last_updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// CompletionMilestone defines milestones for profile completion
type CompletionMilestone struct {
	gorm.Model

	// Milestone configuration
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Completion threshold (percentage)
	Threshold int `gorm:"not null;index" json:"threshold"`

	// Rewards
	RewardType  string `gorm:"type:varchar(50)" json:"reward_type"` // badge, points, feature_unlock
	RewardValue string `gorm:"type:varchar(255)" json:"reward_value"`

	// Display
	IconURL  string `gorm:"type:varchar(255)" json:"icon_url,omitempty"`
	IsActive bool   `gorm:"default:true" json:"is_active"`

	// Track which users achieved this
	Achievements []CompletionAchievement `gorm:"foreignKey:MilestoneID" json:"-"`
}

// CompletionAchievement tracks when users achieve milestones
type CompletionAchievement struct {
	gorm.Model
	UserID       uint `gorm:"index;not null" json:"user_id"`
	MilestoneID  uint `gorm:"index;not null" json:"milestone_id"`

	AchievedAt          time.Time `gorm:"not null" json:"achieved_at"`
	ScoreAtAchievement  int       `gorm:"not null" json:"score_at_achievement"`

	User      User               `gorm:"foreignKey:UserID" json:"-"`
	Milestone CompletionMilestone `gorm:"foreignKey:MilestoneID" json:"milestone,omitempty"`
}