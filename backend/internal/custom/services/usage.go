package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/smoothweb/backend/internal/custom/models"
)

// UsageService handles usage recording and querying operations
type UsageService struct {
	db *gorm.DB
}

// NewUsageService creates a new UsageService instance
func NewUsageService(db *gorm.DB) *UsageService {
	return &UsageService{db: db}
}

// UsageSummaryResponse represents the overall usage summary for a user
type UsageSummaryResponse struct {
	TotalRequests      int64   `json:"total_requests"`
	SuccessfulRequests int64   `json:"successful_requests"`
	FailedRequests     int64   `json:"failed_requests"`
	TotalInputTokens   int64   `json:"total_input_tokens"`
	TotalOutputTokens  int64   `json:"total_output_tokens"`
	TotalTokens        int64   `json:"total_tokens"`
	TotalCost          float64 `json:"total_cost"`
	AverageDuration    float64 `json:"average_duration_ms"`
	PeriodStart        string  `json:"period_start"`
	PeriodEnd          string  `json:"period_end"`
}

// DailyUsageResponse represents usage data for a single day
type DailyUsageResponse struct {
	Date           string  `json:"date"`
	Requests       int64   `json:"requests"`
	InputTokens    int64   `json:"input_tokens"`
	OutputTokens   int64   `json:"output_tokens"`
	TotalTokens    int64   `json:"total_tokens"`
	Cost           float64 `json:"cost"`
	AverageDuration float64 `json:"average_duration_ms"`
}

// UsageByKeyResponse represents usage data grouped by proxy key
type UsageByKeyResponse struct {
	KeyID           uint    `json:"key_id"`
	KeyPrefix       string  `json:"key_prefix"`
	KeyName         string  `json:"key_name"`
	Requests        int64   `json:"requests"`
	InputTokens     int64   `json:"input_tokens"`
	OutputTokens    int64   `json:"output_tokens"`
	TotalTokens     int64   `json:"total_tokens"`
	Cost            float64 `json:"cost"`
	AverageDuration float64 `json:"average_duration_ms"`
}

// UsageByProviderResponse represents usage data grouped by provider
type UsageByProviderResponse struct {
	ProviderID      uint    `json:"provider_id"`
	ProviderName    string  `json:"provider_name"`
	ProviderType    string  `json:"provider_type"`
	Requests        int64   `json:"requests"`
	InputTokens     int64   `json:"input_tokens"`
	OutputTokens    int64   `json:"output_tokens"`
	TotalTokens     int64   `json:"total_tokens"`
	Cost            float64 `json:"cost"`
	AverageDuration float64 `json:"average_duration_ms"`
}

// UsageByModelResponse represents usage data grouped by model
type UsageByModelResponse struct {
	Model           string  `json:"model"`
	Requests        int64   `json:"requests"`
	InputTokens     int64   `json:"input_tokens"`
	OutputTokens    int64   `json:"output_tokens"`
	TotalTokens     int64   `json:"total_tokens"`
	Cost            float64 `json:"cost"`
	AverageDuration float64 `json:"average_duration_ms"`
}

// UsageRecordResponse represents a single usage record
type UsageRecordResponse struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	ProxyKeyID      uint      `json:"proxy_key_id"`
	ProviderID      uint      `json:"provider_id"`
	Model           string    `json:"model"`
	InputTokens     int       `json:"input_tokens"`
	OutputTokens    int       `json:"output_tokens"`
	TotalTokens     int       `json:"total_tokens"`
	Cost            float64   `json:"cost"`
	RequestDuration int       `json:"request_duration_ms"`
	StatusCode      int       `json:"status_code"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	// Related info for convenience
	KeyPrefix    string `json:"key_prefix,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	ProviderType string `json:"provider_type,omitempty"`
}

// RecordUsageRequest represents the data needed to record a usage event
type RecordUsageRequest struct {
	UserID               uint
	ProxyKeyID           uint
	ProviderID           uint
	Model                string
	InputTokens          int
	OutputTokens         int
	TotalTokens          int
	RequestDuration      int // milliseconds
	StatusCode           int
	ErrorMessage         string
	InputCostPerMillion  float64
	OutputCostPerMillion float64
}

// UsageQueryParams represents query parameters for filtering usage data
type UsageQueryParams struct {
	StartDate  *time.Time
	EndDate    *time.Time
	ProviderID *uint
	KeyID      *uint
	Model      *string
	Limit      int
	Offset     int
}

