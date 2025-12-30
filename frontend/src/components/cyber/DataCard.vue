<template>
  <div class="bg-cyber-gray/30 backdrop-blur-sm border border-cyber-cyan/30 rounded-cyber p-6">
    <div class="flex items-start gap-4">
      <div
        v-if="icon"
        :class="[
          'flex items-center justify-center rounded-cyber w-12 h-12',
          color === 'cyan' && 'bg-cyber-cyan/20 border-cyber-cyan/30 text-cyber-cyan',
          color === 'green' && 'bg-cyber-green/20 border-cyber-green/30 text-cyber-green',
          color === 'pink' && 'bg-cyber-pink/20 border-cyber-pink/30 text-cyber-pink',
          color === 'purple' && 'bg-cyber-purple/20 border-cyber-purple/30 text-cyber-purple',
          color === 'orange' && 'bg-cyber-orange/20 border-cyber-orange/30 text-cyber-orange',
          color === 'yellow' && 'bg-cyber-yellow/20 border-cyber-yellow/30 text-cyber-yellow',
        ]"
      >
        <component :is="iconComponent" />
      </div>
      
      <div class="flex-1">
        <div class="text-2xl font-cyber text-cyber-light-gray mb-2">{{ title }}</div>
        <div
          v-if="loading"
          class="text-4xl font-mono-cyber text-cyber-cyan"
        >
          Loading...
        </div>
        <div
          v-else
          :class="[
            'text-4xl font-bold',
            (typeof value === 'number' && 'text-cyber-cyan'),
            (typeof value === 'string' && 'text-cyber-white'),
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

interface Props {
  title: string
  value: number | string
  icon?: any
  loading?: boolean
  color?: 'cyan' | 'green' 'pink' | 'purple' | 'orange' | 'yellow'
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  color: 'cyan',
})

const iconComponents: Record<string, Component> = {
  Users: () => import('@/components/icons/lucide/Users.vue'),
  Activity: () => import('@/components/icons/lucide/Activity.vue'),
  Shield: () => import('@/components/icons/lucide/Shield.vue'),
  CheckCircle: () => import('@/components/icons/lucide/CheckCircle.vue'),
  UserPlus: () => import('@/components/icons/lucide/UserPlus.vue'),
  TrendingUp: () => import('@/components/icons/lucide/TrendingUp.vue'),
}

const iconComponent = computed(() => {
  if (!props.icon) return null
  return iconComponents[props.icon as string]
})
</script>
