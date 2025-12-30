package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/smoothweb/backend/internal/models"
	"github.com/smoothweb/backend/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	db           *gorm.DB
	jwtService   *JWTService
	rbacEnforcer RBACEnforcer
}

// RBACEnforcer interface to avoid circular dependency
type RBACEnforcer interface {
	AddRoleForUser(userID, role string) (bool, error)
}

func NewService(db *gorm.DB, jwtService *JWTService, rbacEnforcer RBACEnforcer) *Service {
	return &Service{
		db:           db,
		jwtService:   jwtService,
		rbacEnforcer: rbacEnforcer,
	}
}

func (s *Service) RegisterUser(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	var existingUser models.User
	result := s.db.Where("email = ?", req.Email).Or("username = ?", req.Username).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("user already exists")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var userCount int64
	s.db.Model(&models.User{}).Count(&userCount)

	role := "user"
	status := "active"
	if userCount == 0 {
		role = "admin"
	} else {
		registrationEnabled, err := s.getBoolSetting(models.AppSettingRegistrationEnabledKey, true)
		if err != nil {
			return nil, fmt.Errorf("failed to read registration settings: %w", err)
		}
		if !registrationEnabled {
			return nil, errors.New("registration is currently disabled")
		}

		autoApprove, err := s.getBoolSetting(models.AppSettingAutoApproveNewUsersKey, true)
		if err != nil {
			return nil, fmt.Errorf("failed to read registration settings: %w", err)
		}
		if !autoApprove {
			status = "pending"
		}
	}

	user := models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         role,
		Status:       status,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Add role assignment to Casbin
	if s.rbacEnforcer != nil {
		userIDStr := fmt.Sprintf("%d", user.ID)
		if _, err := s.rbacEnforcer.AddRoleForUser(userIDStr, user.Role); err != nil {
			// Log error but don't fail registration
			fmt.Printf("Failed to add role assignment for user %d: %v\n", user.ID, err)
		}
	}

	if user.Status != "active" {
		return &models.RegisterResponse{
			User:     user,
			Approved: false,
			Message:  "Registration received. Your account is pending approval.",
		}, nil
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.RegisterResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
		Approved:     true,
	}, nil
}

func (s *Service) LoginUser(req *models.LoginRequest) (*models.AuthResponse, error) {
	var user models.User
	result := s.db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if user.Status != "active" {
		return nil, errors.New("account is pending approval")
	}

	now := time.Now()
	user.LastActiveAt = &now
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update last active time: %w", err)
	}

	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *Service) getBoolSetting(key string, defaultValue bool) (bool, error) {
	var setting models.AppSetting
	if err := s.db.Where("key = ?", key).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return defaultValue, nil
		}
		return defaultValue, err
	}

	if setting.Value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseBool(setting.Value)
	if err != nil {
		return defaultValue, nil
	}
	return parsed, nil
}

func (s *Service) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	userID, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var user models.User
	result := s.db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	newAccessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &models.AuthResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *Service) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", result.Error)
	}
	return &user, nil
}
