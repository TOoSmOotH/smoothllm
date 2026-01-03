import apiClient from '@/api/client'

// Provider type constants
export const ProviderType = {
  OPENAI: 'openai',
  ANTHROPIC: 'anthropic',
  LOCAL: 'local',
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
  input_cost_per_1k: number
  output_cost_per_1k: number
  created_at: string
  updated_at: string
}

export interface CreateProviderRequest {
  name: string
  provider_type: ProviderTypeValue
  base_url?: string
  api_key: string
  is_active?: boolean
  default_model?: string
  input_cost_per_1k?: number
  output_cost_per_1k?: number
}

export interface UpdateProviderRequest {
  name?: string
  provider_type?: ProviderTypeValue
  base_url?: string
  api_key?: string
  is_active?: boolean
  default_model?: string
  input_cost_per_1k?: number
  output_cost_per_1k?: number
}

export interface TestConnectionResponse {
  message: string
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
}
