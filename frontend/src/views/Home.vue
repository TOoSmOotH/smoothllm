<template>
  <AppContainer>
    <!-- Hero Section -->
    <div class="flex flex-col items-center justify-center min-h-[calc(100vh-200px)] py-12 sm:py-16 md:py-20">
      <div class="text-center max-w-4xl mx-auto px-4">
        <!-- Main Title -->
        <h1 class="font-display font-bold text-4xl sm:text-5xl md:text-6xl lg:text-7xl xl:text-8xl mb-4 sm:mb-6 text-text-primary">
          {{ home.title }}
        </h1>
        
        <p class="text-lg sm:text-xl md:text-2xl text-text-secondary font-display mb-6 sm:mb-8">
          {{ home.subtitle }}
        </p>
        
        <p class="text-text-tertiary text-sm sm:text-base md:text-lg mb-8 sm:mb-12 max-w-2xl mx-auto leading-relaxed">
          {{ home.description }}
        </p>
        
        <!-- CTA Buttons -->
        <div class="flex flex-col sm:flex-row gap-3 sm:gap-4 justify-center mb-12 sm:mb-16">
          <button
            v-for="cta in home.ctas"
            :key="cta.to"
            @click="router.push(cta.to)"
            :class="ctaClass(cta.variant)"
          >
            {{ cta.label }}
          </button>
        </div>
        
        <!-- Features Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 sm:gap-6 mt-12 sm:mt-16">
          <div v-for="feature in home.features" :key="feature.title" class="card p-4 sm:p-6 text-center">
            <div
              class="w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center mx-auto mb-3 sm:mb-4"
              :class="featureTone(feature.tone).bg"
            >
              <component
                :is="featureIcon(feature.icon)"
                class="w-6 h-6"
                :class="featureTone(feature.tone).text"
              />
            </div>
            <h3
              class="font-display font-semibold mb-2 text-base sm:text-lg"
              :class="featureTone(feature.tone).text"
            >
              {{ feature.title }}
            </h3>
            <p class="text-text-tertiary text-xs sm:text-sm leading-relaxed">
              {{ feature.description }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </AppContainer>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import AppContainer from '@/components/layout/AppContainer.vue'
import { appConfig } from '@/config/appConfig'

const router = useRouter()
const home = appConfig.home

const ctaClasses = {
  primary: 'btn-primary px-6 sm:px-8 py-3 sm:py-4 text-base sm:text-lg min-h-[44px]',
  secondary: 'btn-secondary px-6 sm:px-8 py-3 sm:py-4 text-base sm:text-lg min-h-[44px]',
  ghost: 'btn-ghost px-6 sm:px-8 py-3 sm:py-4 text-base sm:text-lg min-h-[44px]',
}

const ctaClass = (variant?: 'primary' | 'secondary' | 'ghost') => {
  return ctaClasses[variant ?? 'primary'] ?? ctaClasses.primary
}

const toneClasses = {
  primary: { bg: 'bg-primary-500/10', text: 'text-primary-500' },
  secondary: { bg: 'bg-secondary-500/10', text: 'text-secondary-500' },
  success: { bg: 'bg-success-500/10', text: 'text-success-500' },
}

const featureTone = (tone?: 'primary' | 'secondary' | 'success') => {
  return toneClasses[tone ?? 'primary'] ?? toneClasses.primary
}

const LockIcon = {
  template: '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>',
}

const SparklesIcon = {
  template: '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01" /></svg>',
}

const BoltIcon = {
  template: '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>',
}

const DefaultFeatureIcon = {
  template: '<svg fill="none" viewBox="0 0 24 24" stroke="currentColor"><circle cx="12" cy="12" r="4" /></svg>',
}

const featureIcons: Record<string, { template: string }> = {
  lock: LockIcon,
  sparkles: SparklesIcon,
  bolt: BoltIcon,
}

const featureIcon = (key: string) => {
  return featureIcons[key] ?? DefaultFeatureIcon
}
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
