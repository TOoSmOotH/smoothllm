<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />
    <div class="flex-1 py-8 px-6 sm:px-8">
      <div class="max-w-7xl mx-auto">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between mb-8">
        <div>
          <h1 class="text-3xl font-display text-text-primary mb-2">Admin Users</h1>
          <p class="text-text-muted">User management and control</p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <Button variant="primary" @click="showCreateModal = true">
            Add User
          </Button>
          <Button variant="outline" @click="router.push('/admin')">
            Back to Admin Dashboard
          </Button>
        </div>
      </div>
 
      <div class="mb-6">
        <Input
          v-model="searchQuery"
          placeholder="Search users by email or username..."
          class="max-w-md"
        />
      </div>
 
      <div class="mb-6">
        <div class="bg-bg-secondary border border-border-subtle rounded-lg p-4 mb-4">
          <h3 class="font-display text-lg text-text-primary mb-4">Filters</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <select
              v-model="filters.role"
              class="bg-bg-tertiary border border-border-default rounded-md text-text-primary p-3 font-sans focus:outline-none focus:border-primary-500"
            >
              <option value="">All Roles</option>
              <option value="admin">Admin</option>
              <option value="user">User</option>
            </select>
            <select
              v-model="filters.sort"
              class="bg-bg-tertiary border border-border-default rounded-md text-text-primary p-3 font-sans focus:outline-none focus:border-primary-500"
            >
              <option value="created_at">Newest</option>
              <option value="updated_at">Updated</option>
              <option value="email">Email</option>
              <option value="username">Username</option>
            </select>
          </div>
        </div>
 
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border-default">
                <th class="p-4 text-left font-medium text-text-secondary text-sm">Email</th>
                <th class="p-4 text-left font-medium text-text-secondary text-sm">Username</th>
                <th class="p-4 text-left font-medium text-text-secondary text-sm">Role</th>
                <th class="p-4 text-left font-medium text-text-secondary text-sm">Created</th>
                <th class="p-4 text-left font-medium text-text-secondary text-sm">Status</th>
                <th class="p-4 text-center font-medium text-text-secondary text-sm">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="user in paginatedUsers"
                :key="user.id"
                class="border-b border-border-subtle hover:bg-bg-tertiary transition-colors"
              >
                <td class="p-4">
                  <div class="text-text-primary">{{ user.email }}</div>
                </td>
                <td class="p-4">
                  <div class="text-text-muted">{{ user.username || '-' }}</div>
                </td>
                <td class="p-4">
                  <span
                    :class="[
                      'text-sm font-medium px-2 py-1 rounded-md',
                      user.role === 'admin'
                        ? 'bg-secondary-500/10 text-secondary-500'
                        : 'bg-success-500/10 text-success-500'
                    ]"
                  >{{ user.role }}</span>
                </td>
                <td class="p-4">
                  <div class="text-text-tertiary text-sm">
                    {{ formatDate(user.created_at) }}
                  </div>
                </td>
                <td class="p-4">
                  <span
                    :class="[
                      'text-sm font-medium px-2 py-1 rounded-md',
                      user.status === 'active'
                        ? 'bg-success-500/10 text-success-500'
                        : user.status === 'pending'
                          ? 'bg-warning-500/10 text-warning-500'
                          : 'bg-error-500/10 text-error-500'
                    ]"
                  >{{ statusLabel(user.status) }}</span>
                </td>
                <td class="p-4 text-center">
                  <div class="flex items-center justify-center gap-2">
                    <button
                      v-if="user.status === 'pending'"
                      @click="approveUser(user)"
                      class="p-2 text-success-500 hover:text-success-600 hover:bg-success-500/10 rounded-md transition-colors"
                      :disabled="loading"
                    >
                      Approve
                    </button>
                    <button
                      @click="confirmDeleteUser(user)"
                      class="p-2 text-error-500 hover:text-error-600 hover:bg-error-500/10 rounded-md transition-colors"
                      :disabled="loading"
                    >
                      Delete
                    </button>
                    <button
                      @click="toggleRoleMenu(user)"
                      class="p-2 text-primary-500 hover:text-primary-600 hover:bg-primary-500/10 rounded-md transition-colors"
                      :disabled="loading"
                    >
                      Change Role
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="pagination.totalPages > 1" class="flex items-center justify-between mt-6">
          <Button
            variant="outline"
            @click="changePage(currentPage - 1)"
            :disabled="currentPage === 1"
          >
            Previous
          </Button>
          <span class="text-text-muted text-sm">
            Page {{ currentPage }} of {{ pagination.totalPages }}
          </span>
          <Button
            variant="outline"
            @click="changePage(currentPage + 1)"
            :disabled="currentPage === pagination.totalPages"
          >
            Next
          </Button>
        </div>
      </div>
    </div>
 
    <!-- Create User Modal -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
    >
      <div class="bg-bg-secondary border border-border-default rounded-lg p-6 max-w-md w-full shadow-xl">
        <h2 class="font-display text-xl text-text-primary mb-4">Add User</h2>
        <div class="space-y-4">
          <Input
            id="new-email"
            v-model="newUser.email"
            type="email"
            label="Email"
            placeholder="user@example.com"
          />
          <Input
            id="new-username"
            v-model="newUser.username"
            type="text"
            label="Username"
            placeholder="username"
          />
          <Input
            id="new-password"
            v-model="newUser.password"
            type="password"
            label="Temporary Password"
            placeholder="Minimum 8 characters"
          />
          <label class="block text-sm font-medium text-text-secondary">
            Role
            <select
              v-model="newUser.role"
              class="mt-2 w-full bg-bg-tertiary border border-border-default rounded-md px-4 py-3 text-text-primary focus:outline-none focus:border-primary-500"
            >
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </label>
        </div>
        <div class="flex gap-3 mt-6">
          <Button variant="primary" class="flex-1" @click="createUser" :disabled="loading">
            Create user
          </Button>
          <Button variant="ghost" class="flex-1" @click="closeCreateModal">
            Cancel
          </Button>
        </div>
      </div>
    </div>

    <!-- Role Change Modal -->
    <div
      v-if="showRoleModal"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
    >
      <div class="bg-bg-secondary border border-border-default rounded-lg p-6 max-w-sm w-full shadow-xl">
        <h2 class="font-display text-xl text-text-primary mb-4">Change User Role</h2>
        <p class="text-text-muted mb-6 text-sm">
          {{ selectedUser?.email }}: {{ selectedUser?.username || selectedUser?.email }}
        </p>
        <div class="space-y-3">
          <Button
            variant="primary"
            class="w-full"
            @click="changeUserRole('admin')"
          >
            Make Admin
          </Button>
          <Button
            variant="secondary"
            class="w-full"
            @click="changeUserRole('user')"
          >
            Make User
          </Button>
          <Button
            variant="ghost"
            class="w-full"
            @click="showRoleModal = false"
          >
            Cancel
          </Button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div
      v-if="showDeleteModal"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
    >
      <div class="bg-bg-secondary border border-error-500/50 rounded-lg p-6 max-w-sm w-full shadow-xl">
        <h2 class="font-display text-xl text-error-500 mb-4">Confirm Delete</h2>
        <p class="text-text-muted mb-4">
          Are you sure you want to delete this user? This action cannot be undone.
        </p>
        <div class="space-y-3">
          <p class="text-text-muted text-sm mb-2">
            <span class="font-mono">{{ userToDelete?.email }}</span>
          </p>
          <p class="text-text-muted text-xs mb-2">
            {{ userToDelete?.username || userToDelete?.email }}
          </p>
          <div class="flex gap-3">
            <Button
              variant="destructive"
              class="flex-1"
              @click="executeDelete"
              :disabled="loading"
            >
              Yes, Delete
            </Button>
            <Button
              variant="ghost"
              class="flex-1"
              @click="showDeleteModal = false"
            >
              Cancel
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
 
