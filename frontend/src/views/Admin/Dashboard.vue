<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />
    <div class="flex-1 py-8 px-6 sm:px-8">
      <div class="max-w-7xl mx-auto">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h1 class="text-3xl font-display text-text-primary mb-2">Admin Dashboard</h1>
          <p class="text-text-muted">Platform overview and statistics</p>
        </div>
        <Button variant="outline" @click="router.push('/dashboard')">
          Back to Dashboard
        </Button>
      </div>
 
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Card
          title="Total Users"
          :value="stats.total_users"
          icon="Users"
          :loading="loading"
          color="primary"
        />
        
        <Card
          title="Active Users (7 days)"
          :value="stats.active_users"
          icon="Activity"
          :loading="loading"
          color="success"
        />
        
        <Card
          title="Admin Users"
          :value="stats.admin_users"
          icon="Shield"
          :loading="loading"
          color="secondary"
        />
        
        <Card
          title="Profiles Complete"
          :value="stats.profiles_completed"
          icon="CheckCircle"
          :loading="loading"
          color="info"
        />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card
          title="New Users Today"
          :value="stats.new_users_today"
          icon="UserPlus"
          :loading="loading"
          color="warning"
        />
        
        <Card
          title="New Users This Week"
          :value="stats.new_users_week"
          icon="Users"
          :loading="loading"
          color="primary"
        />
        
        <Card
          title="Avg Completion"
          :value="`${Math.round(stats.avg_completion)}%`"
          icon="TrendingUp"
          :loading="loading"
          color="success"
        />
        
        <div class="bg-bg-secondary border border-border-subtle rounded-lg p-6">
          <h3 class="font-display text-xl text-text-primary mb-4">Last Updated</h3>
          <p class="text-text-tertiary">{{ formatDate(stats.last_updated) }}</p>
        </div>
      </div>
 
      <div class="mt-8">
        <div class="flex flex-wrap items-center gap-4">
          <Button variant="primary" @click="fetchStats">
            <RefreshCw class="w-4 h-4" />
            Refresh Statistics
          </Button>
          
          <Button variant="outline" @click="router.push('/admin/users')">
            <Users class="w-4 h-4" />
            Manage Users
          </Button>

          <Button variant="outline" @click="router.push('/admin/settings')">
            <Settings class="w-4 h-4" />
            Theme Settings
          </Button>
        </div>
      </div>
      </div>
    </div>
  </div>
</template>
 
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import AppHeader from '@/components/layout/AppHeader.vue'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import { RefreshCw, Users, CheckCircle, Activity, Shield, UserPlus, TrendingUp, Settings } from 'lucide-vue-next'
import { adminApi } from '@/api/admin'
import type { Statistics } from '@/api/admin'
 
const router = useRouter()
const loading = ref(false)
const stats = ref<Statistics>({
  total_users: 0,
  active_users: 0,
  admin_users: 0,
  profiles_completed: 0,
  new_users_today: 0,
  new_users_week: 0,
  avg_completion: 0,
  last_updated: '',
})
 
const fetchStats = async () => {
  loading.value = true
  try {
    const response = await adminApi.getStatistics()
    stats.value = response
    toast.success('Statistics updated successfully')
  } catch (err: any) {
    toast.error('Failed to fetch statistics')
    console.error('Failed to fetch statistics:', err)
  } finally {
    loading.value = false
  }
}
 
const formatDate = (date: string) => {
  const d = new Date(date)
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
 
onMounted(() => {
  fetchStats()
})
</script>
