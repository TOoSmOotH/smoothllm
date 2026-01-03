<template>
  <div class="relative">
    <label
      v-if="label"
      :for="id"
      :class="labelClasses"
    >
      {{ label }}
    </label>
    
    <input
      :id="id"
      :type="type"
      :value="modelValue"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      :required="required"
      :autocomplete="autocomplete"
      :class="inputClasses"
      @input="handleInput"
      @focus="handleFocus"
      @blur="handleBlur"
      @keyup.enter="$emit('enter')"
    />
    
    <!-- Error message -->
    <p v-if="error" class="mt-1 text-xs text-error-500 font-medium">
      {{ error }}
    </p>
    
    <!-- Helper text -->
    <p v-if="helperText && !error" class="mt-1 text-xs text-text-tertiary">
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
  required?: boolean
  autocomplete?: string
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
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
    'font-sans',
    'transition-all',
    'duration-200',
    'bg-bg-secondary',
    'border',
    'border-border-default',
    'rounded-md',
    'text-text-primary',
    'placeholder-text-muted',
    'focus:outline-none',
    'focus:border-primary-500',
    'focus:ring-2',
    'focus:ring-primary-500/10',
    'disabled:opacity-50',
    'disabled:cursor-not-allowed',
  ]

  const sizeClasses = {
    sm: ['px-3', 'py-2', 'text-sm'],
    md: ['px-4', 'py-3', 'text-base'],
    lg: ['px-6', 'py-4', 'text-lg'],
  }

  const stateClasses = []
  if (props.error) {
    stateClasses.push(
      'border-error-500',
      'focus:border-error-500',
      'focus:ring-error-500/10'
    )
  }

  return [
    ...baseClasses,
    ...sizeClasses[props.size],
    ...stateClasses,
  ].filter(Boolean).join(' ')
})

const labelClasses = computed(() => {
  const baseClasses = [
    'block',
    'text-sm',
    'font-medium',
    'mb-2',
    'transition-colors',
    'duration-200',
  ]

  const stateClasses = []
  if (!props.error) {
    stateClasses.push('text-text-secondary')
  }
  if (props.error) {
    stateClasses.push('text-error-500')
  }

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
