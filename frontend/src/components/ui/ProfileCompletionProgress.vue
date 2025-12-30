<template>
  <div class="bg-bg-secondary border border-border-subtle rounded-lg p-4 sm:p-6">
    <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 mb-4 sm:mb-6">
      <h3 class="font-display text-lg sm:text-xl text-text-primary">Profile Completion</h3>
      <Button
        v-if="!isComplete"
        variant="secondary"
        size="sm"
        @click="handleRecalculate"
        :loading="loading"
        class="min-h-[36px] min-w-[36px]"
      >
        Recalculate
      </Button>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="animate-spin text-primary-500">
        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    </div>

    <div v-else-if="completionScore">
      <div class="mb-4 sm:mb-6">
        <div class="flex flex-col sm:flex-row items-start sm:items-end justify-between gap-2 mb-2">
          <div>
            <p class="text-3xl sm:text-4xl font-display font-bold text-primary-500">
              {{ Math.round(percentage) }}%
            </p>
            <p class="text-xs sm:text-sm text-text-muted font-mono">
              {{ completedCount }} / {{ totalCount }} fields completed
            </p>
          </div>
          <div v-if="isComplete" class="flex items-center gap-2 text-success-500">
            <svg class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
            </svg>
            <span class="font-sans text-lg">Complete!</span>
          </div>
        </div>

        <div class="relative h-3 sm:h-4 bg-bg-tertiary rounded-full overflow-hidden">
          <div
            class="absolute top-0 left-0 h-full transition-all duration-700 ease-out"
            :class="progressBarClass"
            :style="{ width: `${percentage}%` }"
          >
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent animate-pulse"></div>
          </div>
        </div>
      </div>

      <div v-if="nextRecommended.length > 0" class="mb-4 sm:mb-6">
        <p class="text-xs sm:text-sm font-medium text-text-tertiary mb-2 sm:mb-3">
          Recommended next steps:
        </p>
        <div class="space-y-2">
          <div
            v-for="field in nextRecommended.slice(0, 3)"
            :key="field.field"
            class="flex items-center justify-between bg-bg-tertiary border border-border-default rounded-md p-2 sm:p-3"
          >
            <span class="text-text-primary font-mono">{{ field.label }}</span>
            <span class="text-primary-500 font-medium">+{{ field.points }} pts</span>
          </div>
        </div>
      </div>

      <div v-if="Object.keys(categoryBreakdown).length > 0" class="space-y-2 sm:space-y-3">
        <p class="text-xs sm:text-sm font-medium text-text-tertiary">
          Category breakdown:
        </p>
        <div
          v-for="(info, category) in categoryBreakdown"
          :key="category"
          class="space-y-1"
        >
          <div class="flex items-center justify-between text-xs sm:text-sm">
            <span class="text-text-muted font-mono capitalize">
              {{ formatCategoryName(category) }}
            </span>
            <span class="text-text-primary font-medium">
              {{ info.completed }} / {{ Math.round(info.total / info.max_points * info.completed) }} completed
            </span>
          </div>
          <div class="relative h-1.5 sm:h-2 bg-bg-tertiary rounded-full overflow-hidden">
            <div
              class="absolute top-0 left-0 h-full transition-all duration-500 ease-out bg-secondary-500"
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
import Button from '@/components/ui/Button.vue'
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
    return 'bg-success-500'
  } else if (percentage.value >= 50) {
    return 'bg-primary-500'
  } else if (percentage.value >= 25) {
    return 'bg-warning-500'
  } else {
    return 'bg-error-500'
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
