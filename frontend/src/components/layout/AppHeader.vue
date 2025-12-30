<template>
  <header class="bg-bg-secondary border-b border-border-subtle sticky top-0 z-50">
    <div class="container-lg">
      <div class="flex items-center justify-between h-16">
        <!-- Logo/Brand -->
        <div class="flex items-center space-x-3">
          <router-link to="/" class="flex items-center space-x-3 group">
            <div class="w-8 h-8 bg-primary-500 rounded-lg flex items-center justify-center transition-all duration-200 group-hover:bg-primary-600">
              <span class="text-white font-display font-bold text-lg">{{ brand.shortName }}</span>
            </div>
            <span class="font-display font-semibold text-xl text-text-primary group-hover:text-primary-500 transition-colors duration-200">
              {{ brand.name }}
            </span>
          </router-link>
        </div>

        <!-- Desktop Navigation -->
        <nav class="hidden md:flex items-center space-x-1">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="px-4 py-3 rounded-md text-sm font-medium transition-all duration-200"
            :class="isActive(item.path)
              ? 'text-primary-500 bg-primary-500/10'
              : 'text-text-secondary hover:text-text-primary hover:bg-bg-tertiary'"
          >
            {{ item.label }}
          </router-link>
        </nav>

        <!-- User Menu / Auth Buttons -->
        <div class="flex items-center space-x-3">
          <!-- Authenticated User -->
          <div v-if="authStore.isAuthenticated" class="flex items-center space-x-3">
            <span class="hidden sm:block text-sm text-text-secondary">
              {{ authStore.user?.username }}
            </span>
            <button
              @click="handleLogout"
              class="btn-ghost text-sm px-4 py-2 min-h-[44px] min-w-[44px]"
            >
              Sign Out
            </button>
          </div>

          <!-- Unauthenticated User - Desktop -->
          <div v-else class="hidden md:flex items-center space-x-2">
            <router-link
              to="/login"
              class="btn-ghost text-sm px-4 py-2 min-h-[44px] min-w-[44px]"
            >
              Sign In
            </router-link>
            <router-link
              to="/register"
              class="btn-primary text-sm px-4 py-2 min-h-[44px] min-w-[44px]"
            >
              Get Started
            </router-link>
          </div>

          <!-- Mobile Menu Button -->
          <button
            @click="toggleMobileMenu"
            class="md:hidden p-2 rounded-md text-text-secondary hover:text-text-primary hover:bg-bg-tertiary transition-colors duration-200 min-h-[44px] min-w-[44px]"
            aria-label="Toggle menu"
            aria-expanded="mobileMenuOpen"
          >
            <MenuIcon v-if="!mobileMenuOpen" class="w-6 h-6" />
            <XIcon v-else class="w-6 h-6" />
          </button>
        </div>
      </div>

      <!-- Mobile Navigation Menu -->
      <nav
        v-if="mobileMenuOpen"
        class="md:hidden border-t border-border-subtle py-4"
      >
        <div class="space-y-2">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            @click="mobileMenuOpen = false"
            class="block px-4 py-3 rounded-md text-sm font-medium transition-all duration-200"
            :class="isActive(item.path)
              ? 'text-primary-500 bg-primary-500/10'
              : 'text-text-secondary hover:text-text-primary hover:bg-bg-tertiary'"
          >
            {{ item.label }}
          </router-link>
          
          <!-- Mobile Auth Buttons -->
          <div v-if="!authStore.isAuthenticated" class="pt-4 space-y-2">
            <router-link
              to="/login"
              @click="mobileMenuOpen = false"
              class="block w-full text-center px-4 py-3 rounded-md text-sm font-medium text-primary-500 hover:bg-bg-tertiary transition-all duration-200 min-h-[44px]"
            >
              Sign In
            </router-link>
            <router-link
              to="/register"
              @click="mobileMenuOpen = false"
              class="block w-full text-center px-4 py-3 rounded-md text-sm font-medium bg-primary-500 text-white hover:bg-primary-600 transition-all duration-200 min-h-[44px]"
            >
              Get Started
            </router-link>
          </div>
        </div>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { appConfig } from '@/config/appConfig'

// Simple SVG icons as components
const MenuIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16" /></svg>'
}

const XIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>'
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const mobileMenuOpen = ref(false)

const brand = appConfig.brand
const navItems = appConfig.navigation.header

const isActive = (path: string) => {
  return route.path === path
}

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
