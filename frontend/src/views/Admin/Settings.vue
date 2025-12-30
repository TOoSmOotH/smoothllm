<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />
    <div class="flex-1 py-10 px-6 sm:px-8">
      <div class="max-w-4xl mx-auto space-y-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-3xl font-display text-text-primary">Admin Settings</h1>
          <p class="text-text-tertiary">Configure global UI preferences for this deployment.</p>
        </div>
        <Button variant="outline" @click="router.push('/admin')">
          Back to Admin Dashboard
        </Button>
      </div>

      <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6 sm:p-8">
        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-6">
          <div class="flex-1">
            <h2 class="text-xl font-display text-text-primary mb-2">Theme Selector</h2>
            <p class="text-text-tertiary">
              Choose the visual theme that best fits the project you are templating.
            </p>
            <div class="mt-4">
              <label class="block text-sm font-medium text-text-secondary mb-2" for="theme">
                Active theme
              </label>
              <select
                id="theme"
                v-model="selectedTheme"
                class="w-full bg-bg-tertiary border border-border-default rounded-md px-4 py-3 text-text-primary focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10"
              >
                <option v-for="option in themeOptions" :key="option.id" :value="option.id">
                  {{ option.label }}
                </option>
              </select>
              <p class="mt-2 text-sm text-text-muted">
                {{ activeTheme?.description }}
              </p>
            </div>
          </div>
          <div class="bg-bg-tertiary border border-border-subtle rounded-lg p-4 sm:p-5 w-full sm:w-60">
            <p class="text-xs uppercase tracking-wide text-text-tertiary mb-3">Preview</p>
            <div class="space-y-3">
              <div class="h-2 rounded-full bg-primary-500"></div>
              <div class="h-2 rounded-full bg-secondary-500"></div>
              <div class="h-2 rounded-full bg-info-500"></div>
            </div>
            <div class="mt-4 text-xs text-text-tertiary">
              This updates live and persists for all admins.
            </div>
          </div>
        </div>
      </div>

      <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6 sm:p-8">
        <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-6">
          <div class="flex-1">
            <h2 class="text-xl font-display text-text-primary mb-2">User Access</h2>
            <p class="text-text-tertiary">
              Control whether registrations are open and if new users require approval.
            </p>
            <div class="mt-5 space-y-4">
              <label class="flex items-center justify-between gap-4 bg-bg-tertiary border border-border-default rounded-lg px-4 py-3">
                <span class="text-sm text-text-secondary">Accept new registrations</span>
                <input
                  type="checkbox"
                  v-model="registrationEnabled"
                  class="w-5 h-5 accent-primary-500"
                />
              </label>
              <label class="flex items-center justify-between gap-4 bg-bg-tertiary border border-border-default rounded-lg px-4 py-3">
                <span class="text-sm text-text-secondary">Auto-approve new users</span>
                <input
                  type="checkbox"
                  v-model="autoApproveNewUsers"
                  class="w-5 h-5 accent-primary-500"
                />
              </label>
            </div>
          </div>
          <div class="bg-bg-tertiary border border-border-subtle rounded-lg p-4 sm:p-5 w-full lg:w-64 space-y-4">
            <p class="text-xs uppercase tracking-wide text-text-tertiary">Registration</p>
            <div class="text-sm text-text-muted space-y-2">
              <p v-if="!registrationEnabled">Registrations are closed.</p>
              <p v-else-if="autoApproveNewUsers">New users are auto-approved.</p>
              <p v-else>New users require manual approval.</p>
            </div>
            <Button variant="primary" class="w-full" @click="saveRegistrationSettings">
              Save access settings
            </Button>
          </div>
        </div>
      </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import AppHeader from '@/components/layout/AppHeader.vue'
import Button from '@/components/ui/Button.vue'
import { useThemeStore, themeOptions, type ThemeId } from '@/stores/theme'
import { settingsApi } from '@/api/settings'

const router = useRouter()
const themeStore = useThemeStore()
const registrationEnabled = ref(true)
const autoApproveNewUsers = ref(true)

const selectedTheme = computed({
  get: () => themeStore.theme,
  set: async (value) => {
    try {
      await themeStore.setThemeForAll(value as ThemeId)
      toast.success('Theme updated for everyone')
    } catch (err) {
      toast.error('Failed to update theme')
    }
  },
})

const activeTheme = computed(() => {
  return themeOptions.find((option) => option.id === themeStore.theme)
})

const loadRegistrationSettings = async () => {
  try {
    const settings = await settingsApi.getRegistrationSettings()
    registrationEnabled.value = settings.registration_enabled
    autoApproveNewUsers.value = settings.auto_approve_new_users
  } catch (err) {
    toast.error('Failed to load registration settings')
  }
}

const saveRegistrationSettings = async () => {
  try {
    await settingsApi.updateRegistrationSettings({
      registration_enabled: registrationEnabled.value,
      auto_approve_new_users: autoApproveNewUsers.value,
    })
    toast.success('Registration settings updated')
  } catch (err) {
    toast.error('Failed to update registration settings')
  }
}

onMounted(() => {
  loadRegistrationSettings()
})
</script>
