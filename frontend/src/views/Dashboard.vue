<template>
  <AppLayout>
    <!-- User Info Card -->
    <div class="card mb-6 sm:mb-8">
      <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">User Profile</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <div>
          <p class="text-text-tertiary text-sm font-medium mb-1">Username</p>
          <p class="text-text-primary font-mono">{{ authStore.user?.username }}</p>
        </div>
        <div>
          <p class="text-text-tertiary text-sm font-medium mb-1">Email</p>
          <p class="text-text-primary font-mono">{{ authStore.user?.email }}</p>
        </div>
        <div>
          <p class="text-text-tertiary text-sm font-medium mb-1">Role</p>
          <p class="text-primary-500 font-mono uppercase font-medium">{{ authStore.user?.role }}</p>
        </div>
        <div>
          <p class="text-text-tertiary text-sm font-medium mb-1">Member Since</p>
          <p class="text-text-primary font-mono">{{ formatDate(authStore.user?.created_at) }}</p>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-6 sm:mb-8">
      <button
        @click="router.push('/profile')"
        class="card hover:border-primary-500/50 transition-all duration-200 group min-h-[100px]"
      >
        <div class="flex flex-col items-center justify-center h-20 sm:h-24 space-y-2">
          <div class="w-12 h-12 bg-primary-500/10 rounded-lg flex items-center justify-center group-hover:bg-primary-500/20 transition-colors duration-200">
            <svg class="w-6 h-6 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          <span class="text-text-primary font-medium">View Profile</span>
        </div>
      </button>
      
      <button
        @click="router.push('/profile/edit')"
        class="card hover:border-secondary-500/50 transition-all duration-200 group min-h-[100px]"
      >
        <div class="flex flex-col items-center justify-center h-20 sm:h-24 space-y-2">
          <div class="w-12 h-12 bg-secondary-500/10 rounded-lg flex items-center justify-center group-hover:bg-secondary-500/20 transition-colors duration-200">
            <svg class="w-6 h-6 text-secondary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <span class="text-text-primary font-medium">Edit Profile</span>
        </div>
      </button>
      
      <button
        v-if="authStore.isAdmin"
        @click="router.push('/admin')"
        class="card hover:border-primary-500/50 transition-all duration-200 group min-h-[100px]"
      >
        <div class="flex flex-col items-center justify-center h-20 sm:h-24 space-y-2">
          <div class="w-12 h-12 bg-primary-500/10 rounded-lg flex items-center justify-center group-hover:bg-primary-500/20 transition-colors duration-200">
            <svg class="w-6 h-6 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <span class="text-text-primary font-medium">Admin Panel</span>
        </div>
      </button>
    </div>

    <!-- System Status -->
    <div class="card">
      <h2 class="font-display font-semibold text-xl text-text-primary" style="margin-bottom: var(--space-6)">System Status</h2>
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
        <div class="text-center bg-bg-tertiary rounded-lg" style="padding: var(--space-4)">
          <div class="text-success-500 font-display font-semibold text-2xl mb-2">Online</div>
          <p class="text-text-tertiary text-sm">Backend Service</p>
        </div>
        <div class="text-center bg-bg-tertiary rounded-lg" style="padding: var(--space-4)">
          <div class="text-success-500 font-display font-semibold text-2xl mb-2">Active</div>
          <p class="text-text-tertiary text-sm">Database</p>
        </div>
        <div class="text-center bg-bg-tertiary rounded-lg" style="padding: var(--space-4)">
          <div class="text-success-500 font-display font-semibold text-2xl mb-2">Secure</div>
          <p class="text-text-tertiary text-sm">Authentication</p>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'

const router = useRouter()
const authStore = useAuthStore()

const formatDate = (dateString?: string) => {
  if (!dateString) return 'Unknown'
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>