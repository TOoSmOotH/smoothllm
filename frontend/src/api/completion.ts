import apiClient from './client'

export interface CompletionFieldInfo {
  field: string
  label: string
  points: number
  is_filled: boolean
}

export interface CategoryInfo {
  completed: number
  total: number
  points: number
  max_points: number
}

export interface CompletionScoreResponse {
  user_id: number
  score: number
  max_score: number
  percentage: number
  is_complete: boolean
  completed_fields: string[]
  missing_fields: CompletionFieldInfo[]
  category_breakdown: Record<string, CategoryInfo>
  next_recommended: CompletionFieldInfo[]
}

export interface MilestoneResponse {
  id: number
  name: string
  description: string
  threshold: number
  reward_type: string
  reward_value: string
  icon_url?: string
  is_active: boolean
}

export interface LeaderboardEntry {
  user_id: number
  username: string
  score: number
  rank: number
  completed_at: string
}

/**
 * Get current user's profile completion score
 */
export const getCompletionScore = async (): Promise<CompletionScoreResponse> => {
  const response = await apiClient.get('/completion')
  return response.data?.data ?? response.data
}

/**
 * Recalculate profile completion score
 */
export const recalculateCompletionScore = async (): Promise<CompletionScoreResponse> => {
  const response = await apiClient.post('/completion/recalculate')
  return response.data?.data ?? response.data
}

/**
 * Get all available milestones
 */
export const getMilestones = async (): Promise<MilestoneResponse[]> => {
  const response = await apiClient.get('/completion/milestones')
  return response.data?.data ?? response.data
}

/**
 * Get completion leaderboard
 */
export const getLeaderboard = async (limit?: number): Promise<LeaderboardEntry[]> => {
  const params = limit ? { limit } : {}
  const response = await apiClient.get('/completion/leaderboard', { params })
  return response.data?.data ?? response.data
}
