<template>
  <AppLayout>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-display text-text-primary mb-2">API Keys</h1>
        <p class="text-text-muted">Manage your proxy API keys for LLM access</p>
      </div>
      <Button variant="primary" @click="openCreateModal" :disabled="!providersStore.hasProviders">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Create Key
      </Button>
    </div>

    <!-- No Providers Warning -->
    <div
      v-if="!providersStore.hasProviders && providersStore.isInitialized"
      class="bg-warning-500/10 border border-warning-500/20 rounded-lg p-4 mb-6"
    >
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-warning-500 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <p class="text-text-primary font-medium">No providers configured</p>
          <p class="text-text-muted text-sm">You need to configure at least one LLM provider before creating API keys.</p>
        </div>
        <Button variant="outline" size="sm" @click="router.push('/providers')" class="ml-auto">
          Add Provider
        </Button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="apiKeysStore.loading && !apiKeysStore.isInitialized" class="flex justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="!apiKeysStore.hasApiKeys && providersStore.hasProviders"
      class="bg-bg-secondary border border-border-subtle rounded-lg p-12 text-center"
    >
      <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
        </svg>
      </div>
      <h3 class="text-lg font-display text-text-primary mb-2">No API Keys</h3>
      <p class="text-text-muted mb-6">Create your first proxy API key to start using the LLM proxy.</p>
      <Button variant="primary" @click="openCreateModal">Create Your First Key</Button>
    </div>

    <!-- API Keys List -->
    <div v-else-if="apiKeysStore.hasApiKeys" class="space-y-4">
      <div
        v-for="apiKey in apiKeysStore.apiKeys"
        :key="apiKey.id"
        class="bg-bg-secondary border border-border-subtle rounded-lg p-6 hover:border-border-default transition-colors duration-200"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <h3 class="font-display font-semibold text-lg text-text-primary">{{ apiKey.name || 'Unnamed Key' }}</h3>
              <span
                :class="[
                  'px-2 py-0.5 rounded-full text-xs font-medium uppercase',
                  apiKey.is_active
                    ? 'bg-success-500/10 text-success-500'
                    : 'bg-text-muted/10 text-text-muted'
                ]"
              >
                {{ apiKey.is_active ? 'Active' : 'Revoked' }}
              </span>
              <span
                v-if="isExpired(apiKey)"
                class="px-2 py-0.5 rounded-full text-xs font-medium uppercase bg-error-500/10 text-error-500"
              >
                Expired
              </span>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-sm">
              <div>
                <p class="text-text-tertiary mb-1">Key Prefix</p>
                <p class="text-text-primary font-mono">{{ apiKey.key_prefix }}</p>
              </div>
              <div>
                <p class="text-text-tertiary mb-1">Provider</p>
                <p class="text-text-primary font-mono">{{ getProviderName(apiKey.provider_id) }}</p>
              </div>
              <div>
                <p class="text-text-tertiary mb-1">Last Used</p>
                <p class="text-text-primary font-mono">{{ formatDate(apiKey.last_used_at) }}</p>
              </div>
              <div>
                <p class="text-text-tertiary mb-1">Expires</p>
                <p class="text-text-primary font-mono">{{ formatDate(apiKey.expires_at) || 'Never' }}</p>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2 ml-4">
            <Button variant="ghost" size="sm" @click="openEditModal(apiKey)">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </Button>
            <Button
              v-if="apiKey.is_active"
              variant="ghost"
              size="sm"
              @click="confirmRevoke(apiKey)"
            >
              <svg class="w-4 h-4 text-warning-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
              </svg>
            </Button>
            <Button variant="ghost" size="sm" @click="confirmDelete(apiKey)">
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
        <div class="relative z-10 bg-bg-primary border border-border-subtle rounded-lg shadow-xl max-w-lg w-full max-h-[90vh] overflow-y-auto">
          <div class="p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">
              {{ editingKey ? 'Edit API Key' : 'Create API Key' }}
            </h2>

            <div class="space-y-4">
              <Input
                v-model="form.name"
                label="Key Name"
                placeholder="My API Key"
                helper-text="A friendly name to identify this key"
                :error="errors.name"
              />

              <div v-if="!editingKey">
                <label class="block text-sm font-medium text-text-secondary mb-2">Provider</label>
                <select
                  v-model="form.provider_id"
                  class="w-full font-sans bg-bg-secondary border border-border-default rounded-md text-text-primary px-4 py-3 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all duration-200"
                >
                  <option :value="0" disabled>Select a provider...</option>
                  <option
                    v-for="provider in providersStore.activeProviders"
                    :key="provider.id"
                    :value="provider.id"
                  >
                    {{ provider.name }} ({{ provider.provider_type }})
                  </option>
                </select>
                <p v-if="errors.provider_id" class="mt-1 text-xs text-error-500 font-medium">{{ errors.provider_id }}</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-text-secondary mb-2">Expiration Date (optional)</label>
                <input
                  v-model="form.expires_at"
                  type="date"
                  class="w-full font-sans bg-bg-secondary border border-border-default rounded-md text-text-primary px-4 py-3 focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all duration-200"
                />
                <p class="mt-1 text-xs text-text-tertiary">Leave empty for a key that never expires</p>
              </div>
            </div>

            <div class="flex justify-end gap-3 mt-8 pt-6 border-t border-border-subtle">
              <Button variant="ghost" @click="closeModal">Cancel</Button>
              <Button
                variant="primary"
                @click="handleSubmit"
                :loading="submitting"
              >
                {{ editingKey ? 'Save Changes' : 'Create Key' }}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- New Key Created Modal -->
    <Teleport to="body">
      <div
        v-if="showNewKeyModal && newlyCreatedKey"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <div class="fixed inset-0 bg-black/50 z-0"></div>
        <div class="relative z-10 bg-bg-primary border border-border-subtle rounded-lg shadow-xl max-w-lg w-full">
          <div class="p-6">
            <div class="flex items-center gap-4 mb-6">
              <div class="w-12 h-12 bg-success-500/10 rounded-full flex items-center justify-center">
                <svg class="w-6 h-6 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <h3 class="font-display text-lg text-text-primary">API Key Created</h3>
                <p class="text-sm text-text-muted">Copy your key now - it won't be shown again!</p>
              </div>
            </div>

            <div class="bg-bg-tertiary rounded-lg p-4 mb-6">
              <div class="flex items-center justify-between gap-4">
                <code class="text-sm text-text-primary font-mono break-all">{{ newlyCreatedKey.key }}</code>
                <Button variant="ghost" size="sm" @click="copyKey" class="flex-shrink-0">
                  <svg v-if="!keyCopied" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                  <svg v-else class="w-4 h-4 text-success-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </Button>
              </div>
            </div>

            <div class="bg-warning-500/10 border border-warning-500/20 rounded-lg p-4 mb-6">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-warning-500 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <div>
                  <p class="text-text-primary font-medium text-sm">Save this key securely</p>
                  <p class="text-text-muted text-sm">This is the only time you'll see this key. Store it in a secure location.</p>
                </div>
              </div>
            </div>

            <div class="flex justify-end">
              <Button variant="primary" @click="closeNewKeyModal">Done</Button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Revoke Confirmation Modal -->
    <Teleport to="body">
      <div
        v-if="showRevokeModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="closeRevokeModal"
      >
        <div class="fixed inset-0 bg-black/50 z-0" @click="closeRevokeModal"></div>
        <div class="relative z-10 bg-bg-primary border border-border-subtle rounded-lg shadow-xl max-w-md w-full">
          <div class="p-6">
            <div class="flex items-center gap-4 mb-4">
              <div class="w-12 h-12 bg-warning-500/10 rounded-full flex items-center justify-center">
                <svg class="w-6 h-6 text-warning-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                </svg>
              </div>
              <div>
                <h3 class="font-display text-lg text-text-primary">Revoke API Key</h3>
                <p class="text-sm text-text-muted">This will immediately stop this key from working.</p>
              </div>
            </div>
            <p class="text-text-secondary mb-6">
              Are you sure you want to revoke <strong class="text-text-primary">{{ revokingKey?.name || revokingKey?.key_prefix }}</strong>?
              The key will be deactivated but can be reactivated later.
            </p>
            <div class="flex justify-end gap-3">
              <Button variant="ghost" @click="closeRevokeModal">Cancel</Button>
              <Button variant="secondary" @click="handleRevoke" :loading="revoking">
                Revoke Key
              </Button>
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
                <h3 class="font-display text-lg text-text-primary">Delete API Key</h3>
                <p class="text-sm text-text-muted">This action cannot be undone.</p>
              </div>
            </div>
            <p class="text-text-secondary mb-6">
              Are you sure you want to permanently delete <strong class="text-text-primary">{{ deletingKey?.name || deletingKey?.key_prefix }}</strong>?
              All usage history for this key will also be deleted.
            </p>
            <div class="flex justify-end gap-3">
              <Button variant="ghost" @click="closeDeleteModal">Cancel</Button>
              <Button variant="destructive" @click="handleDelete" :loading="deleting">
                Delete Key
              </Button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useApiKeysStore } from '@/stores/apiKeys'
