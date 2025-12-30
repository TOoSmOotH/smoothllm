import apiClient from './client'

export interface SocialLink {
  id: number
  user_id: number
  platform: string
  url: string
  visible: boolean
  order: number
  created_at: string
}

export interface CreateSocialLinkRequest {
  platform: string
  url: string
  visible?: boolean
  order?: number
}

export interface UpdateSocialLinkRequest {
  platform?: string
  url?: string
  visible?: boolean
  order?: number
}

export interface ReorderSocialLinksRequest {
  updates: Array<{
    id: number
    order: number
  }>
}

/**
 * Get all social links for current user
 */
export const getSocialLinks = async (): Promise<SocialLink[]> => {
  const response = await apiClient.get('/social')
  return response.data
}

/**
 * Create a new social link
 */
export const createSocialLink = async (data: CreateSocialLinkRequest): Promise<SocialLink> => {
  const response = await apiClient.post('/social', data)
  return response.data
}

/**
 * Update a social link (full update)
 */
export const updateSocialLink = async (id: number, data: UpdateSocialLinkRequest): Promise<SocialLink> => {
  const response = await apiClient.put(`/social/${id}`, data)
  return response.data
}

/**
 * Partially update a social link
 */
export const patchSocialLink = async (id: number, data: UpdateSocialLinkRequest): Promise<SocialLink> => {
  const response = await apiClient.patch(`/social/${id}`, data)
  return response.data
}

/**
 * Delete a social link
 */
export const deleteSocialLink = async (id: number): Promise<void> => {
  await apiClient.delete(`/social/${id}`)
}

/**
 * Reorder social links
 */
export const reorderSocialLinks = async (updates: ReorderSocialLinksRequest): Promise<SocialLink[]> => {
  const response = await apiClient.patch('/social/reorder', updates)
  return response.data
}
