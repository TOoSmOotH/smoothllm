<template>
  <div class="bg-cyber-gray/30 backdrop-blur-sm border border-cyber-cyan/30 rounded-cyber p-6">
    <div class="flex items-center justify-between mb-6">
      <h3 class="font-cyber text-xl text-cyber-cyan">Profile Completion</h3>
      <CyberButton
        v-if="!isComplete"
        variant="secondary"
        size="sm"
        @click="handleRecalculate"
        :loading="loading"
      >
        Recalculate
      </CyberButton>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="animate-spin text-cyber-cyan">
        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    </div>

    <div v-else-if="completionScore">
      <div class="mb-6">
        <div class="flex items-end justify-between mb-2">
          <div>
            <p class="text-4xl font-cyber font-bold text-cyber-cyan">
              {{ Math.round(percentage) }}%
            </p>
            <p class="text-sm text-cyber-light-gray/60 font-mono-cyber">
              {{ completedCount }} / {{ totalCount }} fields completed
            </p>
          </div>
          <div v-if="isComplete" class="flex items-center gap-2 text-cyber-green">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
            </svg>
            <span class="font-cyber text-lg">Complete!</span>
          </div>
        </div>

        <div class="relative h-4 bg-cyber-dark rounded-full overflow-hidden">
          <div
            class="absolute top-0 left-0 h-full transition-all duration-700 ease-out"
            :class="progressBarClass"
            :style="{ width: `${percentage}%` }"
          >
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent animate-pulse"></div>
          </div>
        </div>
      </div>

      <div v-if="nextRecommended.length > 0" class="mb-6">
        <p class="text-sm font-cyber text-cyber-light-gray/80 mb-3">
          Recommended next steps:
        </p>
        <div class="space-y-2">
          <div
            v-for="field in nextRecommended.slice(0, 3)"
            :key="field.field"
            class="flex items-center justify-between bg-cyber-dark/50 border border-cyber-cyan/20 rounded-cyber p-3"
          >
            <span class="text-white font-mono-cyber">{{ field.label }}</span>
            <span class="text-cyber-cyan font-cyber">+{{ field.points }} pts</span>
          </div>
        </div>
      </div>

      <div v-if="Object.keys(categoryBreakdown).length > 0" class="space-y-3">
        <p class="text-sm font-cyber text-cyber-light-gray/80">
          Category breakdown:
        </p>
        <div
          v-for="(info, category) in categoryBreakdown"
          :key="category"
          class="space-y-1"
        >
          <div class="flex items-center justify-between text-sm">
            <span class="text-cyber-light-gray font-mono-cyber capitalize">
              {{ formatCategoryName(category) }}
            </span>
            <span class="text-white font-cyber">
              {{ info.completed }} / {{ Math.round(info.total / info.max_points * info.completed) }} completed
            </span>
          </div>
          <div class="relative h-2 bg-cyber-dark rounded-full overflow-hidden">
            <div
              class="absolute top-0 left-0 h-full transition-all duration-500 ease-out bg-cyber-purple"
              :style="{ width: `${getCategoryPercentage(category)}%` }"
            ></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { toast } from 'vue-sonner'
import CyberButton from '@/components/cyber/CyberButton.vue'
import { useCompletion } from '@/composables/useCompletion'
const { completionScore, loading, recalculateScore, getCategoryPercentage } = useCompletion()

const percentage = computed(() => completionScore.value?.percentage || 0)
const isComplete = computed(() => completionScore.value?.is_complete || false)
const completedCount = computed(() => completionScore.value?.completed_fields.length || 0)
const totalCount = computed(() => {
  return (
    (completionScore.value?.completed_fields.length || 0) +
    (completionScore.value?.missing_fields.length || 0)
  )
})
const nextRecommended = computed(() => completionScore.value?.next_recommended || [])
const categoryBreakdown = computed(() => completionScore.value?.category_breakdown || {})

const progressBarClass = computed(() => {
  if (percentage.value >= 75) {
    return 'bg-gradient-to-r from-cyber-green to-green-500'
  } else if (percentage.value >= 50) {
    return 'bg-gradient-to-r from-cyber-cyan to-cyan-500'
  } else if (percentage.value >= 25) {
    return 'bg-gradient-to-r from-cyber-yellow to-yellow-500'
  } else {
    return 'bg-gradient-to-r from-cyber-pink to-pink-500'
  }
})

const formatCategoryName = (category: string): string => {
  const names: Record<string, string> = {
    'basic_info': 'Basic Info',
    'contact_info': 'Contact Info',
    'personal_info': 'Personal Info',
    'professional': 'Professional',
    'extras': 'Extras',
  }
  return names[category] || category
}

const handleRecalculate = async () => {
  try {
    await recalculateScore()
    toast.success('Profile score recalculated successfully')
  } catch (err) {
    console.error('Failed to recalculate score:', err)
  }
}
</script>