<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import AppHeader from '@/components/layout/AppHeader.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import { adminApi } from '@/api/admin'
import type { User } from '@/types/user'
 
const router = useRouter()
const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const filters = ref({
  role: '',
  sort: 'created_at',
})
 
const users = ref<User[]>([])
const pagination = ref({
  total: 0,
  totalPages: 0,
})
 
const showRoleModal = ref(false)
const showDeleteModal = ref(false)
const showCreateModal = ref(false)
const selectedUser = ref<User | null>(null)
const userToDelete = ref<User | null>(null)
const newUser = ref({
  email: '',
  username: '',
  password: '',
  role: 'user' as 'admin' | 'user',
})
 
const paginatedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return users.value.slice(start, end)
})
 
const fetchUsers = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      limit: pageSize.value,
    }
    
    if (searchQuery.value) {
      params.search = searchQuery.value
    }
    
    if (filters.value.role) {
      params.role = filters.value.role
    }
    
    if (filters.value.sort) {
      params.sort = filters.value.sort
    }
 
    const response = await adminApi.listUsers(params)
    
    users.value = response.users
    pagination.value = {
      total: response.pagination.total,
      totalPages: response.pagination.pages,
    }
  } catch (err: any) {
    toast.error('Failed to fetch users')
    console.error('Failed to fetch users:', err)
  } finally {
    loading.value = false
  }
}
 
