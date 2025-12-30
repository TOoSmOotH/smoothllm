<template>
  <div class="min-h-screen bg-bg-primary flex flex-col">
    <AppHeader />
    <div class="flex-1 flex items-center justify-center bg-gradient-to-br from-bg-primary via-bg-primary to-bg-tertiary px-6 py-8 sm:px-8">
      <div class="relative z-10 w-full max-w-md">
      <!-- Register Card -->
        <div class="bg-bg-secondary border border-border-subtle rounded-2xl p-6 shadow-xl backdrop-blur-sm">
        <!-- Header -->
          <div class="text-center mb-8">
          <h1 class="text-4xl font-display font-bold text-text-primary mb-2">
            Register
          </h1>
          <p class="text-text-muted text-sm">
            Create your account
          </p>
          </div>

        <!-- Error Message -->
          <div v-if="authStore.error" class="mb-6 p-4 border border-error-500/50 bg-error-500/10 rounded-md">
          <p class="text-error-500 text-sm text-center">
            {{ authStore.error }}
          </p>
          </div>

        <!-- Register Form -->
          <form @submit.prevent="handleRegister" class="space-y-5">
          <!-- Username Input -->
            <Input
            id="username"
            v-model="form.username"
            type="text"
            label="Username"
            placeholder="username"
            :error="errors.username"
            helper-text="3-50 characters, letters, numbers, and underscores"
            required
            @enter="handleRegister"
          />

          <!-- Email Input -->
            <Input
            id="email"
            v-model="form.email"
            type="email"
            label="Email Address"
            placeholder="user@example.com"
            :error="errors.email"
            required
            @enter="handleRegister"
          />

          <!-- Password Input -->
            <Input
            id="password"
            v-model="form.password"
            type="password"
            label="Password"
            placeholder="Enter a secure password"
            :error="errors.password"
            helper-text="Minimum 8 characters"
            required
            @enter="handleRegister"
          />

          <!-- Confirm Password Input -->
            <Input
            id="confirmPassword"
            v-model="form.confirmPassword"
            type="password"
            label="Confirm Password"
            placeholder="Re-enter your password"
            :error="errors.confirmPassword"
            required
            @enter="handleRegister"
          />

          <!-- Terms and Conditions -->
            <div class="space-y-2">
            <label class="flex items-start space-x-2 cursor-pointer min-h-[44px]">
              <input
                type="checkbox"
                v-model="form.agreeToTerms"
                class="w-5 h-5 accent-primary-500 mt-0.5 flex-shrink-0 rounded border-border-default focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-bg-secondary"
              />
              <span class="text-text-muted text-sm leading-relaxed">
                I agree to
                <button type="button" class="text-primary-500 hover:text-primary-600 font-medium transition-colors duration-200">
                  Terms of Service
                </button>
                and
                <button type="button" class="text-primary-500 hover:text-primary-600 font-medium transition-colors duration-200">
                  Privacy Policy
                </button>
              </span>
            </label>
            <p v-if="errors.agreeToTerms" class="text-error-500 text-xs">
              {{ errors.agreeToTerms }}
            </p>
            </div>

          <!-- Submit Button -->
            <Button
            type="submit"
            variant="primary"
            size="lg"
            :loading="authStore.loading"
            :disabled="!form.agreeToTerms"
            class="w-full"
            @click="handleRegister"
          >
            <span v-if="!authStore.loading">Create Account</span>
            <span v-else>Creating account...</span>
            </Button>
          </form>

        <!-- Divider -->
          <div class="relative my-6">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-border-default"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-bg-secondary text-text-muted font-medium">
              OR
            </span>
          </div>
          </div>

        <!-- Social Login (Placeholder) -->
          <div class="space-y-3 flex flex-col items-center">
          <Button
            variant="outline"
            size="lg"
            class="w-full"
            @click="handleSocialLogin('github')"
          >
            <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            Continue with GitHub
          </Button>
          
          <Button
            variant="outline"
            size="lg"
            class="w-full"
            @click="handleSocialLogin('google')"
          >
            <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Continue with Google
          </Button>
          </div>

        <!-- Login Link -->
          <div class="text-center mt-6">
          <p class="text-text-muted text-sm">
            Already have an account?
            <router-link
              to="/login"
              class="text-primary-500 hover:text-primary-600 font-semibold transition-colors duration-200 ml-1"
            >
              Sign in
            </router-link>
          </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/layout/AppHeader.vue'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import type { RegisterRequest } from '@/types/auth'

const router = useRouter()
const authStore = useAuthStore()

// Form data
const form = reactive<RegisterRequest & { confirmPassword: string; agreeToTerms: boolean }>({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeToTerms: false,
})

// Form errors
const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeToTerms: '',
})

// Validate form
const validateForm = (): boolean => {
  // Clear previous errors
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = ''
  })

  let isValid = true

  // Username validation
  if (!form.username) {
    errors.username = 'Username is required'
    isValid = false
  } else if (form.username.length < 3) {
    errors.username = 'Username must be at least 3 characters'
    isValid = false
  } else if (form.username.length > 50) {
    errors.username = 'Username must be less than 50 characters'
    isValid = false
  } else if (!/^[a-zA-Z0-9_]+$/.test(form.username)) {
    errors.username = 'Username can only contain letters, numbers, and underscores'
    isValid = false
  }

  // Email validation
  if (!form.email) {
    errors.email = 'Email is required'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Invalid email format'
    isValid = false
  }

  // Password validation
  if (!form.password) {
    errors.password = 'Password is required'
    isValid = false
  } else if (form.password.length < 8) {
    errors.password = 'Password must be at least 8 characters'
    isValid = false
  }

  // Confirm password validation
  if (!form.confirmPassword) {
    errors.confirmPassword = 'Please confirm your password'
    isValid = false
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match'
    isValid = false
  }

  // Terms validation
  if (!form.agreeToTerms) {
    errors.agreeToTerms = 'You must agree to terms and conditions'
    isValid = false
  }

  return isValid
}

// Handle registration
const handleRegister = async () => {
  if (!validateForm()) return
  
  try {
    const response = await authStore.register({
      username: form.username,
      email: form.email,
      password: form.password,
    })

    if (response.approved) {
      router.push('/dashboard')
    } else {
      toast.success(response.message || 'Account created. Awaiting approval.')
      router.push('/login')
    }
  } catch (error) {
    // Error is handled by store
    console.error('Registration failed:', error)
  }
}

// Handle social login (placeholder)
const handleSocialLogin = (provider: string) => {
  console.log(`Social login with ${provider} - not implemented yet`)
}

onMounted(() => {
  // Clear any existing auth errors
  authStore.error = null
})
</script>

<style scoped>
</style>
