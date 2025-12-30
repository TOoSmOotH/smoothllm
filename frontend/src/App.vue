<template>
  <div id="app" class="min-h-screen bg-bg-primary">
    <RouterView />
    <Toaster :theme="toastTheme" rich-colors close-button />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterView } from 'vue-router'
import { Toaster } from 'vue-sonner'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const toastTheme = computed(() => {
  const darkThemes = new Set([
    'warm-editorial-dark',
    'gothic-punk',
    'cyberpunk',
    'vaporware',
    'cassette-futurism-dark',
  ])
  return darkThemes.has(themeStore.theme) ? 'dark' : 'light'
})
</script>

<style scoped>
#app {
  position: relative;
  overflow-x: hidden;
}
</style>
