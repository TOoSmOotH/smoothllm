import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'

export function useAuth() {
  const authStore = useAuthStore()
  const { user, token, loading, error, isAuthenticated, isAdmin } = storeToRefs(authStore)

  return {
    user,
    token,
    loading,
    error,
    isAuthenticated,
    isAdmin,
    authStore,
  }
}
