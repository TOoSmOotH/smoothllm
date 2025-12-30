export interface User {
  id: number
  email: string
  username: string
  role: 'admin' | 'user'
  status: 'active' | 'pending' | 'disabled'
  created_at: string
  updated_at: string
  last_active_at?: string
  email_verified_at?: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface AuthResponse {
  token: string
  refresh_token: string
  user: User
}

export interface RegisterResponse {
  token?: string
  refresh_token?: string
  user: User
  approved: boolean
  message?: string
}

export interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  isInitialized: boolean
  loading: boolean
  error: string | null
}
