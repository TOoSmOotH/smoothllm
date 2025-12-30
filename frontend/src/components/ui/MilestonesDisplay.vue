<template>
  <div class="bg-bg-secondary border border-border-subtle rounded-lg p-4 sm:p-6">
    <div class="flex items-center justify-between mb-4 sm:mb-6">
      <h3 class="font-display text-lg sm:text-xl text-text-primary">Milestones</h3>
      <Button variant="ghost" size="sm" @click="handleRefresh" :loading="loading" class="min-h-[36px] min-w-[36px]">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
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

    <div v-else-if="milestones.length === 0" class="text-center py-8 text-text-muted">
      <p class="font-mono">No milestones available</p>
    </div>

    <div v-else class="space-y-3 sm:space-y-4">
      <div
        v-for="milestone in sortedMilestones"
        :key="milestone.id"
        :class="[
          'relative border rounded-md p-3 sm:p-4 transition-all duration-200',
          milestone.is_achieved
            ? 'border-success-500/50 bg-success-500/10'
            : milestone.is_next
            ? 'border-primary-500/50 bg-primary-500/10 ring-2 ring-primary-500/20'
            : 'border-border-default bg-bg-tertiary'
        ]"
      >
        <div class="flex items-start gap-3 sm:gap-4">
          <div class="flex-shrink-0">
            <div
              :class="[
                'w-10 h-10 sm:w-12 sm:h-12 rounded-lg flex items-center justify-center text-xl sm:text-2xl',
                milestone.is_achieved
                  ? 'bg-success-500 text-white'
                  : milestone.is_next
                  ? 'bg-primary-500 text-white'
                  : 'bg-bg-elevated text-text-muted'
              ]"
            >
              <span v-if="milestone.is_achieved">✓</span>
              <span v-else-if="milestone.is_next">→</span>
              <span v-else>○</span>
            </div>
          </div>

          <div class="flex-grow min-w-0">
            <div class="flex items-start justify-between gap-2">
              <div>
                <h4 class="font-sans text-base sm:text-lg text-text-primary mb-1">{{ milestone.name }}</h4>
                <p class="text-xs sm:text-sm text-text-muted font-mono mb-2">
                  {{ milestone.description }}
                </p>
              </div>
              <div v-if="milestone.icon_url" class="flex-shrink-0">
                <img :src="milestone.icon_url" :alt="milestone.name" class="w-10 h-10 sm:w-12 sm:h-12 rounded-lg" />
              </div>
            </div>

            <div class="flex items-center gap-3 sm:gap-4 text-xs sm:text-sm">
              <span class="text-text-tertiary font-mono">
                Threshold: <span class="text-primary-500 font-semibold">{{ milestone.threshold }} pts</span>
              </span>
              <span class="text-text-tertiary font-mono">
                Reward: <span class="text-secondary-500 font-semibold">{{ milestone.reward_value }}</span>
              </span>
            </div>
          </div>
        </div>

        <div v-if="milestone.is_achieved" class="absolute top-2 right-2">
          <span class="bg-success-500 text-white text-xs font-medium px-2 py-1 rounded">
            Achieved!
          </span>
        </div>

        <div v-else-if="milestone.is_next" class="absolute top-2 right-2">
          <span class="bg-primary-500 text-white text-xs font-medium px-2 py-1 rounded">
            Next
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import Button from '@/components/ui/Button.vue'
import { useCompletion } from '@/composables/useCompletion'
const { milestones, loading, fetchMilestones, completionScore, getNextMilestone, getAchievedMilestones } = useCompletion()

const sortedMilestones = computed(() => {
  const nextMilestone = getNextMilestone()
  const achievedMilestones = getAchievedMilestones()

  return milestones.value.map((milestone) => ({
    ...milestone,
    is_achieved: achievedMilestones.some((m) => m.id === milestone.id),
    is_next: nextMilestone?.id === milestone.id,
  }))
})

const handleRefresh = async () => {
  try {
    await fetchMilestones()
    toast.success('Milestones refreshed successfully')
  } catch (err) {
    console.error('Failed to refresh milestones:', err)
  }
}

onMounted(() => {
  fetchMilestones()
})
</script>
