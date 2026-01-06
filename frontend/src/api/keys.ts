import apiClient from '@/api/client'

// Key response from the backend (returned for list/get/update operations)
export interface AllowedProvider {
  provider_id: number
  provider_name: string
  provider_type: string
  models: string[]
}

export interface KeyResponse {
  id: number
  user_id: number
  key_prefix: string
  name: string
  is_active: boolean
  last_used_at: string | null
  expires_at: string | null
  created_at: string
  updated_at: string
  allowed_providers: AllowedProvider[]
}

// Key create response - includes the full key (only returned once on creation)
export interface KeyCreateResponse extends KeyResponse {
  key: string
}

export interface ProviderSelection {
  provider_id: number
  models: string[]
}

export interface CreateKeyRequest {
  name: string
  expires_at?: string
  allowed_providers: ProviderSelection[]
}

export interface UpdateKeyRequest {
  name?: string
  is_active?: boolean
  expires_at?: string
  allowed_providers?: ProviderSelection[]
}

export interface RevokeKeyResponse {
  message: string
}

export const keysApi = {
  async listKeys(): Promise<KeyResponse[]> {
    const response = await apiClient.get('/keys')
    return response.data
  },

  async getKey(id: number): Promise<KeyResponse> {
    const response = await apiClient.get(`/keys/${id}`)
    return response.data
  },

  async createKey(payload: CreateKeyRequest): Promise<KeyCreateResponse> {
    const response = await apiClient.post('/keys', payload)
    return response.data
  },

  async updateKey(id: number, payload: UpdateKeyRequest): Promise<KeyResponse> {
    const response = await apiClient.put(`/keys/${id}`, payload)
    return response.data
  },

  async deleteKey(id: number): Promise<void> {
    await apiClient.delete(`/keys/${id}`)
  },

  async revokeKey(id: number): Promise<RevokeKeyResponse> {
    const response = await apiClient.post(`/keys/${id}/revoke`)
    return response.data
  },
}