// RecordUsage records a new usage event
func (s *UsageService) RecordUsage(req *RecordUsageRequest) (*models.UsageRecord, error) {
	// Calculate total tokens if not provided
	totalTokens := req.TotalTokens
	if totalTokens == 0 {
		totalTokens = req.InputTokens + req.OutputTokens
	}

	record := &models.UsageRecord{
		UserID:          req.UserID,
		ProxyKeyID:      req.ProxyKeyID,
		ProviderID:      req.ProviderID,
		ModelName:       req.Model,
		InputTokens:     req.InputTokens,
		OutputTokens:    req.OutputTokens,
		TotalTokens:     totalTokens,
		RequestDuration: req.RequestDuration,
		StatusCode:      req.StatusCode,
		ErrorMessage:    req.ErrorMessage,
	}

	// Calculate cost based on provider rates (cost per million tokens)
	record.Cost = record.CalculateCost(req.InputCostPerMillion, req.OutputCostPerMillion)

	if err := s.db.Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to record usage: %w", err)
	}

	return record, nil
}

// RecordUsageAsync records usage asynchronously (non-blocking)
func (s *UsageService) RecordUsageAsync(req *RecordUsageRequest) {
	go func() {
		_, err := s.RecordUsage(req)
		if err != nil {
			// Log the error but don't block the response
			// In production, you'd use a proper logging framework
			fmt.Printf("Failed to record usage: %v\n", err)
		}
	}()
}

// GetUsageSummary returns the overall usage summary for a user
func (s *UsageService) GetUsageSummary(userID uint, params *UsageQueryParams) (*UsageSummaryResponse, error) {
	query := s.db.Model(&models.UsageRecord{}).Where("user_id = ?", userID)
	query = s.applyFilters(query, params)

	var result struct {
		TotalRequests      int64
		SuccessfulRequests int64
		FailedRequests     int64
		TotalInputTokens   int64
		TotalOutputTokens  int64
		TotalTokens        int64
		TotalCost          float64
		TotalDuration      int64
		MinDate            time.Time
		MaxDate            time.Time
	}

	// Get aggregate stats
	if err := query.Select(`
		COUNT(*) as total_requests,
		SUM(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 ELSE 0 END) as successful_requests,
		SUM(CASE WHEN status_code >= 400 OR error_message != '' THEN 1 ELSE 0 END) as failed_requests,
		COALESCE(SUM(input_tokens), 0) as total_input_tokens,
		COALESCE(SUM(output_tokens), 0) as total_output_tokens,
		COALESCE(SUM(total_tokens), 0) as total_tokens,
		COALESCE(SUM(cost), 0) as total_cost,
		COALESCE(SUM(request_duration), 0) as total_duration,
		MIN(created_at) as min_date,
		MAX(created_at) as max_date
	`).Scan(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to get usage summary: %w", err)
	}

	// Calculate average duration
	var avgDuration float64
	if result.TotalRequests > 0 {
		avgDuration = float64(result.TotalDuration) / float64(result.TotalRequests)
	}

	// Format dates
	periodStart := ""
	periodEnd := ""
	if !result.MinDate.IsZero() {
		periodStart = result.MinDate.Format(time.RFC3339)
	}
	if !result.MaxDate.IsZero() {
		periodEnd = result.MaxDate.Format(time.RFC3339)
	}

	return &UsageSummaryResponse{
		TotalRequests:      result.TotalRequests,
		SuccessfulRequests: result.SuccessfulRequests,
		FailedRequests:     result.FailedRequests,
		TotalInputTokens:   result.TotalInputTokens,
		TotalOutputTokens:  result.TotalOutputTokens,
		TotalTokens:        result.TotalTokens,
		TotalCost:          result.TotalCost,
		AverageDuration:    avgDuration,
		PeriodStart:        periodStart,
		PeriodEnd:          periodEnd,
	}, nil
}

// GetDailyUsage returns usage data broken down by day
func (s *UsageService) GetDailyUsage(userID uint, params *UsageQueryParams) ([]DailyUsageResponse, error) {
	query := s.db.Model(&models.UsageRecord{}).Where("user_id = ?", userID)
	query = s.applyFilters(query, params)

	var results []struct {
		Date          string
		Requests      int64
		InputTokens   int64
		OutputTokens  int64
		TotalTokens   int64
		Cost          float64
		TotalDuration int64
	}

	// Group by date (using SQLite date function)
	if err := query.Select(`
		DATE(created_at) as date,
		COUNT(*) as requests,
		COALESCE(SUM(input_tokens), 0) as input_tokens,
		COALESCE(SUM(output_tokens), 0) as output_tokens,
		COALESCE(SUM(total_tokens), 0) as total_tokens,
		COALESCE(SUM(cost), 0) as cost,
		COALESCE(SUM(request_duration), 0) as total_duration
	`).Group("DATE(created_at)").Order("date DESC").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get daily usage: %w", err)
	}

	responses := make([]DailyUsageResponse, len(results))
	for i, r := range results {
		var avgDuration float64
		if r.Requests > 0 {
			avgDuration = float64(r.TotalDuration) / float64(r.Requests)
		}
		responses[i] = DailyUsageResponse{
			Date:            r.Date,
			Requests:        r.Requests,
			InputTokens:     r.InputTokens,
			OutputTokens:    r.OutputTokens,
			TotalTokens:     r.TotalTokens,
			Cost:            r.Cost,
			AverageDuration: avgDuration,
		}
	}

	return responses, nil
}

