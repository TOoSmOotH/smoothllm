<template>
  <div class="bg-cyber-gray/30 backdrop-blur-sm border border-cyber-cyan/30 rounded-cyber p-6">
    <div class="flex items-center justify-between mb-6">
      <h3 class="font-cyber text-xl text-cyber-cyan">Milestones</h3>
      <CyberButton variant="ghost" size="sm" @click="handleRefresh" :loading="loading">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
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

    <div v-else-if="milestones.length === 0" class="text-center py-8 text-cyber-light-gray/60">
      <p class="font-mono-cyber">No milestones available</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="milestone in sortedMilestones"
        :key="milestone.id"
        :class="[
          'relative border rounded-cyber p-4 transition-all duration-300',
          milestone.is_achieved
            ? 'border-cyber-green/50 bg-cyber-green/10'
            : milestone.is_next
            ? 'border-cyber-cyan/50 bg-cyber-cyan/10 ring-2 ring-cyber-cyan/30'
            : 'border-cyber-light-gray/20 bg-cyber-dark/30'
        ]"
      >
        <div class="flex items-start gap-4">
          <div class="flex-shrink-0">
            <div
              :class="[
                'w-12 h-12 rounded-lg flex items-center justify-center text-2xl',
                milestone.is_achieved
                  ? 'bg-cyber-green text-cyber-black'
                  : milestone.is_next
                  ? 'bg-cyber-cyan text-cyber-black'
                  : 'bg-cyber-light-gray/40 text-cyber-light-gray/60'
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
                <h4 class="font-cyber text-lg text-white mb-1">{{ milestone.name }}</h4>
                <p class="text-sm text-cyber-light-gray/60 font-mono-cyber mb-2">
                  {{ milestone.description }}
                </p>
              </div>
              <div v-if="milestone.icon_url" class="flex-shrink-0">
                <img :src="milestone.icon_url" :alt="milestone.name" class="w-12 h-12 rounded-lg" />
              </div>
            </div>

            <div class="flex items-center gap-4 text-sm">
              <span class="text-cyber-light-gray/80 font-mono-cyber">
                Threshold: <span class="text-cyber-cyan font-semibold">{{ milestone.threshold }} pts</span>
              </span>
              <span class="text-cyber-light-gray/80 font-mono-cyber">
                Reward: <span class="text-cyber-purple font-semibold">{{ milestone.reward_value }}</span>
              </span>
            </div>
          </div>
        </div>

        <div v-if="milestone.is_achieved" class="absolute top-2 right-2">
          <span class="bg-cyber-green text-cyber-black text-xs font-cyber px-2 py-1 rounded">
            Achieved!
          </span>
        </div>

        <div v-else-if="milestone.is_next" class="absolute top-2 right-2">
          <span class="bg-cyber-cyan text-cyber-black text-xs font-cyber px-2 py-1 rounded">
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
import CyberButton from '@/components/cyber/CyberButton.vue'
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
