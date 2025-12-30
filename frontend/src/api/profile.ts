import apiClient from '@/api/client'

export interface MediaFile {
  id: number
  file_name: string
  original_name: string
  file_path: string
  thumbnail_path?: string
  mime_type: string
  file_size: number
  width?: number
  height?: number
  file_type: string
  is_public: boolean
}

export interface ProfileResponse {
  id: number
  user_id: number
  username: string
  first_name?: string | null
  last_name?: string | null
  display_name?: string | null
  avatar?: MediaFile | null
  cover_photo?: MediaFile | null
  phone?: string | null
  website?: string | null
  bio?: string | null
  location?: string | null
  city?: string | null
  state?: string | null
  country?: string | null
  timezone?: string | null
  birthday?: string | null
  gender?: string | null
  pronouns?: string | null
  language?: string | null
  job_title?: string | null
  company?: string | null
  industry?: string | null
  linkedin_url?: string | null
  portfolio_url?: string | null
  interests?: string[]
  skills?: string[]
  custom_fields?: Record<string, unknown>
  created_at: string
  updated_at: string
}

export type ProfileUpdateRequest = Partial<
  Pick<
    ProfileResponse,
    | 'first_name'
    | 'last_name'
    | 'display_name'
    | 'phone'
    | 'website'
    | 'bio'
    | 'location'
    | 'city'
    | 'state'
    | 'country'
    | 'timezone'
    | 'birthday'
    | 'gender'
    | 'pronouns'
    | 'language'
    | 'job_title'
    | 'company'
    | 'industry'
    | 'linkedin_url'
    | 'portfolio_url'
    | 'interests'
    | 'skills'
    | 'custom_fields'
  >
>

export const profileApi = {
  async getMyProfile(): Promise<ProfileResponse> {
    const response = await apiClient.get('/profile')
    return response.data?.data
  },
  async updateProfile(payload: ProfileUpdateRequest): Promise<ProfileResponse> {
    const response = await apiClient.put('/profile', payload)
    return response.data?.data
  },
}