import { useProvidersStore } from '@/stores/providers'
import type { KeyResponse, KeyCreateResponse, CreateKeyRequest } from '@/api/keys'
import AppLayout from '@/components/layout/AppLayout.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'

const router = useRouter()
const apiKeysStore = useApiKeysStore()
const providersStore = useProvidersStore()

// Modal state
const showModal = ref(false)
const showNewKeyModal = ref(false)
const showRevokeModal = ref(false)
const showDeleteModal = ref(false)
const editingKey = ref<KeyResponse | null>(null)
const revokingKey = ref<KeyResponse | null>(null)
const deletingKey = ref<KeyResponse | null>(null)
const newlyCreatedKey = ref<KeyCreateResponse | null>(null)
const keyCopied = ref(false)

// Form state
const form = ref({
  name: '',
  provider_id: 0,
  expires_at: '',
})

const errors = ref<Record<string, string>>({})

// Loading states
const submitting = ref(false)
const revoking = ref(false)
const deleting = ref(false)

// Helper functions
const getProviderName = (providerId: number): string => {
  const provider = providersStore.getProviderById(providerId)
  return provider?.name || 'Unknown Provider'
}

const formatDate = (dateString?: string | null): string => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const isExpired = (key: KeyResponse): boolean => {
  if (!key.expires_at) return false
  return new Date(key.expires_at) < new Date()
}

