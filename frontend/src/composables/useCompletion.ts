import { ref, computed } from 'vue'
import { toast } from 'vue-sonner'
import {
  CompletionScoreResponse,
  MilestoneResponse,
  LeaderboardEntry,
  getCompletionScore,
  recalculateCompletionScore,
  getMilestones,
  getLeaderboard,
} from '@/api/completion'

export function useCompletion() {
  const completionScore = ref<CompletionScoreResponse | null>(null)
  const milestones = ref<MilestoneResponse[]>([])
  const leaderboard = ref<LeaderboardEntry[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const percentage = computed(() => completionScore.value?.percentage || 0)
  const isComplete = computed(() => completionScore.value?.is_complete || false)
  const completedCount = computed(() => completionScore.value?.completed_fields.length || 0)
  const totalCount = computed(() => {
    return (
      (completionScore.value?.completed_fields.length || 0) +
      (completionScore.value?.missing_fields.length || 0)
    )
  })

  const fetchCompletionScore = async () => {
    loading.value = true
    error.value = null
    try {
      completionScore.value = await getCompletionScore()
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch completion score'
      toast.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const recalculateScore = async () => {
    loading.value = true
    error.value = null
    try {
      completionScore.value = await recalculateCompletionScore()
      toast.success('Profile score recalculated successfully')
    } catch (err: any) {
      error.value = err.message || 'Failed to recalculate score'
      toast.error(error.value)
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchMilestones = async () => {
    loading.value = true
    error.value = null
    try {
      milestones.value = await getMilestones()
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch milestones'
      toast.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const fetchLeaderboard = async (limit?: number) => {
    loading.value = true
    error.value = null
    try {
      leaderboard.value = await getLeaderboard(limit)
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch leaderboard'
      toast.error(error.value)
    } finally {
      loading.value = false
    }
  }

  const getCategoryPercentage = (category: string): number => {
    if (!completionScore.value) return 0
    const categoryInfo = completionScore.value.category_breakdown[category]
    if (!categoryInfo) return 0
    return categoryInfo.total > 0 ? (categoryInfo.points / categoryInfo.total) * 100 : 0
  }

  const getNextMilestone = (): MilestoneResponse | null => {
    if (!completionScore.value || !milestones.value.length) return null

    const next = milestones.value.find(
      (m) => m.threshold > completionScore.value!.score && m.is_active
    )

    return next || null
  }

  const getAchievedMilestones = (): MilestoneResponse[] => {
    if (!completionScore.value || !milestones.value.length) return []

    return milestones.value.filter(
      (m) => m.threshold <= completionScore.value!.score && m.is_active
    )
  }

  return {
    completionScore,
    milestones,
    leaderboard,
    loading,
    error,
    percentage,
    isComplete,
    completedCount,
    totalCount,
    fetchCompletionScore,
    recalculateScore,
    fetchMilestones,
    fetchLeaderboard,
    getCategoryPercentage,
    getNextMilestone,
    getAchievedMilestones,
  }
}