// GetUsageByKey returns usage data grouped by proxy key
func (s *UsageService) GetUsageByKey(userID uint, params *UsageQueryParams) ([]UsageByKeyResponse, error) {
	query := s.db.Model(&models.UsageRecord{}).
		Joins("LEFT JOIN proxy_api_keys ON usage_records.proxy_key_id = proxy_api_keys.id").
		Where("usage_records.user_id = ?", userID)
	query = s.applyFilters(query, params)

	var results []struct {
		KeyID         uint
		KeyPrefix     string
		KeyName       string
		Requests      int64
		InputTokens   int64
		OutputTokens  int64
		TotalTokens   int64
		Cost          float64
		TotalDuration int64
	}

	if err := query.Select(`
		usage_records.proxy_key_id as key_id,
		proxy_api_keys.key_prefix as key_prefix,
		proxy_api_keys.name as key_name,
		COUNT(*) as requests,
		COALESCE(SUM(usage_records.input_tokens), 0) as input_tokens,
		COALESCE(SUM(usage_records.output_tokens), 0) as output_tokens,
		COALESCE(SUM(usage_records.total_tokens), 0) as total_tokens,
		COALESCE(SUM(usage_records.cost), 0) as cost,
		COALESCE(SUM(usage_records.request_duration), 0) as total_duration
	`).Group("usage_records.proxy_key_id").Order("requests DESC").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get usage by key: %w", err)
	}

	responses := make([]UsageByKeyResponse, len(results))
	for i, r := range results {
		var avgDuration float64
		if r.Requests > 0 {
			avgDuration = float64(r.TotalDuration) / float64(r.Requests)
		}
		responses[i] = UsageByKeyResponse{
			KeyID:           r.KeyID,
			KeyPrefix:       r.KeyPrefix,
			KeyName:         r.KeyName,
			Requests:        r.Requests,
			InputTokens:     r.InputTokens,
			OutputTokens:    r.OutputTokens,
			TotalTokens:     r.TotalTokens,
			Cost:            r.Cost,
			AverageDuration: avgDuration,
		}
	}

	return responses, nil
}

// GetUsageByProvider returns usage data grouped by provider
func (s *UsageService) GetUsageByProvider(userID uint, params *UsageQueryParams) ([]UsageByProviderResponse, error) {
	query := s.db.Model(&models.UsageRecord{}).
		Joins("LEFT JOIN providers ON usage_records.provider_id = providers.id").
		Where("usage_records.user_id = ?", userID)
	query = s.applyFilters(query, params)

	var results []struct {
		ProviderID    uint
		ProviderName  string
		ProviderType  string
		Requests      int64
		InputTokens   int64
		OutputTokens  int64
		TotalTokens   int64
		Cost          float64
		TotalDuration int64
	}

	if err := query.Select(`
		usage_records.provider_id as provider_id,
		providers.name as provider_name,
		providers.provider_type as provider_type,
		COUNT(*) as requests,
		COALESCE(SUM(usage_records.input_tokens), 0) as input_tokens,
		COALESCE(SUM(usage_records.output_tokens), 0) as output_tokens,
		COALESCE(SUM(usage_records.total_tokens), 0) as total_tokens,
		COALESCE(SUM(usage_records.cost), 0) as cost,
		COALESCE(SUM(usage_records.request_duration), 0) as total_duration
	`).Group("usage_records.provider_id").Order("requests DESC").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get usage by provider: %w", err)
	}

	responses := make([]UsageByProviderResponse, len(results))
	for i, r := range results {
		var avgDuration float64
		if r.Requests > 0 {
			avgDuration = float64(r.TotalDuration) / float64(r.Requests)
		}
		responses[i] = UsageByProviderResponse{
			ProviderID:      r.ProviderID,
			ProviderName:    r.ProviderName,
			ProviderType:    r.ProviderType,
			Requests:        r.Requests,
			InputTokens:     r.InputTokens,
			OutputTokens:    r.OutputTokens,
			TotalTokens:     r.TotalTokens,
			Cost:            r.Cost,
			AverageDuration: avgDuration,
		}
	}

	return responses, nil
}

