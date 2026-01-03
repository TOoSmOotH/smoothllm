import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  usageApi,
  type UsageQueryParams,
  type UsageSummaryResponse,
  type DailyUsageResponse,
  type UsageByKeyResponse,
  type UsageByProviderResponse,
  type UsageByModelResponse,
  type UsageRecordResponse,
} from '@/api/usage'

export const useUsageStore = defineStore('usage', () => {
  // State
  const summary = ref<UsageSummaryResponse | null>(null)
  const dailyUsage = ref<DailyUsageResponse[]>([])
  const usageByKey = ref<UsageByKeyResponse[]>([])
  const usageByProvider = ref<UsageByProviderResponse[]>([])
  const usageByModel = ref<UsageByModelResponse[]>([])
  const recentRecords = ref<UsageRecordResponse[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isInitialized = ref(false)
  const currentParams = ref<UsageQueryParams | null>(null)

  // Getters
  const hasUsageData = computed(() => summary.value !== null && summary.value.total_requests > 0)
  const totalRequests = computed(() => summary.value?.total_requests ?? 0)
  const totalCost = computed(() => summary.value?.total_cost ?? 0)
  const totalTokens = computed(() => summary.value?.total_tokens ?? 0)
  const successRate = computed(() => {
    if (!summary.value || summary.value.total_requests === 0) return 0
    return (summary.value.successful_requests / summary.value.total_requests) * 100
  })

  const getUsageByKeyId = computed(() => {
    return (keyId: number) => usageByKey.value.find((u) => u.key_id === keyId)
  })

  const getUsageByProviderId = computed(() => {
    return (providerId: number) => usageByProvider.value.find((u) => u.provider_id === providerId)
  })

  const getUsageByModelName = computed(() => {
    return (model: string) => usageByModel.value.find((u) => u.model === model)
  })

  // Actions
  const fetchUsageSummary = async (params?: UsageQueryParams): Promise<UsageSummaryResponse> => {
    loading.value = true
    error.value = null

    try {
      const data = await usageApi.getUsageSummary(params)
      summary.value = data
      currentParams.value = params ?? null
      isInitialized.value = true
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch usage summary'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchDailyUsage = async (params?: UsageQueryParams): Promise<DailyUsageResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await usageApi.getDailyUsage(params)
      dailyUsage.value = data
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch daily usage'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchUsageByKey = async (params?: UsageQueryParams): Promise<UsageByKeyResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await usageApi.getUsageByKey(params)
      usageByKey.value = data
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch usage by key'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchUsageByProvider = async (params?: UsageQueryParams): Promise<UsageByProviderResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await usageApi.getUsageByProvider(params)
      usageByProvider.value = data
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch usage by provider'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchUsageByModel = async (params?: UsageQueryParams): Promise<UsageByModelResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await usageApi.getUsageByModel(params)
      usageByModel.value = data
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch usage by model'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchRecentUsage = async (params?: UsageQueryParams): Promise<UsageRecordResponse[]> => {
    loading.value = true
    error.value = null

    try {
      const data = await usageApi.getRecentUsage(params)
      recentRecords.value = data
      return data
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch recent usage'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchAllUsageData = async (params?: UsageQueryParams): Promise<void> => {
    loading.value = true
    error.value = null

    try {
      const [summaryData, dailyData, byKeyData, byProviderData, byModelData, recentData] =
        await Promise.all([
          usageApi.getUsageSummary(params),
          usageApi.getDailyUsage(params),
          usageApi.getUsageByKey(params),
          usageApi.getUsageByProvider(params),
          usageApi.getUsageByModel(params),
          usageApi.getRecentUsage(params),
        ])

      summary.value = summaryData
      dailyUsage.value = dailyData
      usageByKey.value = byKeyData
      usageByProvider.value = byProviderData
      usageByModel.value = byModelData
      recentRecords.value = recentData
      currentParams.value = params ?? null
      isInitialized.value = true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to fetch usage data'
      throw err
    } finally {
      loading.value = false
    }
  }

  const refreshUsageData = async (): Promise<void> => {
    await fetchAllUsageData(currentParams.value ?? undefined)
  }

  const clearError = () => {
    error.value = null
  }

  const reset = () => {
    summary.value = null
    dailyUsage.value = []
    usageByKey.value = []
    usageByProvider.value = []
    usageByModel.value = []
    recentRecords.value = []
    loading.value = false
    error.value = null
    isInitialized.value = false
    currentParams.value = null
  }

  return {
    // State
    summary,
    dailyUsage,
    usageByKey,
    usageByProvider,
    usageByModel,
    recentRecords,
    loading,
    error,
    isInitialized,
    currentParams,

    // Getters
    hasUsageData,
    totalRequests,
    totalCost,
    totalTokens,
    successRate,
    getUsageByKeyId,
    getUsageByProviderId,
    getUsageByModelName,

    // Actions
    fetchUsageSummary,
    fetchDailyUsage,
    fetchUsageByKey,
    fetchUsageByProvider,
    fetchUsageByModel,
    fetchRecentUsage,
    fetchAllUsageData,
    refreshUsageData,
    clearError,
    reset,
  }
})
