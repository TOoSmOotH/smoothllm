package services

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/models"
)

type CompletionService struct {
	db *gorm.DB
}

func NewCompletionService(db *gorm.DB) *CompletionService {
	return &CompletionService{db: db}
}

type CompletionScoreResponse struct {
	UserID            uint                    `json:"user_id"`
	Score             int                     `json:"score"`
	MaxScore          int                     `json:"max_score"`
	Percentage        float64                 `json:"percentage"`
	IsComplete        bool                    `json:"is_complete"`
	CompletedFields   []string                `json:"completed_fields"`
	MissingFields     []CompletionFieldInfo   `json:"missing_fields"`
	CategoryBreakdown map[string]CategoryInfo `json:"category_breakdown"`
	NextRecommended   []CompletionFieldInfo   `json:"next_recommended"`
}

type CompletionFieldInfo struct {
	Field    string `json:"field"`
	Label    string `json:"label"`
	Points   int    `json:"points"`
	IsFilled bool   `json:"is_filled"`
}

type CategoryInfo struct {
	Completed int `json:"completed"`
	Total     int `json:"total"`
	Points    int `json:"points"`
	MaxPoints int `json:"max_points"`
}

type MilestoneResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Threshold   int    `json:"threshold"`
	RewardType  string `json:"reward_type"`
	RewardValue string `json:"reward_value"`
	IconURL     string `json:"icon_url,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type LeaderboardEntry struct {
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Score       int       `json:"score"`
	Rank        int       `json:"rank"`
	CompletedAt time.Time `json:"completed_at"`
}

func (s *CompletionService) GetCompletionScore(userID uint) (*CompletionScoreResponse, error) {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return s.buildDefaultCompletionResponse(userID), nil
	}

	var completion models.ProfileCompletion
	if err := s.db.Where("user_id = ?", userID).First(&completion).Error; err != nil {
		return s.buildDefaultCompletionResponse(userID), nil
	}

	return s.buildCompletionResponse(&profile, &completion), nil
}

func (s *CompletionService) RecalculateCompletionScore(userID uint) (*CompletionScoreResponse, error) {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	score := s.calculateCompletionScore(&profile)

	var completion models.ProfileCompletion
	if err := s.db.Where("user_id = ?", userID).First(&completion).Error; err != nil {
		return nil, fmt.Errorf("failed to get completion: %w", err)
	}

	completion.Score = score
	completion.LastUpdatedAt = time.Now()
	completion.CompletedFields = s.buildCompletedFields(&profile)

	if err := s.db.Save(&completion).Error; err != nil {
		return nil, fmt.Errorf("failed to update completion: %w", err)
	}

	return s.buildCompletionResponse(&profile, &completion), nil
}

func (s *CompletionService) GetMilestones() ([]MilestoneResponse, error) {
	var milestones []models.CompletionMilestone
	if err := s.db.Where("is_active = ?", true).Order("threshold ASC").Find(&milestones).Error; err != nil {
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}

	responses := make([]MilestoneResponse, len(milestones))
	for i, m := range milestones {
		responses[i] = MilestoneResponse{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			Threshold:   m.Threshold,
			RewardType:  m.RewardType,
			RewardValue: m.RewardValue,
			IconURL:     m.IconURL,
			IsActive:    m.IsActive,
		}
	}

	return responses, nil
}

func (s *CompletionService) GetLeaderboard(limit int) ([]LeaderboardEntry, error) {
	var completions []models.ProfileCompletion
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	if err := s.db.Order("score DESC").Limit(limit).Find(&completions).Error; err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	leaderboard := make([]LeaderboardEntry, len(completions))
	for i, completion := range completions {
		var user models.User
		if err := s.db.Where("id = ?", completion.UserID).First(&user).Error; err != nil {
			continue
		}

		leaderboard[i] = LeaderboardEntry{
			UserID:      completion.UserID,
			Username:    user.Username,
			Score:       completion.Score,
			Rank:        i + 1,
			CompletedAt: completion.LastUpdatedAt,
		}
	}

	return leaderboard, nil
}

func (s *CompletionService) buildCompletionResponse(profile *models.UserProfile, completion *models.ProfileCompletion) *CompletionScoreResponse {
	maxScore := 100
	percentage := float64(completion.Score) / float64(maxScore) * 100

	return &CompletionScoreResponse{
		UserID:            profile.UserID,
		Score:             completion.Score,
		MaxScore:          maxScore,
		Percentage:        percentage,
		IsComplete:        completion.Score >= maxScore,
		CategoryBreakdown: s.buildCategoryBreakdown(profile, completion),
		NextRecommended:   s.buildNextRecommendedFields(profile),
		CompletedFields:   s.buildStringSliceFromJSON(completion.CompletedFields),
		MissingFields:     s.buildMissingFields(profile, s.buildStringSliceFromJSON(completion.CompletedFields)),
	}
}

func (s *CompletionService) buildDefaultCompletionResponse(userID uint) *CompletionScoreResponse {
	return &CompletionScoreResponse{
		UserID:            userID,
		Score:             0,
		MaxScore:          100,
		Percentage:        0,
		IsComplete:        false,
		CategoryBreakdown: s.buildDefaultCategoryBreakdown(),
		NextRecommended:   s.buildNextRecommendedFieldsForNewUser(),
	}
}

func (s *CompletionService) calculateCompletionScore(profile *models.UserProfile) int {
	score := 0

	if profile.FirstName != nil && *profile.FirstName != "" {
		score += 10
	}
	if profile.LastName != nil && *profile.LastName != "" {
		score += 10
	}
	if profile.DisplayName != nil && *profile.DisplayName != "" {
		score += 10
	}
	if profile.Bio != nil && *profile.Bio != "" {
		score += 10
	}
	if profile.Phone != nil && *profile.Phone != "" {
		score += 5
	}
	if profile.Website != nil && *profile.Website != "" {
		score += 5
	}
	if profile.Location != nil && *profile.Location != "" {
		score += 10
	}
	if profile.Birthday != nil {
		score += 5
	}
	if profile.Gender != nil && *profile.Gender != "" {
		score += 5
	}
	if profile.Pronouns != nil && *profile.Pronouns != "" {
		score += 5
	}
	if profile.Language != nil && *profile.Language != "" {
		score += 5
	}
	if profile.JobTitle != nil && *profile.JobTitle != "" {
		score += 5
	}
	if profile.Company != nil && *profile.Company != "" {
		score += 5
	}
	if profile.LinkedInURL != nil && *profile.LinkedInURL != "" {
		score += 5
	}
	if profile.PortfolioURL != nil && *profile.PortfolioURL != "" {
		score += 5
	}

	if profile.Skills != nil {
		var skills []string
		json.Unmarshal(profile.Skills, &skills)
		if len(skills) > 0 {
			score += 5
		}
	}

	if profile.Interests != nil {
		var interests []string
		json.Unmarshal(profile.Interests, &interests)
		if len(interests) > 0 {
			score += 5
		}
	}

	return score
}

func (s *CompletionService) buildCompletedFields(profile *models.UserProfile) datatypes.JSON {
	completed := []string{}

	if profile.FirstName != nil && *profile.FirstName != "" {
		completed = append(completed, "first_name")
	}
	if profile.LastName != nil && *profile.LastName != "" {
		completed = append(completed, "last_name")
	}
	if profile.DisplayName != nil && *profile.DisplayName != "" {
		completed = append(completed, "display_name")
	}
	if profile.Bio != nil && *profile.Bio != "" {
		completed = append(completed, "bio")
	}
	if profile.Phone != nil && *profile.Phone != "" {
		completed = append(completed, "phone")
	}
	if profile.Website != nil && *profile.Website != "" {
		completed = append(completed, "website")
	}
	if profile.Location != nil && *profile.Location != "" {
		completed = append(completed, "location")
	}
	if profile.Birthday != nil {
		completed = append(completed, "birthday")
	}
	if profile.Gender != nil && *profile.Gender != "" {
		completed = append(completed, "gender")
	}
	if profile.Pronouns != nil && *profile.Pronouns != "" {
		completed = append(completed, "pronouns")
	}
	if profile.Language != nil && *profile.Language != "" {
		completed = append(completed, "language")
	}
	if profile.JobTitle != nil && *profile.JobTitle != "" {
		completed = append(completed, "job_title")
	}
	if profile.Company != nil && *profile.Company != "" {
		completed = append(completed, "company")
	}
	if profile.LinkedInURL != nil && *profile.LinkedInURL != "" {
		completed = append(completed, "linkedin_url")
	}
	if profile.PortfolioURL != nil && *profile.PortfolioURL != "" {
		completed = append(completed, "portfolio_url")
	}

	data, _ := json.Marshal(completed)
	return datatypes.JSON(data)
}

func (s *CompletionService) buildMissingFields(profile *models.UserProfile, completed []string) []CompletionFieldInfo {
	allFields := s.getAllFieldDefinitions()
	missing := []CompletionFieldInfo{}

	for _, field := range allFields {
		isFilled := s.isFieldFilled(profile, field.Field)
		if !isFilled {
			missing = append(missing, CompletionFieldInfo{
				Field:    field.Field,
				Label:    field.Label,
				Points:   field.Points,
				IsFilled: false,
			})
		}
	}

	return missing
}

func (s *CompletionService) buildMissingFieldsForNewUser() []CompletionFieldInfo {
	allFields := s.getAllFieldDefinitions()
	missing := make([]CompletionFieldInfo, len(allFields))

	for i, field := range allFields {
		missing[i] = CompletionFieldInfo{
			Field:    field.Field,
			Label:    field.Label,
			Points:   field.Points,
			IsFilled: false,
		}
	}

	return missing
}

func (s *CompletionService) getAllFieldDefinitions() []CompletionFieldInfo {
	return []CompletionFieldInfo{
		{"first_name", "First Name", 10, false},
		{"last_name", "Last Name", 10, false},
		{"display_name", "Display Name", 10, false},
		{"bio", "Bio", 10, false},
		{"phone", "Phone Number", 5, false},
		{"website", "Website", 5, false},
		{"location", "Location", 10, false},
		{"birthday", "Birthday", 5, false},
		{"gender", "Gender", 5, false},
		{"pronouns", "Pronouns", 5, false},
		{"language", "Language", 5, false},
		{"job_title", "Job Title", 5, false},
		{"company", "Company", 5, false},
		{"linkedin_url", "LinkedIn URL", 5, false},
		{"portfolio_url", "Portfolio URL", 5, false},
		{"skills", "Skills", 5, false},
		{"interests", "Interests", 5, false},
	}
}

func (s *CompletionService) isFieldFilled(profile *models.UserProfile, fieldName string) bool {
	switch fieldName {
	case "first_name":
		return profile.FirstName != nil && *profile.FirstName != ""
	case "last_name":
		return profile.LastName != nil && *profile.LastName != ""
	case "display_name":
		return profile.DisplayName != nil && *profile.DisplayName != ""
	case "bio":
		return profile.Bio != nil && *profile.Bio != ""
	case "phone":
		return profile.Phone != nil && *profile.Phone != ""
	case "website":
		return profile.Website != nil && *profile.Website != ""
	case "location":
		return profile.Location != nil && *profile.Location != ""
	case "birthday":
		return profile.Birthday != nil
	case "gender":
		return profile.Gender != nil && *profile.Gender != ""
	case "pronouns":
		return profile.Pronouns != nil && *profile.Pronouns != ""
	case "language":
		return profile.Language != nil && *profile.Language != ""
	case "job_title":
		return profile.JobTitle != nil && *profile.JobTitle != ""
	case "company":
		return profile.Company != nil && *profile.Company != ""
	case "linkedin_url":
		return profile.LinkedInURL != nil && *profile.LinkedInURL != ""
	case "portfolio_url":
		return profile.PortfolioURL != nil && *profile.PortfolioURL != ""
	case "skills":
		return profile.Skills != nil
	case "interests":
		return profile.Interests != nil
	default:
		return false
	}
}

func (s *CompletionService) buildCategoryBreakdown(profile *models.UserProfile, completion *models.ProfileCompletion) map[string]CategoryInfo {
	basic := s.calculateCategoryPoints(profile, []string{"first_name", "last_name", "display_name", "bio"}, 10)
	contact := s.calculateCategoryPoints(profile, []string{"phone", "website", "location"}, 5)
	personal := s.calculateCategoryPoints(profile, []string{"birthday", "gender", "pronouns", "language"}, 5)
	professional := s.calculateCategoryPoints(profile, []string{"job_title", "company", "linkedin_url", "portfolio_url"}, 5)
	extras := s.calculateCategoryPoints(profile, []string{"skills", "interests"}, 5)

	return map[string]CategoryInfo{
		"basic_info":    basic,
		"contact_info":  contact,
		"personal_info": personal,
		"professional":  professional,
		"extras":        extras,
	}
}

func (s *CompletionService) calculateCategoryPoints(profile *models.UserProfile, fields []string, maxPoints int) CategoryInfo {
	total := len(fields) * maxPoints
	completed := 0

	for _, field := range fields {
		if s.isFieldFilled(profile, field) {
			completed++
		}
	}

	return CategoryInfo{
		Completed: completed,
		Total:     total,
		Points:    completed * maxPoints,
		MaxPoints: total,
	}
}

func (s *CompletionService) buildDefaultCategoryBreakdown() map[string]CategoryInfo {
	return map[string]CategoryInfo{
		"basic_info":    {Completed: 0, Total: 40, Points: 0, MaxPoints: 40},
		"contact_info":  {Completed: 0, Total: 20, Points: 0, MaxPoints: 20},
		"personal_info": {Completed: 0, Total: 20, Points: 0, MaxPoints: 20},
		"professional":  {Completed: 0, Total: 20, Points: 0, MaxPoints: 20},
		"extras":        {Completed: 0, Total: 10, Points: 0, MaxPoints: 10},
	}
}

func (s *CompletionService) buildNextRecommendedFields(profile *models.UserProfile) []CompletionFieldInfo {
	allFields := s.getAllFieldDefinitions()
	recommended := []CompletionFieldInfo{}
	highPriorityFields := []string{"first_name", "last_name", "display_name", "bio", "location"}

	for _, priorityField := range highPriorityFields {
		if !s.isFieldFilled(profile, priorityField) {
			for _, field := range allFields {
				if field.Field == priorityField {
					recommended = append(recommended, field)
					break
				}
			}
		}
	}

	if len(recommended) >= 3 {
		return recommended
	}

	for _, field := range allFields {
		if !s.isFieldFilled(profile, field.Field) {
			recommended = append(recommended, field)
			if len(recommended) >= 3 {
				break
			}
		}
	}

	return recommended
}

func (s *CompletionService) buildNextRecommendedFieldsForNewUser() []CompletionFieldInfo {
	allFields := s.getAllFieldDefinitions()
	recommended := []CompletionFieldInfo{}

	for i := 0; i < 3; i++ {
		for _, field := range allFields {
			if field.Field == []string{"first_name", "last_name", "display_name", "bio"}[i] {
				recommended = append(recommended, field)
				break
			}
		}
	}

	return recommended
}

func (s *CompletionService) buildStringSliceFromJSON(data datatypes.JSON) []string {
	if data == nil {
		return []string{}
	}

	var slice []string
	json.Unmarshal(data, &slice)
	return slice
}
