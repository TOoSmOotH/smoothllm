import axios from 'axios'

const rawBaseUrl = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'
const normalizedBaseUrl = rawBaseUrl.replace(/\/+$/, '')
const apiBase =
  normalizedBaseUrl.endsWith('/api/v1') ? normalizedBaseUrl : `${normalizedBaseUrl}/api/v1`
const API_BASE_URL = `${apiBase}/admin`

export interface Statistics {
  total_users: number
  active_users: number
  admin_users: number
  profiles_completed: number
  new_users_today: number
  new_users_week: number
  avg_completion: number
  last_updated: string
}

export const adminApi = {
  async getStatistics(): Promise<Statistics> {
    const response = await axios.get(`${API_BASE_URL}/stats`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    return response.data.data
  },

  async listUsers(params?: {
    page?: number
    limit?: number
    role?: string
    search?: string
  }): Promise<any> {
    const response = await axios.get(`${API_BASE_URL}/users`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
      params,
    })
    return response.data.data
  },

  async createUser(payload: {
    email: string
    username: string
    password: string
    role?: 'admin' | 'user'
  }): Promise<any> {
    const response = await axios.post(`${API_BASE_URL}/users`, payload, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
    return response.data.data
  },

  async deleteUser(userId: number): Promise<void> {
    await axios.delete(`${API_BASE_URL}/users/${userId}`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
  },

  async changeUserRole(userId: number, role: 'admin' | 'user'): Promise<void> {
    await axios.patch(`${API_BASE_URL}/users/${userId}/role`, { role }, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
  },

  async approveUser(userId: number): Promise<void> {
    await axios.patch(`${API_BASE_URL}/users/${userId}/approve`, null, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('access_token')}`,
      },
    })
  },
}
