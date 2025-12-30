<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />
    <div class="flex-1 flex items-center justify-center px-6 sm:px-8">
      <div class="relative z-10 text-center">
        <div class="bg-bg-secondary border border-border-subtle rounded-lg p-12">
          <h1 class="text-6xl font-display font-bold text-text-primary mb-4">
            404
          </h1>
          <h2 class="text-2xl font-display text-text-primary mb-4">Page Not Found</h2>
          <p class="text-text-muted mb-8 max-w-md">
            The page you're looking for doesn't exist or has been moved.
          </p>
          
          <div class="space-y-4">
            <Button
              variant="primary"
              size="lg"
              @click="router.push('/')"
            >
              Return Home
            </Button>
            
            <div>
              <router-link
                to="/login"
                class="text-primary-500 hover:text-primary-600 text-sm transition-colors duration-200"
              >
                or sign in to your account
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
 
<script setup lang="ts">
import { useRouter } from 'vue-router'
import AppHeader from '@/components/layout/AppHeader.vue'
import Button from '@/components/ui/Button.vue'
 
const router = useRouter()
</script>
