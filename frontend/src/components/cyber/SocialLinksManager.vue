<template>
  <div class="bg-cyber-gray/30 backdrop-blur-sm border border-cyber-cyan/30 rounded-cyber p-6">
    <div class="flex items-center justify-between mb-6">
      <h3 class="font-cyber text-xl text-cyber-cyan">Social Links</h3>
      <CyberButton variant="primary" size="sm" @click="showAddModal = true" :loading="loading">
        + Add Link
      </CyberButton>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="animate-spin text-cyber-cyan">
        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    </div>

    <div v-else-if="sortedLinks.length === 0" class="text-center py-8 text-cyber-light-gray/60">
      <p class="font-mono-cyber">No social links added yet</p>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="link in sortedLinks"
        :key="link.id"
        class="group flex items-center gap-4 bg-cyber-dark/50 border border-cyber-cyan/20 rounded-cyber p-4 hover:border-cyber-cyan/50 transition-all duration-300"
      >
        <div class="flex-shrink-0">
          <span class="text-2xl">{{ getPlatformIcon(link.platform) }}</span>
        </div>

        <div class="flex-grow min-w-0">
          <p class="font-cyber text-white truncate">{{ link.platform }}</p>
          <a
            v-if="link.url"
            :href="link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-cyber-light-gray/60 hover:text-cyber-cyan transition-colors truncate block"
          >
            {{ link.url }}
          </a>
        </div>

        <div class="flex items-center gap-2 flex-shrink-0">
          <button
            @click="handleToggleVisibility(link)"
            :class="[
              'p-2 rounded transition-colors duration-300',
              link.visible ? 'text-cyber-green hover:bg-cyber-green/10' : 'text-cyber-light-gray/40 hover:bg-cyber-light-gray/10'
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
            class="p-2 text-cyber-cyan hover:bg-cyber-cyan/10 rounded transition-colors duration-300"
            title="Edit"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>

          <button
            @click="handleDelete(link)"
            class="p-2 text-cyber-pink hover:bg-cyber-pink/10 rounded transition-colors duration-300"
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
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div class="bg-cyber-gray border border-cyber-cyan/50 rounded-cyber p-6 max-w-md w-full mx-4 shadow-cyber-cyan">
        <h3 class="font-cyber text-xl text-cyber-cyan mb-6">
          {{ showEditModal ? 'Edit Social Link' : 'Add Social Link' }}
        </h3>

        <CyberInput
          v-model="formData.platform"
          label="Platform"
          placeholder="e.g., Twitter, GitHub, LinkedIn"
          :error="errors.platform"
          class="mb-4"
        />

        <CyberInput
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
            class="w-4 h-4 accent-cyber-cyan"
          />
          <label for="visible" class="text-sm text-cyber-light-gray font-mono-cyber">
            Make visible on profile
          </label>
        </div>

        <div class="flex justify-end gap-3">
          <CyberButton variant="ghost" @click="closeModal">Cancel</CyberButton>
          <CyberButton
            variant="primary"
            @click="handleSubmit"
            :loading="submitting"
          >
            {{ showEditModal ? 'Update' : 'Add' }}
          </CyberButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { toast } from 'vue-sonner'
import CyberButton from '@/components/cyber/CyberButton.vue'
import CyberInput from '@/components/cyber/CyberInput.vue'
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
