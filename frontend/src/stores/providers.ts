import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  providersApi,
  type ProviderResponse,
  type CreateProviderRequest,
  type UpdateProviderRequest,
} from '@/api/providers'

export const useProvidersStore = defineStore('providers', () => {
  // State
  const providers = ref<ProviderResponse[]>([])
  const selectedProvider = ref<ProviderResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isInitialized = ref(false)

  // Getters
  const hasProviders = computed(() => providers.value.length > 0)
  const activeProviders = computed(() => providers.value.filter((p) => p.is_active))
  const providerCount = computed(() => providers.value.length)
  const activeProviderCount = computed(() => activeProviders.value.length)

  const getProviderById = computed(() => {
    return (id: number) => providers.value.find((p) => p.id === id)
  })

  const getProvidersByType = computed(() => {
    return (type: string) => providers.value.filter((p) => p.provider_type === type)
  })

  // Actions
  const fetchProviders = async (): Promise<ProviderResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await providersApi.listProviders()
      providers.value = data
      isInitialized.value = true
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch providers'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchProvider = async (id: number): Promise<ProviderResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await providersApi.getProvider(id)
      selectedProvider.value = data

      // Update in list if exists
      const index = providers.value.findIndex((p) => p.id === id)
      if (index !== -1) {
        providers.value[index] = data
      }

      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch provider'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createProvider = async (payload: CreateProviderRequest): Promise<ProviderResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await providersApi.createProvider(payload)
      providers.value.push(data)
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to create provider'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateProvider = async (
    id: number,
    payload: UpdateProviderRequest
  ): Promise<ProviderResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await providersApi.updateProvider(id, payload)

      // Update in list
      const index = providers.value.findIndex((p) => p.id === id)
      if (index !== -1) {
        providers.value[index] = data
      }

      // Update selected if it's the same
      if (selectedProvider.value?.id === id) {
        selectedProvider.value = data
      }

      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to update provider'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteProvider = async (id: number): Promise<void> => {
    loading.value = true
    error.value = null

    try {
      await providersApi.deleteProvider(id)

      // Remove from list
      providers.value = providers.value.filter((p) => p.id !== id)

      // Clear selected if it was deleted
      if (selectedProvider.value?.id === id) {
        selectedProvider.value = null
      }
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to delete provider'
      throw err
    } finally {
      loading.value = false
    }
  }

  const testConnection = async (id: number): Promise<string> => {
    loading.value = true
    error.value = null

    try {
      const response = await providersApi.testConnection(id)
      return response.message
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Connection test failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const testConnectionWithCredentials = async (payload: CreateProviderRequest): Promise<string> => {
    loading.value = true
    error.value = null

    try {
      const response = await providersApi.testConnectionWithCredentials(payload)
      return response.message
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Connection test failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchAvailableModels = async (id: number): Promise<string[]> => {
    loading.value = true
    error.value = null

    try {
      const response = await providersApi.fetchAvailableModels(id)
      return response.models
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch available models'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchAvailableModelsWithCredentials = async (payload: CreateProviderRequest): Promise<string[]> => {
    loading.value = true
    error.value = null

    try {
      const response = await providersApi.fetchAvailableModelsWithCredentials(payload)
      return response.models
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch available models'
      throw err
    } finally {
      loading.value = false
    }
  }

  const selectProvider = (provider: ProviderResponse | null) => {
    selectedProvider.value = provider
  }

  const clearError = () => {
    error.value = null
  }

  const reset = () => {
    providers.value = []
    selectedProvider.value = null
    loading.value = false
    error.value = null
    isInitialized.value = false
  }

  return {
    // State
    providers,
    selectedProvider,
    loading,
    error,
    isInitialized,

    // Getters
    hasProviders,
    activeProviders,
    providerCount,
    activeProviderCount,
    getProviderById,
    getProvidersByType,

    // Actions
    fetchProviders,
    fetchProvider,
    createProvider,
    updateProvider,
    deleteProvider,
    testConnection,
    testConnectionWithCredentials,
    fetchAvailableModels,
    fetchAvailableModelsWithCredentials,
    selectProvider,
    clearError,
    reset,
  }
})
