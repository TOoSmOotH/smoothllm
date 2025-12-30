import apiClient from '@/api/client'

export interface ThemeResponse {
  theme: string
}

export interface RegistrationSettings {
  registration_enabled: boolean
  auto_approve_new_users: boolean
}

export const settingsApi = {
  async getTheme(): Promise<string> {
    const response = await apiClient.get('/settings/theme')
    return response.data?.data?.theme
  },
  async setTheme(theme: string): Promise<void> {
    await apiClient.put('/admin/settings/theme', { theme })
  },
  async getRegistrationSettings(): Promise<RegistrationSettings> {
    const response = await apiClient.get('/admin/settings/registration')
    return response.data?.data
  },
  async updateRegistrationSettings(payload: RegistrationSettings): Promise<void> {
    await apiClient.put('/admin/settings/registration', payload)
  },
}
