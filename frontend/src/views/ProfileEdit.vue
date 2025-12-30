<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />
    <div class="flex-1 py-8 px-6 sm:px-8">
      <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-display text-text-primary mb-2">Edit Profile</h1>
          <p class="text-text-muted">Update your profile information</p>
        </div>
        <Button variant="outline" @click="router.push('/dashboard')">
          Back to Dashboard
        </Button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-6">
          <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">Basic Information</h2>
 
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <Input
                v-model="profile.first_name"
                label="First Name"
                placeholder="John"
                :error="errors.first_name"
                class="mb-4"
              />
 
              <Input
                v-model="profile.last_name"
                label="Last Name"
                placeholder="Doe"
                :error="errors.last_name"
                class="mb-4"
              />
            </div>
 
            <Input
              v-model="profile.display_name"
              label="Display Name"
              placeholder="johndoe"
              :error="errors.display_name"
              class="mb-4"
            />
 
            <div class="mb-4">
              <label class="block text-sm font-medium text-text-secondary mb-2">Bio</label>
              <textarea
                v-model="profile.bio"
                :class="[
                  'w-full font-sans bg-bg-secondary border border-border-default rounded-md text-text-primary p-4 min-h-[120px] focus:outline-none focus:border-primary-500 focus:ring-2 focus:ring-primary-500/10 transition-all duration-200',
                  errors.bio ? 'border-error-500' : ''
                ]"
                placeholder="Tell us about yourself..."
              ></textarea>
              <p v-if="errors.bio" class="mt-1 text-sm text-error-500">{{ errors.bio }}</p>
            </div>
          </div>
 
          <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">Contact Information</h2>
 
            <Input
              v-model="profile.phone"
              type="tel"
              label="Phone"
              placeholder="+1 (555) 123-4567"
              :error="errors.phone"
              class="mb-4"
            />
 
            <Input
              v-model="profile.website"
              type="url"
              label="Website"
              placeholder="https://yourwebsite.com"
              :error="errors.website"
              class="mb-4"
            />
 
            <Input
              v-model="profile.location"
              label="Location"
              placeholder="San Francisco, CA"
              :error="errors.location"
              class="mb-4"
            />
          </div>
 
          <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">Personal Information</h2>
 
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <Input
                v-model="profile.birthday"
                type="date"
                label="Birthday"
                :error="errors.birthday"
                class="mb-4"
              />
 
              <Input
                v-model="profile.gender"
                label="Gender"
                placeholder="e.g., male, female, non-binary"
                :error="errors.gender"
                class="mb-4"
              />
            </div>
 
            <Input
              v-model="profile.pronouns"
              label="Pronouns"
              placeholder="e.g., he/him, she/her, they/them"
              :error="errors.pronouns"
              class="mb-4"
            />
 
            <Input
              v-model="profile.language"
              label="Language"
              placeholder="English"
              :error="errors.language"
              class="mb-4"
            />
          </div>
 
          <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">Professional Information</h2>
 
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <Input
                v-model="profile.job_title"
                label="Job Title"
                placeholder="Software Engineer"
                :error="errors.job_title"
                class="mb-4"
              />
 
              <Input
                v-model="profile.company"
                label="Company"
                placeholder="Tech Corp"
                :error="errors.company"
                class="mb-4"
              />
            </div>
 
            <Input
              v-model="profile.linkedin_url"
              type="url"
              label="LinkedIn URL"
              placeholder="https://linkedin.com/in/yourprofile"
              :error="errors.linkedin_url"
              class="mb-4"
            />
 
            <Input
              v-model="profile.portfolio_url"
              type="url"
              label="Portfolio URL"
              placeholder="https://yourportfolio.com"
              :error="errors.portfolio_url"
              class="mb-4"
            />
          </div>
 
          <SocialLinksManager />
        </div>
 
        <div class="space-y-6">
          <ProfileCompletionProgress />
 
          <MilestonesDisplay />
 
          <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6">
            <h2 class="font-display text-xl text-text-primary mb-6">Actions</h2>
            <div class="space-y-3">
              <Button
                variant="primary"
                class="w-full"
                @click="handleSave"
                :loading="saving"
              >
                Save Changes
              </Button>
 
              <Button
                variant="secondary"
                class="w-full"
                @click="router.push('/profile')"
              >
                View Profile
              </Button>
 
              <Button
                variant="ghost"
                class="w-full"
                @click="router.push('/dashboard')"
              >
                Cancel
              </Button>
            </div>
          </div>
        </div>
      </div>
      </div>
    </div>
  </div>
</template>
 
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { profileApi, type ProfileResponse } from '@/api/profile'
import AppHeader from '@/components/layout/AppHeader.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import SocialLinksManager from '@/components/ui/SocialLinksManager.vue'
import ProfileCompletionProgress from '@/components/ui/ProfileCompletionProgress.vue'
import MilestonesDisplay from '@/components/ui/MilestonesDisplay.vue'
 
const router = useRouter()
 
const saving = ref(false)
const profile = ref({
  first_name: '',
  last_name: '',
  display_name: '',
  bio: '',
  phone: '',
  website: '',
  location: '',
  birthday: '',
  gender: '',
  pronouns: '',
  language: '',
  job_title: '',
  company: '',
  linkedin_url: '',
  portfolio_url: '',
})
 
const errors = ref<Record<string, string>>({})

const buildUpdatePayload = () => {
  const payload = { ...profile.value }
  if (payload.birthday) {
    payload.birthday = new Date(`${payload.birthday}T00:00:00Z`).toISOString()
  } else {
    delete (payload as any).birthday
  }
  return payload
}

const normalizeProfile = (profileData: ProfileResponse) => {
  const formatDate = (value?: string | null) => {
    if (!value) return ''
    const date = new Date(value)
    if (Number.isNaN(date.getTime())) return ''
    return date.toISOString().slice(0, 10)
  }

  return {
    first_name: profileData.first_name ?? '',
    last_name: profileData.last_name ?? '',
    display_name: profileData.display_name ?? '',
    bio: profileData.bio ?? '',
    phone: profileData.phone ?? '',
    website: profileData.website ?? '',
    location: profileData.location ?? '',
    birthday: formatDate(profileData.birthday),
    gender: profileData.gender ?? '',
    pronouns: profileData.pronouns ?? '',
    language: profileData.language ?? '',
    job_title: profileData.job_title ?? '',
    company: profileData.company ?? '',
    linkedin_url: profileData.linkedin_url ?? '',
    portfolio_url: profileData.portfolio_url ?? '',
  }
}

const handleSave = async () => {
  errors.value = {}
  saving.value = true

  try {
    await profileApi.updateProfile(buildUpdatePayload())

    toast.success('Profile updated successfully')
  } catch (err: any) {
    toast.error('Failed to update profile')
    console.error('Failed to update profile:', err)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const profileData = await profileApi.getMyProfile()
    profile.value = normalizeProfile(profileData)
  } catch (err) {
    toast.error('Failed to load profile')
    console.error('Failed to load profile:', err)
  }
})
</script>
