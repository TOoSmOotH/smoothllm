import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  AuthState,
  User,
  RegisterRequest,
  LoginRequest,
  AuthResponse,
  RegisterResponse,
} from '@/types/auth'
import apiClient from '@/api/client'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const isInitialized = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  // Actions
  const initialize = async () => {
    if (isInitialized.value) return

    const storedToken = localStorage.getItem('access_token')
    const storedRefreshToken = localStorage.getItem('refresh_token')

    if (storedToken && storedRefreshToken) {
      token.value = storedToken
      refreshToken.value = storedRefreshToken
      
      try {
        // Verify token by getting current user
        const response = await apiClient.get('/auth/me')
        user.value = response.data
      } catch (err) {
        // Token is invalid, clear everything
        clearAuth()
      }
    }

    isInitialized.value = true
  }

  const register = async (registerData: RegisterRequest): Promise<RegisterResponse> => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.post<RegisterResponse>('/auth/register', registerData)
      const authData = response.data

      if (authData.approved && authData.token && authData.refresh_token) {
        setAuth({
          token: authData.token,
          refresh_token: authData.refresh_token,
          user: authData.user,
        })
      }

      return authData
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Registration failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const login = async (loginData: LoginRequest): Promise<AuthResponse> => {
    loading.value = true
    error.value = null

    try {
      const response = await apiClient.post<AuthResponse>('/auth/login', loginData)
      const authData = response.data

      setAuth(authData)
      return authData
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      // Call logout endpoint if available
      await apiClient.post('/auth/logout')
    } catch (err) {
      // Ignore logout errors, just clear local auth
    } finally {
      clearAuth()
    }
  }

  const refreshTokens = async (): Promise<void> => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await apiClient.post<AuthResponse>('/auth/refresh', {
        refresh_token: refreshToken.value,
      })
      
      const authData = response.data
      setAuth(authData)
    } catch (err) {
      clearAuth()
      throw err
    }
  }

  const setAuth = (authData: AuthResponse) => {
    user.value = authData.user
    token.value = authData.token
    refreshToken.value = authData.refresh_token

    // Store in localStorage
    localStorage.setItem('access_token', authData.token)
    localStorage.setItem('refresh_token', authData.refresh_token)
  }

  const clearAuth = () => {
    user.value = null
    token.value = null
    refreshToken.value = null
    error.value = null

    // Clear localStorage
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  const updateUser = (userData: Partial<User>) => {
    if (user.value) {
      user.value = { ...user.value, ...userData }
    }
  }

  return {
    // State
    user,
    token,
    refreshToken,
    loading,
    error,
    isInitialized,

    // Getters
    isAuthenticated,
    isAdmin,

    // Actions
    initialize,
    register,
    login,
    logout,
    refreshTokens,
    setAuth,
    clearAuth,
    updateUser,
  }
}, {
  persist: {
    key: 'smoothweb-auth',
    storage: localStorage,
    paths: ['user', 'token', 'refreshToken'],
  },
})
