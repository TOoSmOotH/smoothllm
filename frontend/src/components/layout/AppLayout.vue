<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />

    <!-- Main Content -->
    <main :class="mainClasses">
      <div class="container-2xl">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppHeader from '@/components/layout/AppHeader.vue'

const mainClasses = computed(() => {
  return 'flex-1 overflow-auto py-6'
})
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
