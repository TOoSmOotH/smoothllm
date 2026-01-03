<template>
  <AppLayout>
    <div class="space-y-6 sm:space-y-8">
      <div class="bg-bg-secondary border border-border-subtle rounded-xl overflow-hidden">
        <div class="relative h-48 sm:h-56 md:h-64" :style="coverStyle">
          <div class="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent"></div>
        </div>

        <div class="relative px-6 pb-6">
          <div class="-mt-12 sm:-mt-16 md:-mt-20 flex flex-col sm:flex-row sm:items-end gap-4">
            <div class="relative">
              <div
                class="w-24 h-24 sm:w-28 sm:h-28 md:w-32 md:h-32 rounded-full border-4 border-bg-secondary bg-bg-tertiary overflow-hidden flex items-center justify-center text-2xl sm:text-3xl font-display text-text-primary"
              >
                <img
                  v-if="avatarUrl"
                  :src="avatarUrl"
                  alt="Profile avatar"
                  class="w-full h-full object-cover"
                />
                <span v-else>{{ initials }}</span>
              </div>
            </div>

            <div class="flex-1">
              <h1 class="text-3xl sm:text-4xl font-display text-text-primary">
                {{ displayName }}
              </h1>
              <p class="text-text-tertiary mt-1">@{{ profile?.username }}</p>
              <p v-if="headline" class="text-text-secondary mt-2">
                {{ headline }}
              </p>
              <div class="flex flex-wrap gap-2 mt-4">
                <span v-if="profile?.pronouns" class="badge">{{ profile.pronouns }}</span>
                <span v-if="profile?.location" class="badge">{{ profile.location }}</span>
                <span v-if="profile?.language" class="badge">{{ profile.language }}</span>
              </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-3 sm:items-center">
              <Button variant="primary" @click="router.push('/profile/edit')">
                Edit Profile
              </Button>
              <Button variant="outline" @click="router.push('/dashboard')">
                Back to Dashboard
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="loading" class="card text-center py-12">
        <p class="text-text-tertiary">Loading profile...</p>
      </div>

      <div v-else-if="error" class="card text-center py-12">
        <p class="text-error-500 mb-4">{{ error }}</p>
        <Button variant="outline" @click="fetchProfile">Retry</Button>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <div class="card">
            <h2 class="font-display text-xl text-text-primary mb-4">About</h2>
            <p v-if="profile?.bio" class="text-text-secondary leading-relaxed">
              {{ profile.bio }}
            </p>
            <p v-else class="text-text-tertiary">Add a bio to share your story.</p>
          </div>

          <div class="card">
            <h2 class="font-display text-xl text-text-primary mb-4">Contact</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Email</p>
                <p class="text-text-primary">{{ emailLabel }}</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Phone</p>
                <p class="text-text-primary">{{ profile?.phone || 'Not set' }}</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Website</p>
                <a
                  v-if="profile?.website"
                  :href="profile.website"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-primary-500 hover:text-primary-400 transition-colors"
                >
                  {{ profile.website }}
                </a>
                <p v-else class="text-text-primary">Not set</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Location</p>
                <p class="text-text-primary">{{ locationLabel }}</p>
              </div>
            </div>
          </div>

          <div class="card">
            <h2 class="font-display text-xl text-text-primary mb-4">Professional</h2>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Role</p>
                <p class="text-text-primary">{{ profile?.job_title || 'Not set' }}</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Company</p>
                <p class="text-text-primary">{{ profile?.company || 'Not set' }}</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">LinkedIn</p>
                <a
                  v-if="profile?.linkedin_url"
                  :href="profile.linkedin_url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-primary-500 hover:text-primary-400 transition-colors"
                >
                  {{ profile.linkedin_url }}
                </a>
                <p v-else class="text-text-primary">Not set</p>
              </div>
              <div>
                <p class="text-xs uppercase tracking-wide text-text-tertiary">Portfolio</p>
                <a
                  v-if="profile?.portfolio_url"
                  :href="profile.portfolio_url"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-primary-500 hover:text-primary-400 transition-colors"
                >
                  {{ profile.portfolio_url }}
                </a>
                <p v-else class="text-text-primary">Not set</p>
              </div>
            </div>
          </div>

          <SocialLinksManager />
        </div>

        <div class="space-y-6">
          <ProfileCompletionProgress />
          <MilestonesDisplay />

          <div class="card">
            <h2 class="font-display text-xl text-text-primary mb-4">Profile Details</h2>
            <div class="space-y-3">
              <div class="flex items-center justify-between text-sm">
                <span class="text-text-tertiary">Joined</span>
                <span class="text-text-primary">{{ joinedDate }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-text-tertiary">Last Updated</span>
                <span class="text-text-primary">{{ updatedDate }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import AppLayout from '@/components/layout/AppLayout.vue'
import Button from '@/components/ui/Button.vue'
import SocialLinksManager from '@/components/ui/SocialLinksManager.vue'
import ProfileCompletionProgress from '@/components/ui/ProfileCompletionProgress.vue'
import MilestonesDisplay from '@/components/ui/MilestonesDisplay.vue'
import { profileApi, type ProfileResponse, type MediaFile } from '@/api/profile'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const profile = ref<ProfileResponse | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const getApiBase = (): string => {
  const envUrl = (import.meta as any).env?.VITE_API_URL
  if (envUrl) return envUrl
  const { protocol, hostname } = window.location
  return `${protocol}//${hostname}:8080/api/v1`
}
const apiBase = getApiBase()

const getMediaUrl = (media?: MediaFile | null) => {
  if (!media?.id) return ''
  return `${apiBase}/media/${media.id}`
}

const coverStyle = computed(() => {
  const coverUrl = getMediaUrl(profile.value?.cover_photo)
  if (coverUrl) {
    return {
      backgroundImage: `url(${coverUrl})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
    }
  }
  return {
    background:
      'linear-gradient(135deg, rgba(14,165,233,0.2) 0%, rgba(148,163,184,0.1) 50%, rgba(14,116,144,0.2) 100%)',
  }
})

const avatarUrl = computed(() => getMediaUrl(profile.value?.avatar))

const displayName = computed(() => {
  const data = profile.value
  if (!data) return 'Profile'
  if (data.display_name) return data.display_name
  const fullName = [data.first_name, data.last_name].filter(Boolean).join(' ')
  return fullName || data.username || 'Profile'
})

const initials = computed(() => {
  const name = displayName.value
  return name
    .split(' ')
    .map((part) => part.charAt(0).toUpperCase())
    .slice(0, 2)
    .join('')
})

const headline = computed(() => {
  const data = profile.value
  if (!data) return ''
  if (data.job_title && data.company) return `${data.job_title} · ${data.company}`
  return data.job_title || data.company || ''
})

const locationLabel = computed(() => {
  const data = profile.value
  if (!data) return 'Not set'
  if (data.location) return data.location
  const parts = [data.city, data.state, data.country].filter(Boolean)
  return parts.length > 0 ? parts.join(', ') : 'Not set'
})

const emailLabel = computed(() => {
  return authStore.user?.email || 'Not set'
})

const joinedDate = computed(() => {
  const data = profile.value
  if (!data?.created_at) return '—'
  return new Date(data.created_at).toLocaleDateString()
})

const updatedDate = computed(() => {
  const data = profile.value
  if (!data?.updated_at) return '—'
  return new Date(data.updated_at).toLocaleDateString()
})

const fetchProfile = async () => {
  loading.value = true
  error.value = null
  try {
    profile.value = await profileApi.getMyProfile()
  } catch (err) {
    error.value = 'Failed to load profile'
    toast.error('Failed to load profile')
    console.error('Failed to load profile:', err)
  } finally {
    loading.value = false
  }
}

onMounted(fetchProfile)
</script>