// Modal functions
const openCreateModal = () => {
  editingKey.value = null
  resetForm()
  showModal.value = true
}

const openEditModal = (key: KeyResponse) => {
  editingKey.value = key
  form.value = {
    name: key.name || '',
    provider_id: key.provider_id,
    expires_at: key.expires_at ? key.expires_at.split('T')[0] : '',
  }
  errors.value = {}
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingKey.value = null
  resetForm()
}

const closeNewKeyModal = () => {
  showNewKeyModal.value = false
  newlyCreatedKey.value = null
  keyCopied.value = false
}

const confirmRevoke = (key: KeyResponse) => {
  revokingKey.value = key
  showRevokeModal.value = true
}

const closeRevokeModal = () => {
  showRevokeModal.value = false
  revokingKey.value = null
}

const confirmDelete = (key: KeyResponse) => {
  deletingKey.value = key
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deletingKey.value = null
}

const resetForm = () => {
  form.value = {
    name: '',
    provider_id: providersStore.activeProviders[0]?.id || 0,
    expires_at: '',
  }
  errors.value = {}
}

// Validation
const validateForm = (): boolean => {
  errors.value = {}

  if (!editingKey.value && form.value.provider_id === 0) {
    errors.value.provider_id = 'Please select a provider'
  }

  return Object.keys(errors.value).length === 0
}

// Form submission
const handleSubmit = async () => {
  if (!validateForm()) return

  submitting.value = true

  try {
    if (editingKey.value) {
      await apiKeysStore.updateApiKey(editingKey.value.id, {
        name: form.value.name.trim() || undefined,
        expires_at: form.value.expires_at || undefined,
      })
      toast.success('API key updated successfully')
      closeModal()
    } else {
      const payload: CreateKeyRequest = {
        provider_id: form.value.provider_id,
      }

      if (form.value.name.trim()) {
        payload.name = form.value.name.trim()
      }

      if (form.value.expires_at) {
        payload.expires_at = form.value.expires_at
      }

      const result = await apiKeysStore.createApiKey(payload)
      newlyCreatedKey.value = result
      closeModal()
      showNewKeyModal.value = true
    }
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to save API key')
  } finally {
    submitting.value = false
  }
}

// Revoke key
const handleRevoke = async () => {
  if (!revokingKey.value) return

  revoking.value = true

  try {
    await apiKeysStore.revokeApiKey(revokingKey.value.id)
    toast.success('API key revoked successfully')
    closeRevokeModal()
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to revoke API key')
  } finally {
    revoking.value = false
  }
}

// Delete key
const handleDelete = async () => {
  if (!deletingKey.value) return

  deleting.value = true

  try {
    await apiKeysStore.deleteApiKey(deletingKey.value.id)
    toast.success('API key deleted successfully')
    closeDeleteModal()
  } catch (err: unknown) {
    const error = err as { response?: { data?: { error?: string } } }
    toast.error(error.response?.data?.error || 'Failed to delete API key')
  } finally {
    deleting.value = false
  }
}

// Copy key to clipboard
const copyKey = async () => {
  if (!newlyCreatedKey.value?.key) return

  try {
    await navigator.clipboard.writeText(newlyCreatedKey.value.key)
    keyCopied.value = true
    toast.success('API key copied to clipboard')
    setTimeout(() => {
      keyCopied.value = false
    }, 2000)
  } catch {
    toast.error('Failed to copy to clipboard')
  }
}

// Load data on mount
onMounted(async () => {
  try {
    await Promise.all([
      providersStore.fetchProviders(),
      apiKeysStore.fetchApiKeys(),
    ])
  } catch {
    toast.error('Failed to load data')
  }
})
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