// GetUsageByModel returns usage data grouped by model
func (s *UsageService) GetUsageByModel(userID uint, params *UsageQueryParams) ([]UsageByModelResponse, error) {
	query := s.db.Model(&models.UsageRecord{}).Where("user_id = ?", userID)
	query = s.applyFilters(query, params)

	var results []struct {
		Model         string
		Requests      int64
		InputTokens   int64
		OutputTokens  int64
		TotalTokens   int64
		Cost          float64
		TotalDuration int64
	}

	if err := query.Select(`
		model,
		COUNT(*) as requests,
		COALESCE(SUM(input_tokens), 0) as input_tokens,
		COALESCE(SUM(output_tokens), 0) as output_tokens,
		COALESCE(SUM(total_tokens), 0) as total_tokens,
		COALESCE(SUM(cost), 0) as cost,
		COALESCE(SUM(request_duration), 0) as total_duration
	`).Group("model").Order("requests DESC").Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get usage by model: %w", err)
	}

	responses := make([]UsageByModelResponse, len(results))
	for i, r := range results {
		var avgDuration float64
		if r.Requests > 0 {
			avgDuration = float64(r.TotalDuration) / float64(r.Requests)
		}
		responses[i] = UsageByModelResponse{
			Model:           r.Model,
			Requests:        r.Requests,
			InputTokens:     r.InputTokens,
			OutputTokens:    r.OutputTokens,
			TotalTokens:     r.TotalTokens,
			Cost:            r.Cost,
			AverageDuration: avgDuration,
		}
	}

	return responses, nil
}

// GetRecentUsage returns recent usage records with details
func (s *UsageService) GetRecentUsage(userID uint, params *UsageQueryParams) ([]UsageRecordResponse, error) {
	query := s.db.Model(&models.UsageRecord{}).
		Preload("ProxyKey").
		Preload("Provider").
		Where("user_id = ?", userID)
	query = s.applyFilters(query, params)

	// Apply pagination
	limit := 50 // default
	if params != nil && params.Limit > 0 {
		limit = params.Limit
	}
	if limit > 100 {
		limit = 100 // max
	}
	offset := 0
	if params != nil && params.Offset > 0 {
		offset = params.Offset
	}

	var records []models.UsageRecord
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent usage: %w", err)
	}

	responses := make([]UsageRecordResponse, len(records))
	for i, record := range records {
		responses[i] = s.buildUsageRecordResponse(&record)
	}

	return responses, nil
}

// buildUsageRecordResponse creates a UsageRecordResponse from a UsageRecord model
func (s *UsageService) buildUsageRecordResponse(record *models.UsageRecord) UsageRecordResponse {
	response := UsageRecordResponse{
		ID:              record.ID,
		UserID:          record.UserID,
		ProxyKeyID:      record.ProxyKeyID,
		ProviderID:      record.ProviderID,
		Model:           record.ModelName,
		InputTokens:     record.InputTokens,
		OutputTokens:    record.OutputTokens,
		TotalTokens:     record.TotalTokens,
		Cost:            record.Cost,
		RequestDuration: record.RequestDuration,
		StatusCode:      record.StatusCode,
		ErrorMessage:    record.ErrorMessage,
		CreatedAt:       record.CreatedAt,
	}

	// Include related info if loaded
	if record.ProxyKey != nil {
		response.KeyPrefix = record.ProxyKey.KeyPrefix
	}
	if record.Provider != nil {
		response.ProviderName = record.Provider.Name
		response.ProviderType = record.Provider.ProviderType
	}

	return response
}

// applyFilters applies query parameters to the database query
func (s *UsageService) applyFilters(query *gorm.DB, params *UsageQueryParams) *gorm.DB {
	if params == nil {
		return query
	}

	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}
	if params.ProviderID != nil {
		query = query.Where("provider_id = ?", *params.ProviderID)
	}
	if params.KeyID != nil {
		query = query.Where("proxy_key_id = ?", *params.KeyID)
	}
	if params.Model != nil && *params.Model != "" {
		query = query.Where("model = ?", *params.Model)
	}

	return query
}

// GetUsageCount returns the total count of usage records for a user (useful for pagination)
func (s *UsageService) GetUsageCount(userID uint, params *UsageQueryParams) (int64, error) {
	query := s.db.Model(&models.UsageRecord{}).Where("user_id = ?", userID)
	query = s.applyFilters(query, params)

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get usage count: %w", err)
	}

	return count, nil
}
