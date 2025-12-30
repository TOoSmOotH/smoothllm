<template>
  <aside class="h-full flex flex-col">
    <!-- Logo -->
    <div class="h-16 flex items-center justify-center border-b border-border-subtle px-4">
      <router-link to="/" class="flex items-center space-x-2 group" @click="handleNavClick">
        <div class="w-8 h-8 bg-primary-500 rounded-lg flex items-center justify-center flex-shrink-0 transition-all duration-200 group-hover:bg-primary-600">
          <span class="text-white font-display font-bold text-lg">{{ brand.shortName }}</span>
        </div>
        <span
          v-if="!collapsed"
          class="font-display font-semibold text-lg text-text-primary group-hover:text-primary-500 transition-colors duration-200 hidden lg:block"
        >
          {{ brand.name }}
        </span>
      </router-link>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto py-4">
      <div class="px-3 space-y-1">
        <!-- Main Navigation -->
        <div v-if="!collapsed" class="px-3 mb-2 hidden lg:block">
          <span class="text-xs font-semibold text-text-muted uppercase tracking-wider">
            Main
          </span>
        </div>

        <router-link
          v-for="item in mainNavItems"
          :key="item.path"
          :to="item.path"
          @click="handleNavClick"
          class="flex items-center space-x-3 px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 group min-h-[44px]"
          :class="isActive(item.path)
            ? 'text-primary-500 bg-primary-500/10'
            : 'text-text-secondary hover:text-text-primary hover:bg-bg-tertiary'"
          :title="collapsed ? item.label : ''"
        >
          <component :is="item.iconComponent" class="w-5 h-5 flex-shrink-0" />
          <span v-if="!collapsed" class="hidden lg:block">{{ item.label }}</span>
        </router-link>

        <!-- Admin Navigation -->
        <template v-if="authStore.isAdmin">
          <div v-if="!collapsed" class="px-3 mb-2 mt-6 hidden lg:block">
            <span class="text-xs font-semibold text-text-muted uppercase tracking-wider">
              Admin
            </span>
          </div>

          <router-link
            v-for="item in adminNavItems"
            :key="item.path"
            :to="item.path"
            @click="handleNavClick"
            class="flex items-center space-x-3 px-3 py-2 rounded-md text-sm font-medium transition-all duration-200 group min-h-[44px]"
            :class="isActive(item.path)
              ? 'text-primary-500 bg-primary-500/10'
              : 'text-text-secondary hover:text-text-primary hover:bg-bg-tertiary'"
            :title="collapsed ? item.label : ''"
          >
            <component :is="item.iconComponent" class="w-5 h-5 flex-shrink-0" />
            <span v-if="!collapsed" class="hidden lg:block">{{ item.label }}</span>
          </router-link>
        </template>
      </div>
    </nav>

    <!-- Collapse Toggle (Desktop Only) -->
    <div class="p-3 border-t border-border-subtle hidden lg:block">
      <button
        @click="toggleCollapse"
        class="w-full flex items-center justify-center space-x-2 px-3 py-2 rounded-md text-sm font-medium text-text-secondary hover:text-text-primary hover:bg-bg-tertiary transition-all duration-200 min-h-[44px]"
        :title="collapsed ? 'Expand sidebar' : 'Collapse sidebar'"
      >
        <ChevronLeftIcon v-if="!collapsed" class="w-5 h-5" />
        <ChevronRightIcon v-else class="w-5 h-5" />
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { appConfig } from '@/config/appConfig'

// Define emits
const emit = defineEmits<{
  close: []
}>()

// Simple SVG icons as components
const HomeIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>'
}

const DashboardIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" /></svg>'
}

const UserIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>'
}

const SettingsIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>'
}

const UsersIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg>'
}

const DefaultIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="4" /></svg>'
}

const ChevronLeftIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" /></svg>'
}

const ChevronRightIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" /></svg>'
}

const route = useRoute()
const authStore = useAuthStore()

const collapsed = ref(false)

const brand = appConfig.brand

const iconMap: Record<string, { template: string }> = {
  home: HomeIcon,
  dashboard: DashboardIcon,
  profile: UserIcon,
  admin: DashboardIcon,
  users: UsersIcon,
  settings: SettingsIcon,
}

const resolveNavItems = (items: typeof appConfig.navigation.sidebar.main) => {
  return items.map((item) => ({
    ...item,
    iconComponent: iconMap[item.icon ?? ''] ?? DefaultIcon,
  }))
}

const mainNavItems = resolveNavItems(appConfig.navigation.sidebar.main)
const adminNavItems = resolveNavItems(appConfig.navigation.sidebar.admin)

const isActive = (path: string) => {
  return route.path === path
}

const toggleCollapse = () => {
  collapsed.value = !collapsed.value
}

const handleNavClick = () => {
  // On mobile, close sidebar after navigation
  if (window.innerWidth < 1024) {
    emit('close')
  }
}
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
