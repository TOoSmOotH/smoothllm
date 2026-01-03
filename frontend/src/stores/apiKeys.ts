import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  keysApi,
  type KeyResponse,
  type KeyCreateResponse,
  type CreateKeyRequest,
  type UpdateKeyRequest,
} from '@/api/keys'

export const useApiKeysStore = defineStore('apiKeys', () => {
  // State
  const apiKeys = ref<KeyResponse[]>([])
  const selectedKey = ref<KeyResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isInitialized = ref(false)

  // Getters
  const hasApiKeys = computed(() => apiKeys.value.length > 0)
  const activeApiKeys = computed(() => apiKeys.value.filter((k) => k.is_active))
  const apiKeyCount = computed(() => apiKeys.value.length)
  const activeApiKeyCount = computed(() => activeApiKeys.value.length)

  const getApiKeyById = computed(() => {
    return (id: number) => apiKeys.value.find((k) => k.id === id)
  })

  const getApiKeysByProviderId = computed(() => {
    return (providerId: number) => apiKeys.value.filter((k) => k.provider_id === providerId)
  })

  // Actions
  const fetchApiKeys = async (): Promise<KeyResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await keysApi.listKeys()
      apiKeys.value = data
      isInitialized.value = true
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch API keys'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchApiKey = async (id: number): Promise<KeyResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await keysApi.getKey(id)
      selectedKey.value = data

      // Update in list if exists
      const index = apiKeys.value.findIndex((k) => k.id === id)
      if (index !== -1) {
        apiKeys.value[index] = data
      }

      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch API key'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createApiKey = async (payload: CreateKeyRequest): Promise<KeyCreateResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await keysApi.createKey(payload)
      // Add to list (without the full key, just the response info)
      apiKeys.value.push(data)
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create API key'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateApiKey = async (
    id: number,
    payload: UpdateKeyRequest
  ): Promise<KeyResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await keysApi.updateKey(id, payload)

      // Update in list
      const index = apiKeys.value.findIndex((k) => k.id === id)
      if (index !== -1) {
        apiKeys.value[index] = data
      }

      // Update selected if it's the same
      if (selectedKey.value?.id === id) {
        selectedKey.value = data
      }

      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update API key'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteApiKey = async (id: number): Promise<void> => {
    loading.value = true
    error.value = null

    try {
      await keysApi.deleteKey(id)

      // Remove from list
      apiKeys.value = apiKeys.value.filter((k) => k.id !== id)

      // Clear selected if it was deleted
      if (selectedKey.value?.id === id) {
        selectedKey.value = null
      }
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete API key'
      throw err
    } finally {
      loading.value = false
    }
  }

  const revokeApiKey = async (id: number): Promise<string> => {
    loading.value = true
    error.value = null

    try {
      const response = await keysApi.revokeKey(id)

      // Update the key in the list to mark as inactive
      const index = apiKeys.value.findIndex((k) => k.id === id)
      if (index !== -1) {
        apiKeys.value[index] = { ...apiKeys.value[index], is_active: false }
      }

      // Update selected if it's the same
      if (selectedKey.value?.id === id) {
        selectedKey.value = { ...selectedKey.value, is_active: false }
      }

      return response.message
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to revoke API key'
      throw err
    } finally {
      loading.value = false
    }
  }

  const selectApiKey = (key: KeyResponse | null) => {
    selectedKey.value = key
  }

  const clearError = () => {
    error.value = null
  }

  const reset = () => {
    apiKeys.value = []
    selectedKey.value = null
    loading.value = false
    error.value = null
    isInitialized.value = false
  }

  return {
    // State
    apiKeys,
    selectedKey,
    loading,
    error,
    isInitialized,

    // Getters
    hasApiKeys,
    activeApiKeys,
    apiKeyCount,
    activeApiKeyCount,
    getApiKeyById,
    getApiKeysByProviderId,

    // Actions
    fetchApiKeys,
    fetchApiKey,
    createApiKey,
    updateApiKey,
    deleteApiKey,
    revokeApiKey,
    selectApiKey,
    clearError,
    reset,
  }
})