const changePage = (page: number) => {
  currentPage.value = page
  fetchUsers()
}
 
const toggleRoleMenu = (user: User) => {
  selectedUser.value = user
  showRoleModal.value = true
}
 
const confirmDeleteUser = (user: User) => {
  userToDelete.value = user
  showDeleteModal.value = true
}
 
const changeUserRole = async (role: 'admin' | 'user') => {
  if (!selectedUser.value) return
  
  loading.value = true
  try {
    await adminApi.changeUserRole(selectedUser.value.id, role)
    toast.success(`User role changed to ${role}`)
    
    const userIndex = users.value.findIndex(u => u.id === selectedUser.value.id)
    if (userIndex !== -1) {
      users.value[userIndex].role = role
    }
    
    showRoleModal.value = false
    fetchUsers()
  } catch (err: any) {
    toast.error('Failed to change user role')
    console.error('Failed to change user role:', err)
  } finally {
    loading.value = false
  }
}

const approveUser = async (user: User) => {
  loading.value = true
  try {
    await adminApi.approveUser(user.id)
    toast.success('User approved')
    const userIndex = users.value.findIndex(u => u.id === user.id)
    if (userIndex !== -1) {
      users.value[userIndex] = { ...users.value[userIndex], status: 'active' }
    }
  } catch (err: any) {
    toast.error('Failed to approve user')
    console.error('Failed to approve user:', err)
  } finally {
    loading.value = false
  }
}
 
const executeDelete = async () => {
  if (!userToDelete.value) return
  
  loading.value = true
  try {
    await adminApi.deleteUser(userToDelete.value.id)
    toast.success('User deleted successfully')
    
    users.value = users.value.filter(u => u.id !== userToDelete.value.id)
    pagination.value.total--
    
    showDeleteModal.value = false
    userToDelete.value = null
    
    // Adjust pagination if current page becomes empty
    if (paginatedUsers.value.length === 0 && currentPage.value > 1) {
      currentPage.value--
    }
  } catch (err: any) {
    toast.error('Failed to delete user')
    console.error('Failed to delete user:', err)
  } finally {
    loading.value = false
  }
}

const createUser = async () => {
  if (!newUser.value.email || !newUser.value.username || !newUser.value.password) {
    toast.error('Email, username, and password are required')
    return
  }

  loading.value = true
  try {
    await adminApi.createUser({
      email: newUser.value.email,
      username: newUser.value.username,
      password: newUser.value.password,
      role: newUser.value.role,
    })
    toast.success('User created')
    closeCreateModal()
    fetchUsers()
  } catch (err: any) {
    toast.error('Failed to create user')
    console.error('Failed to create user:', err)
  } finally {
    loading.value = false
  }
}

const closeCreateModal = () => {
  showCreateModal.value = false
  newUser.value = {
    email: '',
    username: '',
    password: '',
    role: 'user',
  }
}

const statusLabel = (status: User['status']) => {
  if (status === 'active') return 'Active'
  if (status === 'pending') return 'Pending'
  return 'Disabled'
}
 
const formatDate = (date: string) => {
  const d = new Date(date)
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}
 
watch([currentPage, filters], () => {
  fetchUsers()
})
 
onMounted(() => {
  fetchUsers()
})
</script>
 
<style scoped>
table {
  border-collapse: collapse;
}
 
th {
  white-space: nowrap;
}
</style>
