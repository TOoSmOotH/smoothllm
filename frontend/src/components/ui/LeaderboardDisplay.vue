<template>
  <div class="bg-bg-secondary border border-border-subtle rounded-lg p-4 sm:p-6">
    <div class="flex items-center justify-between mb-4 sm:mb-6">
      <h3 class="font-display text-lg sm:text-xl text-text-primary">Leaderboard</h3>
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

    <div v-else-if="leaderboard.length === 0" class="text-center py-8 text-text-muted">
      <p class="font-mono">No leaderboard entries yet</p>
    </div>

    <div v-else class="space-y-2 sm:space-y-3">
      <div
        v-for="(entry, index) in leaderboard"
        :key="entry.user_id"
        :class="[
          'flex items-center gap-3 sm:gap-4 bg-bg-tertiary border rounded-md p-3 sm:p-4 transition-all duration-200',
          isCurrentUser(entry.user_id)
            ? 'border-primary-500/50 ring-2 ring-primary-500/20'
            : 'border-border-default'
        ]"
      >
        <div
          :class="[
            'flex-shrink-0 w-8 h-8 sm:w-10 sm:h-10 rounded-lg flex items-center justify-center font-sans text-base sm:text-lg font-bold',
            getRankStyle(index)
          ]"
        >
          {{ entry.rank }}
        </div>

        <div class="flex-grow min-w-0">
          <p class="font-sans text-sm sm:text-base text-text-primary truncate">{{ entry.username }}</p>
          <p class="text-xs text-text-muted font-mono">
            {{ formatDate(entry.completed_at) }}
          </p>
        </div>

        <div class="flex-shrink-0 text-right">
          <p class="font-sans text-xl sm:text-2xl font-bold text-primary-500">{{ entry.score }}</p>
          <p class="text-xs text-text-muted font-mono">pts</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { toast } from 'vue-sonner'
import { useAuth } from '@/composables/useAuth'
import Button from '@/components/ui/Button.vue'
import { useCompletion } from '@/composables/useCompletion'
const { user } = useAuth()
const { leaderboard, loading, fetchLeaderboard } = useCompletion()

const currentUserId = computed(() => user.value?.id)

const leaderboard = computed(() => {
  return leaderboard.value.slice(0, 10)
})

const getRankStyle = (index: number) => {
  if (index === 0) {
    return 'bg-warning-500 text-white shadow-md'
  } else if (index === 1) {
    return 'bg-text-tertiary text-text-primary shadow-md'
  } else if (index === 2) {
    return 'bg-warning-600 text-white shadow-md'
  } else {
    return 'bg-bg-elevated text-text-muted'
  }
}

const isCurrentUser = (userId: number): boolean => {
  return userId === currentUserId.value
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return 'Today'
  if (days === 1) return 'Yesterday'
  if (days < 7) return `${days} days ago`
  if (days < 30) return `${Math.floor(days / 7)} weeks ago`
  return `${Math.floor(days / 30)} months ago`
}

const handleRefresh = async () => {
  try {
    await fetchLeaderboard(10)
    toast.success('Leaderboard refreshed successfully')
  } catch (err) {
    console.error('Failed to refresh leaderboard:', err)
  }
}

onMounted(() => {
  fetchLeaderboard(10)
})
</script>
