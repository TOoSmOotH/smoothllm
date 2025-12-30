<template>
  <div :class="cardClasses">
    <div class="flex items-start gap-3 sm:gap-4">
      <div
        v-if="icon"
        :class="[
          'flex items-center justify-center rounded-lg w-10 h-10 sm:w-12 sm:h-12 flex-shrink-0',
          color === 'primary' && 'bg-primary-500/10 text-primary-500',
          color === 'secondary' && 'bg-secondary-500/10 text-secondary-500',
          color === 'success' && 'bg-success-500/10 text-success-500',
          color === 'warning' && 'bg-warning-500/10 text-warning-500',
          color === 'error' && 'bg-error-500/10 text-error-500',
          color === 'info' && 'bg-info-500/10 text-info-500',
        ]"
      >
        <component :is="iconComponent" class="w-5 h-5 sm:w-6 sm:h-6" />
      </div>
      
      <div class="flex-1 min-w-0">
        <div class="text-base sm:text-lg font-semibold text-text-secondary mb-1 sm:mb-2">{{ title }}</div>
        <div
          v-if="loading"
          class="text-2xl sm:text-4xl font-sans text-primary-500"
        >
          Loading...
        </div>
        <div
          v-else
          :class="[
            'text-2xl sm:text-4xl font-bold',
            (typeof value === 'number' && 'text-primary-500'),
            (typeof value === 'string' && 'text-text-primary'),
          ]"
        >
          {{ value }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, type Component } from 'vue'
import { Users, Activity, Shield, CheckCircle, UserPlus, TrendingUp } from 'lucide-vue-next'

interface Props {
  title: string
  value: number | string
  icon?: any
  loading?: boolean
  color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info'
  variant?: 'default' | 'elevated'
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  color: 'primary',
  variant: 'default',
})

const cardClasses = computed(() => {
  const baseClasses = [
    'transition-all',
    'duration-200',
    'rounded-lg',
    'p-6',
  ]

  const variantClasses = {
    default: [
      'bg-bg-secondary',
      'border',
      'border-border-subtle',
      'shadow-sm',
      'hover:shadow-md',
    ],
    elevated: [
      'bg-bg-secondary',
      'border',
      'border-border-default',
      'shadow-md',
      'hover:shadow-lg',
    ],
  }

  return [
    ...baseClasses,
    ...variantClasses[props.variant],
  ].join(' ')
})

const iconComponents: Record<string, Component> = {
  Users,
  Activity,
  Shield,
  CheckCircle,
  UserPlus,
  TrendingUp,
}

const iconComponent = computed(() => {
  if (!props.icon) return null
  return iconComponents[props.icon as string]
})
</script>
