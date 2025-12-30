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
    <span class="relative z-10">
      <slot />
    </span>
    <div v-if="variant !== 'ghost'" class="absolute inset-0 bg-gradient-to-r opacity-0 hover:opacity-20 transition-opacity duration-300" :class="gradientClass"></div>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'accent' | 'destructive' | 'ghost' | 'outline'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  glow?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  glow: true,
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
    'font-cyber',
    'font-semibold',
    'rounded-cyber',
    'transition-all',
    'duration-300',
    'transform',
    'hover:scale-105',
    'active:scale-95',
    'overflow-hidden',
    'border',
  ]

  const sizeClasses = {
    sm: ['px-4', 'py-2', 'text-sm'],
    md: ['px-6', 'py-3', 'text-base'],
    lg: ['px-8', 'py-4', 'text-lg'],
  }

  const variantClasses = {
    primary: [
      'bg-gradient-to-r', 'from-cyber-cyan', 'to-cyan-600',
      'text-cyber-black',
      'border-cyber-cyan/50',
      props.glow ? 'shadow-cyber-cyan' : '',
    ],
    secondary: [
      'bg-gradient-to-r', 'from-cyber-pink', 'to-pink-600',
      'text-white',
      'border-cyber-pink/50',
      props.glow ? 'shadow-cyber-pink' : '',
    ],
    accent: [
      'bg-gradient-to-r', 'from-cyber-purple', 'to-purple-600',
      'text-white',
      'border-cyber-purple/50',
      props.glow ? 'shadow-cyber-purple' : '',
    ],
    destructive: [
      'bg-gradient-to-r', 'from-red-500', 'to-red-600',
      'text-white',
      'border-red-500/50',
      'shadow-red-500/50',
    ],
    ghost: [
      'bg-transparent',
      'text-cyber-cyan',
      'border-cyber-cyan/30',
      'hover:bg-cyber-cyan/10',
    ],
    outline: [
      'bg-transparent',
      'text-cyber-cyan',
      'border-cyber-cyan',
      'hover:bg-cyber-cyan/10',
    ],
  }

  const stateClasses = []
  if (props.disabled) {
    stateClasses.push('opacity-50', 'cursor-not-allowed')
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

const gradientClass = computed(() => {
  const gradients = {
    primary: 'from-cyan-400 to-cyber-cyan',
    secondary: 'from-pink-400 to-cyber-pink',
    accent: 'from-purple-400 to-cyber-purple',
    destructive: 'from-red-400 to-red-500',
    ghost: '',
    outline: '',
  }
  return gradients[props.variant]
})

const handleClick = (event: MouseEvent) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
button {
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

button:hover {
  animation: cyber-pulse 2s ease-in-out infinite;
}

@keyframes cyber-pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}
</style>