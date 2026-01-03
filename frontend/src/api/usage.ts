import apiClient from '@/api/client'

// Query parameters for filtering usage data
export interface UsageQueryParams {
  start_date?: string
  end_date?: string
  provider_id?: number
  key_id?: number
  model?: string
  limit?: number
  offset?: number
}

// Overall usage summary for a user
export interface UsageSummaryResponse {
  total_requests: number
  successful_requests: number
  failed_requests: number
  total_input_tokens: number
  total_output_tokens: number
  total_tokens: number
  total_cost: number
  average_duration_ms: number
  period_start: string
  period_end: string
}

// Usage data for a single day
export interface DailyUsageResponse {
  date: string
  requests: number
  input_tokens: number
  output_tokens: number
  total_tokens: number
  cost: number
  average_duration_ms: number
}

// Usage data grouped by proxy key
export interface UsageByKeyResponse {
  key_id: number
  key_prefix: string
  key_name: string
  requests: number
  input_tokens: number
  output_tokens: number
  total_tokens: number
  cost: number
  average_duration_ms: number
}

// Usage data grouped by provider
export interface UsageByProviderResponse {
  provider_id: number
  provider_name: string
  provider_type: string
  requests: number
  input_tokens: number
  output_tokens: number
  total_tokens: number
  cost: number
  average_duration_ms: number
}

// Usage data grouped by model
export interface UsageByModelResponse {
  model: string
  requests: number
  input_tokens: number
  output_tokens: number
  total_tokens: number
  cost: number
  average_duration_ms: number
}

// A single usage record
export interface UsageRecordResponse {
  id: number
  user_id: number
  proxy_key_id: number
  provider_id: number
  model: string
  input_tokens: number
  output_tokens: number
  total_tokens: number
  cost: number
  request_duration_ms: number
  status_code: number
  error_message?: string
  created_at: string
  // Related info for convenience
  key_prefix?: string
  provider_name?: string
  provider_type?: string
}

// Helper to build query string from params
function buildQueryString(params?: UsageQueryParams): string {
  if (!params) return ''

  const queryParts: string[] = []

  if (params.start_date) queryParts.push(`start_date=${encodeURIComponent(params.start_date)}`)
  if (params.end_date) queryParts.push(`end_date=${encodeURIComponent(params.end_date)}`)
  if (params.provider_id) queryParts.push(`provider_id=${params.provider_id}`)
  if (params.key_id) queryParts.push(`key_id=${params.key_id}`)
  if (params.model) queryParts.push(`model=${encodeURIComponent(params.model)}`)
  if (params.limit) queryParts.push(`limit=${params.limit}`)
  if (params.offset) queryParts.push(`offset=${params.offset}`)

  return queryParts.length > 0 ? `?${queryParts.join('&')}` : ''
}

export const usageApi = {
  async getUsageSummary(params?: UsageQueryParams): Promise<UsageSummaryResponse> {
    const response = await apiClient.get(`/usage${buildQueryString(params)}`)
    return response.data
  },

  async getDailyUsage(params?: UsageQueryParams): Promise<DailyUsageResponse[]> {
    const response = await apiClient.get(`/usage/daily${buildQueryString(params)}`)
    return response.data
  },

  async getUsageByKey(params?: UsageQueryParams): Promise<UsageByKeyResponse[]> {
    const response = await apiClient.get(`/usage/by-key${buildQueryString(params)}`)
    return response.data
  },

  async getUsageByProvider(params?: UsageQueryParams): Promise<UsageByProviderResponse[]> {
    const response = await apiClient.get(`/usage/by-provider${buildQueryString(params)}`)
    return response.data
  },

  async getUsageByModel(params?: UsageQueryParams): Promise<UsageByModelResponse[]> {
    const response = await apiClient.get(`/usage/by-model${buildQueryString(params)}`)
    return response.data
  },

  async getRecentUsage(params?: UsageQueryParams): Promise<UsageRecordResponse[]> {
    const response = await apiClient.get(`/usage/recent${buildQueryString(params)}`)
    return response.data
  },
}
