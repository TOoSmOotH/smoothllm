<template>
  <AppLayout>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-display text-text-primary mb-2">LLM Providers</h1>
        <p class="text-text-muted">Configure your LLM provider API keys and settings</p>
      </div>
      <Button variant="primary" @click="openCreateModal">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Provider
      </Button>
    </div>

    <!-- Loading State -->
    <div v-if="providersStore.loading && !providersStore.isInitialized" class="flex justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="!providersStore.hasProviders"
      class="bg-bg-secondary border border-border-subtle rounded-lg p-12 text-center"
    >
      <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
      </div>
      <h3 class="text-lg font-display text-text-primary mb-2">No Providers Configured</h3>
      <p class="text-text-muted mb-6">Add your first LLM provider to start proxying requests.</p>
      <Button variant="primary" @click="openCreateModal">Add Your First Provider</Button>
    </div>

    <!-- Providers List -->
    <div v-else class="space-y-4">
      <div
        v-for="provider in providersStore.providers"
        :key="provider.id"
        class="bg-bg-secondary border border-border-subtle rounded-lg p-6 hover:border-border-default transition-colors duration-200"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h3 class="font-display font-semibold text-lg text-text-primary">{{ provider.name }}</h3>
              <span
                :class="[
                  'px-2 py-0.5 rounded-full text-xs font-medium uppercase',
                  provider.is_active
                    ? 'bg-success-500/10 text-success-500'
                    : 'bg-text-muted/10 text-text-muted'
                ]"
              >
                {{ provider.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-sm">
              <div>
                <p class="text-text-tertiary mb-1">Type</p>
                <p class="text-text-primary font-mono uppercase">{{ provider.provider_type }}</p>
              </div>
              <div>
                <p class="text-text-tertiary mb-1">Default Model</p>
                <p class="text-text-primary font-mono">{{ provider.default_model || 'Not set' }}</p>
              </div>
              <div>
                <p class="text-text-tertiary mb-1">Input Cost</p>
                <p class="text-text-primary font-mono">${{ provider.input_cost_per_million.toFixed(2) }}/M</p>
              </div>
              <div>
                <p class="text-text-tertiary mb-1">Output Cost</p>
                <p class="text-text-primary font-mono">${{ provider.output_cost_per_million.toFixed(2) }}/M</p>
              </div>
            </div>
            <!-- Token Status for Claude Max providers -->
            <div v-if="provider.provider_type === ProviderType.ANTHROPIC_MAX" class="mt-3 pt-3 border-t border-border-subtle">
              <div class="flex items-center gap-2">
                <span
                  :class="[
                    'px-2 py-0.5 rounded-full text-xs font-medium',
                    provider.oauth_connected
                      ? 'bg-success-500/10 text-success-500'
                      : 'bg-warning-500/10 text-warning-500'
                  ]"
                >
                  {{ provider.oauth_connected ? 'Token Valid' : 'Token Not Configured' }}
                </span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2 ml-4">
            <Button
              variant="ghost"
              size="sm"
              @click="provider.provider_type === ProviderType.ANTHROPIC_MAX && provider.oauth_connected ? testOAuthConnection(provider.id) : testProviderConnection(provider.id)"
              :loading="testingId === provider.id"
              :disabled="provider.provider_type === ProviderType.ANTHROPIC_MAX && !provider.oauth_connected"
              :title="provider.provider_type === ProviderType.ANTHROPIC_MAX && !provider.oauth_connected ? 'Connect OAuth first' : 'Test connection'"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </Button>
            <Button variant="ghost" size="sm" @click="openEditModal(provider)">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </Button>
            <Button variant="ghost" size="sm" @click="confirmDelete(provider)">
              <svg class="w-4 h-4 text-error-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </Button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div
        v-if="showModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="closeModal"
      >
        <div class="fixed inset-0 bg-black/50 z-0" @click="closeModal"></div>
        <div class="relative z-10 bg-bg-primary border border-border-subtle rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
          <div class="p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">
              {{ editingProvider ? 'Edit Provider' : 'Add Provider' }}
            </h2>

            <div class="space-y-4">
              <Input
                v-model="form.name"
                label="Provider Name"
                placeholder="My OpenAI Account"
                :error="errors.name"
              />

              <div>
                <label class="block text-sm font-medium text-text-secondary mb-2">Provider Type</label>
                <select
                  v-model="form.provider_type"
                  class="w-full font-sans bg-bg-secondary border border-border-default rounded-md text-text-primary px-4 py-3 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all duration-200"
                >
                  <option value="openai">OpenAI</option>
                  <option value="anthropic">Anthropic (API Key)</option>
                  <option value="anthropic_max">Claude Max (OAuth)</option>
                  <option value="vllm">vLLM</option>
                  <option value="local">Local / Custom</option>
                  <option value="zai">z.ai (Zhipu)</option>
                </select>
                <p v-if="errors.provider_type" class="mt-1 text-xs text-error-500 font-medium">{{ errors.provider_type }}</p>
              </div>

              <Input
                v-model="form.base_url"
                label="Base URL (optional)"
                placeholder="https://api.openai.com/v1"
                :helper-text="getBaseUrlHelperText"
                :error="errors.base_url"
              />

              <!-- API Key input - not shown for OAuth providers -->
              <Input
                v-if="form.provider_type !== ProviderType.ANTHROPIC_MAX"
                v-model="form.api_key"
                type="password"
                :label="editingProvider ? 'API Key (leave empty to keep current)' : 'API Key'"
                placeholder="sk-..."
                :error="errors.api_key"
              />

              <!-- OAuth token input for Claude Max -->
              <div v-if="form.provider_type === ProviderType.ANTHROPIC_MAX" class="space-y-3">
                <div class="bg-primary-500/10 border border-primary-500/20 rounded-md p-4">
                  <p class="text-sm text-text-secondary mb-2">
                    <strong class="text-text-primary">Claude Max uses OAuth tokens from Claude Code CLI.</strong>
                  </p>
                  <p class="text-xs text-text-tertiary mb-2">
                    First, authenticate with Claude Code: <code class="bg-bg-tertiary px-1 rounded">claude auth login</code>
                  </p>
                  <p class="text-xs text-text-tertiary">
                    Then get your refresh token: <code class="bg-bg-tertiary px-1 rounded">jq -r '.claudeAiOauth.refreshToken' ~/.claude/.credentials.json</code>
                  </p>
                </div>
                <Input
                  v-model="form.api_key"
                  type="password"
                  :label="editingProvider ? 'Refresh Token (leave empty to keep current)' : 'Refresh Token'"
                  placeholder="sk-ant-ort01-..."
                  helper-text="Your Claude Max refresh token (starts with sk-ant-ort01-)"
                  :error="errors.api_key"
                />
              </div>

              <Input
                v-model="form.default_model"
                label="Default Model"
                :placeholder="getDefaultModelPlaceholder"
                :error="errors.default_model"
              />

              <div class="grid grid-cols-2 gap-4">
                <Input
                  v-model="form.input_cost_per_million"
                  type="number"
                  label="Input Cost per Million tokens ($)"
                  placeholder="1.50"
                  :error="errors.input_cost_per_million"
                />
                <Input
                  v-model="form.output_cost_per_million"
                  type="number"
                  label="Output Cost per Million tokens ($)"
                  placeholder="2.00"
                  :error="errors.output_cost_per_million"
                />
              </div>

              <div class="flex items-center gap-3">
                <input
                  id="is_active"
                  v-model="form.is_active"
                  type="checkbox"
                  class="w-4 h-4 rounded border-border-default text-primary-500 focus:ring-primary-500"
                />
                <label for="is_active" class="text-sm text-text-secondary">Enable this provider</label>
              </div>
            </div>

            <div class="flex items-center justify-between mt-8 pt-6 border-t border-border-subtle">
              <Button
                variant="outline"
                @click="testConnectionWithForm"
                :loading="testingConnection"
              >
                Test Connection
              </Button>
              <div class="flex items-center gap-3">
                <Button variant="ghost" @click="closeModal">Cancel</Button>
                <Button
                  variant="primary"
                  @click="handleSubmit"
                  :loading="submitting"
                >
                  {{ editingProvider ? 'Save Changes' : 'Create Provider' }}
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div
        v-if="showDeleteModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="closeDeleteModal"
      >
        <div class="fixed inset-0 bg-black/50 z-0" @click="closeDeleteModal"></div>
        <div class="relative z-10 bg-bg-primary border border-border-subtle rounded-lg shadow-xl max-w-md w-full">
          <div class="p-6">
            <div class="flex items-center gap-4 mb-4">
              <div class="w-12 h-12 bg-error-500/10 rounded-full flex items-center justify-center">
                <svg class="w-6 h-6 text-error-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
              <div>
                <h3 class="font-display text-lg text-text-primary">Delete Provider</h3>
                <p class="text-sm text-text-muted">This action cannot be undone.</p>
              </div>
            </div>
            <p class="text-text-secondary mb-6">
              Are you sure you want to delete <strong class="text-text-primary">{{ deletingProvider?.name }}</strong>?
              All associated API keys will also be deleted.
            </p>
            <div class="flex justify-end gap-3">
              <Button variant="ghost" @click="closeDeleteModal">Cancel</Button>
              <Button variant="destructive" @click="handleDelete" :loading="deleting">
                Delete Provider
              </Button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { useProvidersStore } from '@/stores/providers'
import { ProviderType, providersApi, type ProviderResponse, type CreateProviderRequest } from '@/api/providers'
import AppLayout from '@/components/layout/AppLayout.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

const providersStore = useProvidersStore()

// Modal state
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingProvider = ref<ProviderResponse | null>(null)
const deletingProvider = ref<ProviderResponse | null>(null)

// Form state
const form = ref({
  name: '',
  provider_type: ProviderType.OPENAI as string,
  base_url: '',
  api_key: '',
  default_model: '',
  input_cost_per_million: '',
  output_cost_per_million: '',
  is_active: true,
})

const errors = ref<Record<string, string>>({})

// Loading states
const submitting = ref(false)
const deleting = ref(false)
const testingConnection = ref(false)
const testingId = ref<number | null>(null)

// Computed helpers
const getBaseUrlHelperText = computed(() => {
  switch (form.value.provider_type) {
    case ProviderType.OPENAI:
      return 'Leave empty for default: https://api.openai.com/v1'
    case ProviderType.ANTHROPIC:
      return 'Leave empty for default: https://api.anthropic.com'
    case ProviderType.ANTHROPIC_MAX:
      return 'Leave empty for default: https://api.anthropic.com'
    case ProviderType.VLLM:
      return 'Enter your vLLM server URL (e.g., http://localhost:8000)'
    case ProviderType.LOCAL:
      return 'Enter your local model endpoint URL'
    case ProviderType.ZAI:
      return 'Leave empty for default: https://api.z.ai/api/paas/v4/'
    default:
      return ''
  }
})

const getDefaultModelPlaceholder = computed(() => {
  switch (form.value.provider_type) {
    case ProviderType.OPENAI:
      return 'gpt-4o'
    case ProviderType.ANTHROPIC:
      return 'claude-sonnet-4-20250514'
    case ProviderType.ANTHROPIC_MAX:
      return 'claude-sonnet-4-20250514'
    case ProviderType.VLLM:
      return 'your-model-name'
    case ProviderType.LOCAL:
      return 'your-model-name'
    case ProviderType.ZAI:
      return 'GLM-4.7'
    default:
      return ''
  }
})

// Modal functions
const openCreateModal = () => {
  editingProvider.value = null
  resetForm()
  showModal.value = true
}

const openEditModal = (provider: ProviderResponse) => {
  editingProvider.value = provider
  form.value = {
    name: provider.name,
    provider_type: provider.provider_type,
    base_url: provider.base_url || '',
    api_key: '',
    default_model: provider.default_model || '',
    input_cost_per_million: provider.input_cost_per_million.toString(),
    output_cost_per_million: provider.output_cost_per_million.toString(),
    is_active: provider.is_active,
  }
  errors.value = {}
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingProvider.value = null
  resetForm()
}

const confirmDelete = (provider: ProviderResponse) => {
  deletingProvider.value = provider
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deletingProvider.value = null
}

const resetForm = () => {
  form.value = {
    name: '',
    provider_type: ProviderType.OPENAI,
    base_url: '',
    api_key: '',
    default_model: '',
    input_cost_per_million: '',
    output_cost_per_million: '',
    is_active: true,
  }
  errors.value = {}
}

// Validation
const validateForm = (): boolean => {
  errors.value = {}

  if (!form.value.name.trim()) {
    errors.value.name = 'Provider name is required'
  }

  if (!form.value.provider_type) {
    errors.value.provider_type = 'Provider type is required'
  }

  // API key / refresh token validation
  if (!editingProvider.value && !form.value.api_key.trim()) {
    if (form.value.provider_type === ProviderType.ANTHROPIC_MAX) {
      errors.value.api_key = 'Refresh token is required'
    } else {
      errors.value.api_key = 'API key is required'
    }
  }

  if ((form.value.provider_type === ProviderType.LOCAL || form.value.provider_type === ProviderType.VLLM) && !form.value.base_url.trim()) {
    errors.value.base_url = 'Base URL is required for this provider type'
  }

  return Object.keys(errors.value).length === 0
}

// Form submission
const handleSubmit = async () => {
  if (!validateForm()) return

  submitting.value = true

  try {
    const payload: CreateProviderRequest = {
      name: form.value.name.trim(),
      provider_type: form.value.provider_type as CreateProviderRequest['provider_type'],
      api_key: form.value.api_key,
      is_active: form.value.is_active,
    }

    if (form.value.base_url.trim()) {
      payload.base_url = form.value.base_url.trim()
    }

    if (form.value.default_model.trim()) {
      payload.default_model = form.value.default_model.trim()
    }

    if (form.value.input_cost_per_million) {
      payload.input_cost_per_million = parseFloat(form.value.input_cost_per_million)
    }

    if (form.value.output_cost_per_million) {
      payload.output_cost_per_million = parseFloat(form.value.output_cost_per_million)
    }

    if (editingProvider.value) {
      // For updates, only include api_key if provided
      const updatePayload = { ...payload }
      if (!form.value.api_key.trim()) {
        delete (updatePayload as Partial<CreateProviderRequest>).api_key
      }
      await providersStore.updateProvider(editingProvider.value.id, updatePayload)
      toast.success('Provider updated successfully')
    } else {
      await providersStore.createProvider(payload)
      toast.success('Provider created successfully')
    }

    closeModal()
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to save provider')
  } finally {
    submitting.value = false
  }
}

// Delete provider
const handleDelete = async () => {
  if (!deletingProvider.value) return

  deleting.value = true

  try {
    await providersStore.deleteProvider(deletingProvider.value.id)
    toast.success('Provider deleted successfully')
    closeDeleteModal()
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to delete provider')
  } finally {
    deleting.value = false
  }
}

// Test connection
const testProviderConnection = async (id: number) => {
  testingId.value = id

  try {
    const message = await providersStore.testConnection(id)
    toast.success(message)
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Connection test failed')
  } finally {
    testingId.value = null
  }
}

const testConnectionWithForm = async () => {
  // For OAuth providers, we can't test without connecting first
  if (form.value.provider_type === ProviderType.ANTHROPIC_MAX) {
    toast.info('OAuth providers require connecting via OAuth first')
    return
  }

  if (!form.value.api_key.trim() && !editingProvider.value) {
    toast.error('API key is required to test connection')
    return
  }

  testingConnection.value = true

  try {
    const payload: CreateProviderRequest = {
      name: form.value.name.trim() || 'Test',
      provider_type: form.value.provider_type as CreateProviderRequest['provider_type'],
      api_key: form.value.api_key,
    }

    if (form.value.base_url.trim()) {
      payload.base_url = form.value.base_url.trim()
    }

    const message = await providersStore.testConnectionWithCredentials(payload)
    toast.success(message)
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Connection test failed')
  } finally {
    testingConnection.value = false
  }
}

// OAuth functions
const connectingOAuthId = ref<number | null>(null)
const disconnectingOAuthId = ref<number | null>(null)

const connectOAuth = async (providerId: number) => {
  connectingOAuthId.value = providerId

  try {
    const { authorization_url } = await providersApi.getOAuthAuthorizeUrl(providerId)
    // Open OAuth in a new popup window
    const width = 600
    const height = 700
    const left = window.screenX + (window.outerWidth - width) / 2
    const top = window.screenY + (window.outerHeight - height) / 2
    const popup = window.open(
      authorization_url,
      'oauth_popup',
      `width=${width},height=${height},left=${left},top=${top},scrollbars=yes`
    )

    // Poll for popup closure and refresh providers
    if (popup) {
      const pollTimer = setInterval(async () => {
        if (popup.closed) {
          clearInterval(pollTimer)
          connectingOAuthId.value = null
          // Refresh providers to check if OAuth was connected
          await providersStore.fetchProviders()
        }
      }, 500)
    }
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to start OAuth flow')
    connectingOAuthId.value = null
  }
}

const disconnectOAuth = async (providerId: number) => {
  disconnectingOAuthId.value = providerId

  try {
    await providersApi.disconnectOAuth(providerId)
    toast.success('OAuth disconnected successfully')
    await providersStore.fetchProviders() // Refresh providers list
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to disconnect OAuth')
  } finally {
    disconnectingOAuthId.value = null
  }
}

const testOAuthConnection = async (providerId: number) => {
  testingId.value = providerId

  try {
    const response = await providersApi.testOAuthConnection(providerId)
    toast.success(response.message)
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'OAuth connection test failed')
  } finally {
    testingId.value = null
  }
}

// Check for OAuth callback results in URL
const checkOAuthCallback = () => {
  const urlParams = new URLSearchParams(window.location.search)
  const oauthSuccess = urlParams.get('oauth_success')
  const oauthError = urlParams.get('oauth_error')
  const errorDescription = urlParams.get('error_description')

  if (oauthSuccess === 'true') {
    toast.success('Successfully connected to Claude Max!')
    // Clean up URL params
    window.history.replaceState({}, document.title, window.location.pathname)
  } else if (oauthError) {
    const message = errorDescription || `OAuth error: ${oauthError}`
    toast.error(message)
    // Clean up URL params
    window.history.replaceState({}, document.title, window.location.pathname)
  }
}

// Load providers on mount
onMounted(async () => {
  // Check for OAuth callback results first
  checkOAuthCallback()

  try {
    await providersStore.fetchProviders()
  } catch (err) {
    toast.error('Failed to load providers')
  }
})
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
