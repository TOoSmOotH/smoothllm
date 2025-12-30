<template>
  <div class="bg-bg-secondary border border-border-subtle rounded-lg p-4 sm:p-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4 sm:mb-6">
      <h3 class="font-display text-lg sm:text-xl text-text-primary">Social Links</h3>
      <Button variant="primary" size="sm" @click="showAddModal = true" :loading="loading" class="min-h-[36px] min-w-[36px]">
        + Add Link
      </Button>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="animate-spin text-primary-500">
        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    </div>

    <div v-else-if="sortedLinks.length === 0" class="text-center py-8 text-text-muted">
      <p class="font-mono">No social links added yet</p>
    </div>

    <div v-else class="space-y-2 sm:space-y-3">
      <div
        v-for="link in sortedLinks"
        :key="link.id"
        class="group flex items-center gap-3 sm:gap-4 bg-bg-tertiary border border-border-default rounded-md p-3 sm:p-4 hover:border-primary-500/50 transition-all duration-200"
      >
        <div class="flex-shrink-0">
          <span class="text-xl sm:text-2xl">{{ getPlatformIcon(link.platform) }}</span>
        </div>

        <div class="flex-grow min-w-0">
          <p class="font-sans text-sm sm:text-base text-text-primary truncate">{{ link.platform }}</p>
          <a
            v-if="link.url"
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-xs sm:text-sm text-text-muted hover:text-primary-500 transition-colors truncate block"
          >
            {{ link.url }}
          </a>
        </div>

        <div class="flex items-center gap-2 flex-shrink-0">
          <button
            @click="handleToggleVisibility(link)"
            :class="[
              'p-2 rounded transition-colors duration-200',
              link.visible ? 'text-success-500 hover:bg-success-500/10' : 'text-text-muted hover:bg-bg-elevated'
            ]"
            :title="link.visible ? 'Visible' : 'Hidden'"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
          </button>

          <button
            @click="handleEdit(link)"
            class="p-2 text-primary-500 hover:bg-primary-500/10 rounded transition-colors duration-200"
            title="Edit"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>

          <button
            @click="handleDelete(link)"
            class="p-2 text-error-500 hover:bg-error-500/10 rounded transition-colors duration-200"
            title="Delete"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
      <div class="bg-bg-secondary border border-border-default rounded-lg p-4 sm:p-6 max-w-md w-full mx-auto shadow-xl">
        <h3 class="font-display text-xl text-text-primary mb-6">
          {{ showEditModal ? 'Edit Social Link' : 'Add Social Link' }}
        </h3>

        <Input
          v-model="formData.platform"
          label="Platform"
          placeholder="e.g., Twitter, GitHub, LinkedIn"
          :error="errors.platform"
          class="mb-4"
        />

        <Input
          v-model="formData.url"
          type="url"
          label="URL"
          placeholder="https://..."
          :error="errors.url"
          class="mb-4"
        />

        <div class="flex items-center gap-2 mb-6">
          <input
            type="checkbox"
            id="visible"
            v-model="formData.visible"
            class="w-4 h-4 accent-primary-500"
          />
          <label for="visible" class="text-sm text-text-muted font-mono">
            Make visible on profile
          </label>
        </div>

        <div class="flex justify-end gap-3">
          <Button variant="ghost" @click="closeModal">Cancel</Button>
          <Button
            variant="primary"
            @click="handleSubmit"
            :loading="submitting"
          >
            {{ showEditModal ? 'Update' : 'Add' }}
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { toast } from 'vue-sonner'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { useSocialLinks } from '@/composables/useSocialLinks'
import type { SocialLink, CreateSocialLinkRequest, UpdateSocialLinkRequest } from '@/api/social'
const {
  socialLinks,
  loading,
  addSocialLink,
  editSocialLink,
  removeSocialLink,
  toggleLinkVisibility,
} = useSocialLinks()

const showAddModal = ref(false)
const showEditModal = ref(false)
const editingLink = ref<SocialLink | null>(null)
const submitting = ref(false)
const errors = ref<Record<string, string>>({})
const formData = ref<CreateSocialLinkRequest>({
  platform: '',
  url: '',
  visible: true,
})

const sortedLinks = computed(() => {
  return [...socialLinks.value].sort((a, b) => a.order - b.order)
})

const getPlatformIcon = (platform: string): string => {
  const icons: Record<string, string> = {
    twitter: '🐦',
    x: '❌',
    github: '🐙',
    linkedin: '💼',
    instagram: '📷',
    facebook: '📘',
    youtube: '▶️',
    discord: '💬',
    reddit: '🤖',
    tiktok: '🎵',
    twitch: '📺',
    website: '🌐',
    blog: '📝',
    default: '🔗',
  }
  return icons[platform.toLowerCase()] || icons.default
}

const resetForm = () => {
  formData.value = {
    platform: '',
    url: '',
    visible: true,
  }
  errors.value = {}
}

const validateForm = (): boolean => {
  errors.value = {}

  if (!formData.value.platform?.trim()) {
    errors.value.platform = 'Platform is required'
  }

  if (!formData.value.url?.trim()) {
    errors.value.url = 'URL is required'
  } else if (!formData.value.url.startsWith('http://') && !formData.value.url.startsWith('https://')) {
    errors.value.url = 'URL must start with http:// or https://'
  }

  return Object.keys(errors.value).length === 0
}

const handleSubmit = async () => {
  if (!validateForm()) return

  submitting.value = true

  try {
    if (showEditModal.value && editingLink.value) {
      await editSocialLink(editingLink.value.id, {
        platform: formData.value.platform,
        url: formData.value.url,
        visible: formData.value.visible,
      })
    } else {
      await addSocialLink(formData.value)
    }
    closeModal()
  } catch (err) {
    console.error('Failed to save social link:', err)
  } finally {
    submitting.value = false
  }
}

const handleEdit = (link: SocialLink) => {
  editingLink.value = link
  formData.value = {
    platform: link.platform,
    url: link.url,
    visible: link.visible,
  }
  showEditModal.value = true
}

const handleDelete = async (link: SocialLink) => {
  if (!confirm('Are you sure you want to delete this social link?')) return

  try {
    await removeSocialLink(link.id)
  } catch (err) {
    console.error('Failed to delete social link:', err)
  }
}

const handleToggleVisibility = async (link: SocialLink) => {
  try {
    await toggleLinkVisibility(link.id)
  } catch (err) {
    console.error('Failed to toggle visibility:', err)
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  editingLink.value = null
  resetForm()
}
</script>
