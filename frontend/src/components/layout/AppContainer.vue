<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <!-- Header -->
    <AppHeader v-if="showHeader" />

    <!-- Main Content -->
    <main :class="mainClasses" :style="mainStyle">
      <slot />
    </main>

    <!-- Footer -->
    <AppFooter v-if="showFooter" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppHeader from './AppHeader.vue'
import AppFooter from './AppFooter.vue'

interface Props {
  showHeader?: boolean
  showFooter?: boolean
  containerSize?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | 'full'
  padding?: 'none' | 'sm' | 'md' | 'lg' | 'xl'
}

const props = withDefaults(defineProps<Props>(), {
  showHeader: true,
  showFooter: true,
  containerSize: 'lg',
  padding: 'lg',
})

const mainClasses = computed(() => {
  const classes = ['flex-1']

  // Container size
  const containerClass = `container-${props.containerSize}`
  if (props.containerSize !== 'full') {
    classes.push(containerClass)
  }

  return classes.join(' ')
})

const mainStyle = computed(() => {
  // Padding using design system spacing variables
  const paddingMap = {
    none: { paddingTop: '0', paddingBottom: '0' },
    sm: { paddingTop: 'var(--space-4)', paddingBottom: 'var(--space-4)' },
    md: { paddingTop: 'var(--space-8)', paddingBottom: 'var(--space-8)' },
    lg: { paddingTop: 'var(--space-12)', paddingBottom: 'var(--space-12)' },
    xl: { paddingTop: 'var(--space-16)', paddingBottom: 'var(--space-16)' },
  }
  return paddingMap[props.padding]
})
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
