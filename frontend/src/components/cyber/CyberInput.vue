<template>
  <div class="relative">
    <div class="relative">
      <input
        :id="id"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :class="inputClasses"
        @input="handleInput"
        @focus="handleFocus"
        @blur="handleBlur"
        @keyup.enter="$emit('enter')"
      />
      
      <!-- Corner accents -->
      <div class="absolute top-0 left-0 w-2 h-2 border-t-2 border-l-2 border-cyber-cyan/50"></div>
      <div class="absolute top-0 right-0 w-2 h-2 border-t-2 border-r-2 border-cyber-cyan/50"></div>
      <div class="absolute bottom-0 left-0 w-2 h-2 border-b-2 border-l-2 border-cyber-cyan/50"></div>
      <div class="absolute bottom-0 right-0 w-2 h-2 border-b-2 border-r-2 border-cyber-cyan/50"></div>
    </div>
    
    <!-- Label -->
    <label
      v-if="label"
      :for="id"
      :class="labelClasses"
    >
      {{ label }}
    </label>
    
    <!-- Error message -->
    <p v-if="error" class="mt-1 text-sm text-cyber-pink font-mono-cyber">
      {{ error }}
    </p>
    
    <!-- Helper text -->
    <p v-if="helperText && !error" class="mt-1 text-xs text-cyber-light-gray/60 font-mono-cyber">
      {{ helperText }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
  id?: string
  modelValue: string | number
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url'
  label?: string
  placeholder?: string
  error?: string
  helperText?: string
  disabled?: boolean
  readonly?: boolean
  variant?: 'default' | 'cyber' | 'neon'
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  variant: 'cyber',
  size: 'md',
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  focus: [event: FocusEvent]
  blur: [event: FocusEvent]
  enter: []
}>()

const isFocused = ref(false)

const inputClasses = computed(() => {
  const baseClasses = [
    'w-full',
    'font-mono-cyber',
    'transition-all',
    'duration-300',
    'bg-cyber-dark/50',
    'border',
    'rounded-cyber',
    'text-white',
    'placeholder-cyber-light-gray/50',
    'focus:outline-none',
    'relative',
    'z-10',
  ]

  const sizeClasses = {
    sm: ['px-3', 'py-2', 'text-sm'],
    md: ['px-4', 'py-3', 'text-base'],
    lg: ['px-5', 'py-4', 'text-lg'],
  }

  const variantClasses = {
    default: [
      'border-cyber-border',
      'focus:border-cyber-cyan',
      'focus:shadow-cyber-cyan',
    ],
    cyber: [
      'border-cyber-cyan/30',
      'focus:border-cyber-cyan',
      'focus:shadow-cyber-cyan',
      'backdrop-blur-cyber',
    ],
    neon: [
      'border-cyber-pink/30',
      'focus:border-cyber-pink',
      'focus:shadow-cyber-pink',
      'backdrop-blur-cyber',
    ],
  }

  const stateClasses = []
  if (props.disabled) stateClasses.push('opacity-50 cursor-not-allowed')
  if (props.error) stateClasses.push('border-cyber-pink focus:border-cyber-pink focus:shadow-cyber-pink')
  if (isFocused.value && !props.error) stateClasses.push('animate-cyber-pulse')

  return [
    ...baseClasses,
    ...sizeClasses[props.size],
    ...variantClasses[props.variant],
    ...stateClasses,
  ].filter(Boolean).join(' ')
})

const labelClasses = computed(() => {
  const baseClasses = [
    'block',
    'text-sm',
    'font-cyber',
    'mb-2',
    'transition-colors',
    'duration-300',
  ]

  const stateClasses = []
  if (!props.error) stateClasses.push('text-cyber-cyan')
  if (props.error) stateClasses.push('text-cyber-pink')

  return [...baseClasses, ...stateClasses].join(' ')
})

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

const handleFocus = (event: FocusEvent) => {
  isFocused.value = true
  emit('focus', event)
}

const handleBlur = (event: FocusEvent) => {
  isFocused.value = false
  emit('blur', event)
}
</script>

<style scoped>
input {
  background-image: 
    linear-gradient(90deg, transparent 0%, rgba(0, 243, 255, 0.1) 50%, transparent 100%);
  background-size: 0% 100%;
  background-repeat: no-repeat;
  background-position: center;
  transition: background-size 0.3s ease;
}

input:focus {
  background-size: 100% 100%;
}

input::placeholder {
  color: rgba(160, 160, 176, 0.5);
}

/* Custom scrollbar for number inputs */
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  opacity: 0.5;
  filter: invert(1);
}
</style>