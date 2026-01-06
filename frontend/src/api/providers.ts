import apiClient from '@/api/client'

// Provider type constants
export const ProviderType = {
  OPENAI: 'openai',
  ANTHROPIC: 'anthropic',
  ANTHROPIC_MAX: 'anthropic_max',
  VLLM: 'vllm',
  LOCAL: 'local',
  ZAI: 'zai',
} as const

export type ProviderTypeValue = (typeof ProviderType)[keyof typeof ProviderType]

export interface ProviderResponse {
  id: number
  user_id: number
  name: string
  provider_type: ProviderTypeValue
  base_url: string
  is_active: boolean
  default_model: string
  input_cost_per_million: number
  output_cost_per_million: number
  oauth_connected: boolean
  created_at: string
  updated_at: string
}

export interface CreateProviderRequest {
  name: string
  provider_type: ProviderTypeValue
  base_url?: string
  api_key?: string // Optional for OAuth providers
  is_active?: boolean
  default_model?: string
  input_cost_per_million?: number
  output_cost_per_million?: number
}

export interface UpdateProviderRequest {
  name?: string
  provider_type?: ProviderTypeValue
  base_url?: string
  api_key?: string
  is_active?: boolean
  default_model?: string
  input_cost_per_million?: number
  output_cost_per_million?: number
}

export interface TestConnectionResponse {
  message: string
}

export interface OAuthAuthorizeResponse {
  authorization_url: string
}

export const providersApi = {
  async listProviders(): Promise<ProviderResponse[]> {
    const response = await apiClient.get('/providers')
    return response.data
  },

  async getProvider(id: number): Promise<ProviderResponse> {
    const response = await apiClient.get(`/providers/${id}`)
    return response.data
  },

  async createProvider(payload: CreateProviderRequest): Promise<ProviderResponse> {
    const response = await apiClient.post('/providers', payload)
    return response.data
  },

  async updateProvider(id: number, payload: UpdateProviderRequest): Promise<ProviderResponse> {
    const response = await apiClient.put(`/providers/${id}`, payload)
    return response.data
  },

  async deleteProvider(id: number): Promise<void> {
    await apiClient.delete(`/providers/${id}`)
  },

  async testConnection(id: number): Promise<TestConnectionResponse> {
    const response = await apiClient.post(`/providers/${id}/test`)
    return response.data
  },

  async testConnectionWithCredentials(
    payload: CreateProviderRequest
  ): Promise<TestConnectionResponse> {
    const response = await apiClient.post('/providers/test', payload)
    return response.data
  },

  // OAuth methods for Claude Max
  async getOAuthAuthorizeUrl(providerId: number): Promise<OAuthAuthorizeResponse> {
    const response = await apiClient.get(`/oauth/anthropic/authorize?provider_id=${providerId}`)
    return response.data
  },

  async disconnectOAuth(providerId: number): Promise<{ message: string }> {
    const response = await apiClient.post(`/oauth/anthropic/disconnect/${providerId}`)
    return response.data
  },

  async testOAuthConnection(providerId: number): Promise<TestConnectionResponse> {
    const response = await apiClient.post(`/oauth/anthropic/test/${providerId}`)
    return response.data
  },
}
