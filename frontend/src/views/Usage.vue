<template>
  <AppLayout>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-display text-text-primary mb-2">Usage Statistics</h1>
        <p class="text-text-muted">Track your LLM API usage and costs</p>
      </div>
      <Button variant="ghost" @click="handleRefresh" :disabled="usageStore.loading">
        <svg
          class="w-5 h-5"
          :class="{ 'animate-spin': usageStore.loading }"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </Button>
    </div>

    <!-- Date Range Filter -->
    <div class="card mb-6">
      <div class="flex flex-wrap items-end gap-4">
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">Start Date</label>
          <input
            v-model="filters.start_date"
            type="date"
            class="font-sans bg-bg-secondary border border-border-default rounded-md text-text-primary px-4 py-2 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all duration-200"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">End Date</label>
          <input
            v-model="filters.end_date"
            type="date"
            class="font-sans bg-bg-secondary border border-border-default rounded-md text-text-primary px-4 py-2 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all duration-200"
          />
        </div>
        <Button variant="primary" @click="applyFilters" :disabled="usageStore.loading">
          Apply Filters
        </Button>
        <Button variant="ghost" @click="clearFilters" :disabled="usageStore.loading">
          Clear
        </Button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="usageStore.loading && !usageStore.isInitialized" class="flex justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
    </div>

    <template v-else>
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="card">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-primary-500/10 rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <p class="text-text-tertiary text-sm font-medium">Total Requests</p>
              <p class="text-2xl font-display font-semibold text-text-primary">{{ formatNumber(usageStore.totalRequests) }}</p>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-secondary-500/10 rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-secondary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
              </svg>
            </div>
            <div>
              <p class="text-text-tertiary text-sm font-medium">Total Tokens</p>
              <p class="text-2xl font-display font-semibold text-text-primary">{{ formatNumber(usageStore.totalTokens) }}</p>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-success-500/10 rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-text-tertiary text-sm font-medium">Total Cost</p>
              <p class="text-2xl font-display font-semibold text-text-primary">${{ formatCurrency(usageStore.totalCost) }}</p>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 bg-primary-500/10 rounded-lg flex items-center justify-center">
              <svg class="w-6 h-6 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div>
              <p class="text-text-tertiary text-sm font-medium">Success Rate</p>
              <p class="text-2xl font-display font-semibold text-text-primary">{{ formatPercent(usageStore.successRate) }}%</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Token Breakdown -->
      <div v-if="usageStore.summary" class="card mb-8">
        <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">Token Breakdown</h2>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <div class="text-center bg-bg-tertiary rounded-lg" style="padding: var(--space-4)">
            <div class="text-primary-500 font-display font-semibold text-2xl mb-2">{{ formatNumber(usageStore.summary.total_input_tokens) }}</div>
            <p class="text-text-tertiary text-sm">Input Tokens</p>
          </div>
          <div class="text-center bg-bg-tertiary rounded-lg" style="padding: var(--space-4)">
            <div class="text-secondary-500 font-display font-semibold text-2xl mb-2">{{ formatNumber(usageStore.summary.total_output_tokens) }}</div>
            <p class="text-text-tertiary text-sm">Output Tokens</p>
          </div>
          <div class="text-center bg-bg-tertiary rounded-lg" style="padding: var(--space-4)">
            <div class="text-success-500 font-display font-semibold text-2xl mb-2">{{ formatDuration(usageStore.summary.average_duration_ms) }}</div>
            <p class="text-text-tertiary text-sm">Avg Response Time</p>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-if="!usageStore.hasUsageData && usageStore.isInitialized"
        class="bg-bg-secondary border border-border-subtle rounded-lg p-12 text-center"
      >
        <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <h3 class="text-lg font-display text-text-primary mb-2">No Usage Data</h3>
        <p class="text-text-muted mb-6">Start using the LLM proxy to see your usage statistics here.</p>
        <Button variant="primary" @click="router.push('/keys')">Create API Key</Button>
      </div>

      <!-- Usage by Provider -->
      <div v-if="usageStore.usageByProvider.length > 0" class="card mb-8">
        <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">Usage by Provider</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border-subtle">
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Provider</th>
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Type</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Requests</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Cost</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Avg Duration</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="provider in usageStore.usageByProvider"
                :key="provider.provider_id"
                class="border-b border-border-subtle last:border-b-0 hover:bg-bg-tertiary/50 transition-colors"
              >
                <td class="py-3 px-4 text-text-primary font-medium">{{ provider.provider_name }}</td>
                <td class="py-3 px-4 text-text-secondary font-mono text-sm uppercase">{{ provider.provider_type }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(provider.requests) }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(provider.total_tokens) }}</td>
                <td class="py-3 px-4 text-success-500 text-right font-mono">${{ formatCurrency(provider.cost) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatDuration(provider.average_duration_ms) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Usage by Model -->
      <div v-if="usageStore.usageByModel.length > 0" class="card mb-8">
        <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">Usage by Model</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border-subtle">
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Model</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Requests</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Input Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Output Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Cost</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Avg Duration</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="model in usageStore.usageByModel"
                :key="model.model"
                class="border-b border-border-subtle last:border-b-0 hover:bg-bg-tertiary/50 transition-colors"
              >
                <td class="py-3 px-4 text-text-primary font-mono">{{ model.model }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(model.requests) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatNumber(model.input_tokens) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatNumber(model.output_tokens) }}</td>
                <td class="py-3 px-4 text-success-500 text-right font-mono">${{ formatCurrency(model.cost) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatDuration(model.average_duration_ms) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Usage by API Key -->
      <div v-if="usageStore.usageByKey.length > 0" class="card mb-8">
        <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">Usage by API Key</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border-subtle">
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Key Name</th>
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Key Prefix</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Requests</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Cost</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Avg Duration</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="key in usageStore.usageByKey"
                :key="key.key_id"
                class="border-b border-border-subtle last:border-b-0 hover:bg-bg-tertiary/50 transition-colors"
              >
                <td class="py-3 px-4 text-text-primary font-medium">{{ key.key_name || 'Unnamed Key' }}</td>
                <td class="py-3 px-4 text-text-secondary font-mono text-sm">{{ key.key_prefix }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(key.requests) }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(key.total_tokens) }}</td>
                <td class="py-3 px-4 text-success-500 text-right font-mono">${{ formatCurrency(key.cost) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatDuration(key.average_duration_ms) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Daily Usage -->
      <div v-if="usageStore.dailyUsage.length > 0" class="card mb-8">
        <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">Daily Usage (Last 30 Days)</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border-subtle">
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Date</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Requests</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Input Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Output Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Cost</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="day in usageStore.dailyUsage.slice(0, 30)"
                :key="day.date"
                class="border-b border-border-subtle last:border-b-0 hover:bg-bg-tertiary/50 transition-colors"
              >
                <td class="py-3 px-4 text-text-primary font-mono">{{ formatDateShort(day.date) }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(day.requests) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatNumber(day.input_tokens) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatNumber(day.output_tokens) }}</td>
                <td class="py-3 px-4 text-success-500 text-right font-mono">${{ formatCurrency(day.cost) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Recent Usage Records -->
      <div v-if="usageStore.recentRecords.length > 0" class="card">
        <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">Recent Requests</h2>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border-subtle">
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Time</th>
                <th class="text-left py-3 px-4 text-text-tertiary font-medium text-sm">Model</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Tokens</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Cost</th>
                <th class="text-right py-3 px-4 text-text-tertiary font-medium text-sm">Duration</th>
                <th class="text-center py-3 px-4 text-text-tertiary font-medium text-sm">Status</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="record in usageStore.recentRecords.slice(0, 20)"
                :key="record.id"
                class="border-b border-border-subtle last:border-b-0 hover:bg-bg-tertiary/50 transition-colors"
              >
                <td class="py-3 px-4 text-text-primary font-mono text-sm">{{ formatDateTime(record.created_at) }}</td>
                <td class="py-3 px-4 text-text-primary font-mono text-sm">{{ record.model }}</td>
                <td class="py-3 px-4 text-text-primary text-right font-mono">{{ formatNumber(record.total_tokens) }}</td>
                <td class="py-3 px-4 text-success-500 text-right font-mono">${{ formatCurrency(record.cost) }}</td>
                <td class="py-3 px-4 text-text-secondary text-right font-mono">{{ formatDuration(record.request_duration_ms) }}</td>
                <td class="py-3 px-4 text-center">
                  <span
                    :class="[
                      'px-2 py-0.5 rounded-full text-xs font-medium',
                      record.status_code >= 200 && record.status_code < 300
                        ? 'bg-success-500/10 text-success-500'
                        : 'bg-error-500/10 text-error-500'
                    ]"
                  >
                    {{ record.status_code }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useUsageStore } from '@/stores/usage'
import AppLayout from '@/components/layout/AppLayout.vue'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const usageStore = useUsageStore()

// Filters
const filters = ref({
  start_date: '',
  end_date: '',
})

// Formatting helpers
const formatNumber = (num: number): string => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toLocaleString()
}

const formatCurrency = (amount: number): string => {
  return amount.toFixed(4)
}

const formatPercent = (value: number): string => {
  return value.toFixed(1)
}

const formatDuration = (ms: number): string => {
  if (ms >= 1000) {
    return (ms / 1000).toFixed(2) + 's'
  }
  return ms.toFixed(0) + 'ms'
}

const formatDateShort = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
  })
}

const formatDateTime = (dateString: string): string => {
  return new Date(dateString).toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

// Actions
const applyFilters = async () => {
  try {
    await usageStore.fetchAllUsageData({
      start_date: filters.value.start_date || undefined,
      end_date: filters.value.end_date || undefined,
    })
  } catch {
    toast.error('Failed to fetch usage data')
  }
}

const clearFilters = async () => {
  filters.value = {
    start_date: '',
    end_date: '',
  }
  await applyFilters()
}

const handleRefresh = async () => {
  try {
    await usageStore.refreshUsageData()
    toast.success('Usage data refreshed')
  } catch {
    toast.error('Failed to refresh usage data')
  }
}

// Load data on mount
onMounted(async () => {
  try {
    await usageStore.fetchAllUsageData()
  } catch {
    toast.error('Failed to load usage data')
  }
})
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
