package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

// KeyService handles proxy API key CRUD operations
type KeyService struct {
	db *gorm.DB
}

// NewKeyService creates a new KeyService instance
func NewKeyService(db *gorm.DB) *KeyService {
	return &KeyService{db: db}
}

// KeyResponse represents the key data returned to clients
// Note: The full key is never stored or returned after creation
type KeyResponse struct {
	ID         uint       `json:"id"`
	UserID     uint       `json:"user_id"`
	ProviderID uint       `json:"provider_id"`
	KeyPrefix  string     `json:"key_prefix"`
	Name       string     `json:"name"`
	IsActive   bool       `json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	// Provider info for convenience
	ProviderName string `json:"provider_name,omitempty"`
	ProviderType string `json:"provider_type,omitempty"`
}

// KeyCreateResponse is returned when a key is created, includes the full key once
type KeyCreateResponse struct {
	KeyResponse
	Key string `json:"key"` // Full key, only returned on creation
}

// CreateKeyRequest represents the request to create a new proxy API key
type CreateKeyRequest struct {
	ProviderID uint       `json:"provider_id" binding:"required"`
	Name       string     `json:"name"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// UpdateKeyRequest represents the request to update a proxy API key
type UpdateKeyRequest struct {
	Name      *string `json:"name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ListKeys returns all proxy API keys for a user
func (s *KeyService) ListKeys(userID uint) ([]KeyResponse, error) {
	var keys []models.ProxyAPIKey
	if err := s.db.Preload("Provider").Where("user_id = ?", userID).Find(&keys).Error; err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	responses := make([]KeyResponse, len(keys))
	for i, key := range keys {
		responses[i] = s.buildKeyResponse(&key)
	}

	return responses, nil
}

// GetKey retrieves a key by ID for a specific user
func (s *KeyService) GetKey(userID, keyID uint) (*KeyResponse, error) {
	key, err := s.getKeyByID(userID, keyID)
	if err != nil {
		return nil, err
	}

	// Load provider info
	if err := s.db.Model(key).Association("Provider").Find(&key.Provider); err != nil {
		return nil, fmt.Errorf("failed to load provider: %w", err)
	}

	response := s.buildKeyResponse(key)
	return &response, nil
}

// CreateKey creates a new proxy API key for a user
func (s *KeyService) CreateKey(userID uint, req *CreateKeyRequest) (*KeyCreateResponse, error) {
	// Validate input
	if err := s.validateCreateRequest(userID, req); err != nil {
		return nil, err
	}

	// Generate the full API key
	fullKey, err := s.generateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	// Extract prefix for display (first 12 chars of the random part + "...")
	keyPrefix := s.extractKeyPrefix(fullKey)

	// Hash the full key for storage
	keyHash := s.hashKey(fullKey)

	// Create the key record
	key := models.ProxyAPIKey{
		UserID:     userID,
		ProviderID: req.ProviderID,
		KeyHash:    keyHash,
		KeyPrefix:  keyPrefix,
		Name:       req.Name,
		IsActive:   true,
		ExpiresAt:  req.ExpiresAt,
	}

	if err := s.db.Create(&key).Error; err != nil {
		return nil, fmt.Errorf("failed to create key: %w", err)
	}

	// Load provider info for response
	if err := s.db.Model(&key).Association("Provider").Find(&key.Provider); err != nil {
		return nil, fmt.Errorf("failed to load provider: %w", err)
	}

	// Build response with full key (only time it's returned)
	response := &KeyCreateResponse{
		KeyResponse: s.buildKeyResponse(&key),
		Key:         fullKey,
	}

	return response, nil
}

// UpdateKey updates an existing proxy API key
func (s *KeyService) UpdateKey(userID, keyID uint, req *UpdateKeyRequest) (*KeyResponse, error) {
	key, err := s.getKeyByID(userID, keyID)
	if err != nil {
		return nil, err
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.ExpiresAt != nil {
		updates["expires_at"] = req.ExpiresAt
	}

	if len(updates) > 0 {
		if err := s.db.Model(key).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update key: %w", err)
		}
	}

	// Refresh key data
	if err := s.db.Preload("Provider").First(key, keyID).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh key: %w", err)
	}

	response := s.buildKeyResponse(key)
	return &response, nil
}

// DeleteKey deletes (revokes) a proxy API key
func (s *KeyService) DeleteKey(userID, keyID uint) error {
	key, err := s.getKeyByID(userID, keyID)
	if err != nil {
		return err
	}

	if err := s.db.Delete(key).Error; err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// RevokeKey is an alias for DeleteKey - marks a key as inactive/deleted
func (s *KeyService) RevokeKey(userID, keyID uint) error {
	return s.DeleteKey(userID, keyID)
}

// ValidateKey validates a proxy API key and returns the associated key record
// This is used by the proxy to authenticate incoming requests
func (s *KeyService) ValidateKey(fullKey string) (*models.ProxyAPIKey, error) {
	// Validate key format
	if !strings.HasPrefix(fullKey, models.ProxyAPIKeyPrefix) {
		return nil, fmt.Errorf("invalid key format")
	}

	// Hash the provided key
	keyHash := s.hashKey(fullKey)

	// Look up by hash
	var key models.ProxyAPIKey
	if err := s.db.Preload("Provider").Where("key_hash = ?", keyHash).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid API key")
		}
		return nil, fmt.Errorf("failed to validate key: %w", err)
	}

	// Check if key is valid
	if !key.IsValid() {
		if key.IsExpired() {
			return nil, fmt.Errorf("API key has expired")
		}
		return nil, fmt.Errorf("API key is inactive")
	}

	// Update last used timestamp
	key.UpdateLastUsed()
	s.db.Model(&key).Update("last_used_at", key.LastUsedAt)

	return &key, nil
}

// GetKeyByHash retrieves a key by its hash (for internal use)
func (s *KeyService) GetKeyByHash(keyHash string) (*models.ProxyAPIKey, error) {
	var key models.ProxyAPIKey
	if err := s.db.Preload("Provider").Where("key_hash = ?", keyHash).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("key not found")
		}
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return &key, nil
}

// getKeyByID retrieves a key ensuring it belongs to the user
func (s *KeyService) getKeyByID(userID, keyID uint) (*models.ProxyAPIKey, error) {
	var key models.ProxyAPIKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("key not found")
		}
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return &key, nil
}

// buildKeyResponse creates a KeyResponse from a ProxyAPIKey model
func (s *KeyService) buildKeyResponse(key *models.ProxyAPIKey) KeyResponse {
	response := KeyResponse{
		ID:         key.ID,
		UserID:     key.UserID,
		ProviderID: key.ProviderID,
		KeyPrefix:  key.KeyPrefix,
		Name:       key.Name,
		IsActive:   key.IsActive,
		LastUsedAt: key.LastUsedAt,
		ExpiresAt:  key.ExpiresAt,
		CreatedAt:  key.CreatedAt,
		UpdatedAt:  key.UpdatedAt,
	}

	// Include provider info if loaded
	if key.Provider != nil {
		response.ProviderName = key.Provider.Name
		response.ProviderType = key.Provider.ProviderType
	}

	return response
}

// validateCreateRequest validates the create key request
func (s *KeyService) validateCreateRequest(userID uint, req *CreateKeyRequest) error {
	// Validate provider exists and belongs to user
	var provider models.Provider
	if err := s.db.Where("id = ? AND user_id = ?", req.ProviderID, userID).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("provider not found")
		}
		return fmt.Errorf("failed to validate provider: %w", err)
	}

	// Validate provider is active
	if !provider.IsActive {
		return fmt.Errorf("provider is not active")
	}

	// Validate name length if provided
	if len(req.Name) > 100 {
		return fmt.Errorf("name must be 100 characters or less")
	}

	// Validate expiration date if provided
	if req.ExpiresAt != nil && req.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("expiration date must be in the future")
	}

	return nil
}

// generateKey generates a new secure API key with the standard prefix
func (s *KeyService) generateKey() (string, error) {
	// Generate 32 random bytes (256 bits of entropy)
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode as hex and create the full key
	randomPart := hex.EncodeToString(randomBytes)
	fullKey := models.ProxyAPIKeyPrefix + randomPart

	return fullKey, nil
}

// extractKeyPrefix extracts the display prefix from a full key
// Returns format: "sk-smoothllm-abc123...xyz789"
func (s *KeyService) extractKeyPrefix(fullKey string) string {
	// Remove the base prefix to get the random part
	randomPart := strings.TrimPrefix(fullKey, models.ProxyAPIKeyPrefix)

	// Take first 6 and last 4 chars of the random part for display
	if len(randomPart) > 10 {
		return models.ProxyAPIKeyPrefix + randomPart[:6] + "..." + randomPart[len(randomPart)-4:]
	}

	// Fallback for short keys (shouldn't happen with 32-byte random)
	return models.ProxyAPIKeyPrefix + randomPart[:4] + "..."
}

// hashKey creates a SHA256 hash of the API key
func (s *KeyService) hashKey(fullKey string) string {
	hash := sha256.Sum256([]byte(fullKey))
	return hex.EncodeToString(hash[:])
}

// HashKey is a public wrapper for hashing (useful for testing)
func (s *KeyService) HashKey(fullKey string) string {
	return s.hashKey(fullKey)
}
