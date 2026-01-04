<template>
  <button
    :class="buttonClasses"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <span v-if="loading" class="inline-block animate-spin mr-2">
      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </span>
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'destructive' | 'ghost' | 'outline'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const buttonClasses = computed(() => {
  const baseClasses = [
    'relative',
    'inline-flex',
    'items-center',
    'justify-center',
    'font-sans',
    'font-semibold',
    'rounded-md',
    'transition-all',
    'duration-200',
    'focus:outline-none',
    'focus:ring-2',
    'focus:ring-offset-2',
    'focus:ring-offset-bg-primary',
  ]

  const sizeClasses = {
    sm: ['px-4', 'py-2', 'text-sm'],
    md: ['px-6', 'py-3', 'text-base'],
    lg: ['px-8', 'py-4', 'text-lg'],
  }

  const variantClasses = {
    primary: [
      'bg-primary-500',
      'text-white',
      'hover:bg-primary-600',
      'focus:ring-primary-500',
      'shadow-sm',
      'hover:-translate-y-px',
      'active:translate-y-0',
    ],
    secondary: [
      'bg-transparent',
      'text-primary-500',
      'border',
      'border-primary-500',
      'hover:bg-primary-500',
      'hover:text-white',
      'focus:ring-primary-500',
    ],
    destructive: [
      'bg-error-500',
      'text-white',
      'hover:bg-error-600',
      'focus:ring-error-500',
      'shadow-sm',
      'hover:-translate-y-px',
      'active:translate-y-0',
    ],
    ghost: [
      'bg-transparent',
      'text-text-secondary',
      'hover:bg-bg-tertiary',
      'hover:text-text-primary',
      'focus:ring-text-secondary',
    ],
    outline: [
      'bg-transparent',
      'text-primary-500',
      'border',
      'border-primary-500',
      'hover:bg-primary-500/10',
      'focus:ring-primary-500',
    ],
  }

  const stateClasses = []
  if (props.disabled) {
    stateClasses.push('opacity-50', 'cursor-not-allowed', 'pointer-events-none')
  }
  if (props.loading) {
    stateClasses.push('opacity-75', 'cursor-wait')
  }

  return [
    ...baseClasses,
    ...sizeClasses[props.size],
    ...variantClasses[props.variant],
    ...stateClasses,
  ].filter(Boolean).join(' ')
})

const handleClick = (event: MouseEvent) => {
  emit('click', event)
}
</script>
