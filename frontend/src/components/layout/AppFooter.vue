<template>
  <footer class="bg-bg-secondary border-t border-border-subtle mt-auto">
    <div class="container-lg py-12">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- Brand Section -->
        <div class="md:col-span-1">
          <div class="flex items-center space-x-3 mb-4">
            <div class="w-8 h-8 bg-primary-500 rounded-lg flex items-center justify-center">
              <span class="text-white font-display font-bold text-lg">{{ brand.shortName }}</span>
            </div>
            <span class="font-display font-semibold text-xl text-text-primary">
              {{ brand.name }}
            </span>
          </div>
          <p class="text-sm text-text-secondary mb-4">
            {{ footer.tagline }}
          </p>
          <div class="flex space-x-4">
            <a
              v-for="social in socialLinks"
              :key="social.name"
              :href="social.url"
              target="_blank"
              rel="noopener noreferrer"
              class="text-text-tertiary hover:text-primary-500 transition-colors duration-200"
              :aria-label="social.name"
            >
              <span class="sr-only">{{ social.name }}</span>
              <component :is="social.icon" class="w-5 h-5" />
            </a>
          </div>
        </div>

        <!-- Product Links -->
        <div>
          <h4 class="font-semibold text-text-primary mb-4">Product</h4>
          <ul class="space-y-3">
            <li v-for="link in productLinks" :key="link.label">
              <router-link
                :to="link.path"
                class="text-sm text-text-secondary hover:text-primary-500 transition-colors duration-200"
              >
                {{ link.label }}
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Resources Links -->
        <div>
          <h4 class="font-semibold text-text-primary mb-4">Resources</h4>
          <ul class="space-y-3">
            <li v-for="link in resourceLinks" :key="link.label">
              <a
                :href="link.url"
                target="_blank"
                rel="noopener noreferrer"
                class="text-sm text-text-secondary hover:text-primary-500 transition-colors duration-200"
              >
                {{ link.label }}
              </a>
            </li>
          </ul>
        </div>

        <!-- Legal Links -->
        <div>
          <h4 class="font-semibold text-text-primary mb-4">Legal</h4>
          <ul class="space-y-3">
            <li v-for="link in legalLinks" :key="link.label">
              <router-link
                :to="link.path"
                class="text-sm text-text-secondary hover:text-primary-500 transition-colors duration-200"
              >
                {{ link.label }}
              </router-link>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="border-t border-border-subtle mt-12 pt-8">
        <div class="flex flex-col md:flex-row justify-between items-center space-y-4 md:space-y-0">
          <p class="text-sm text-text-muted">
            &copy; {{ currentYear }} {{ brand.name }}. All rights reserved.
          </p>
          <div class="flex items-center space-x-2 text-sm text-text-muted">
            <span>Built with</span>
            <template v-for="(tech, index) in techStack" :key="tech">
              <span class="text-primary-500 font-medium">{{ tech }}</span>
              <span v-if="index < techStack.length - 1">•</span>
            </template>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { appConfig } from '@/config/appConfig'

// Simple SVG icons as components
const GitHubIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" /></svg>'
}

const TwitterIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z" /></svg>'
}

const LinkIcon = {
  template: '<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M13.828 10.172a4 4 0 010 5.656l-1.414 1.414a4 4 0 01-5.656-5.656l1.414-1.414M10.172 13.828a4 4 0 010-5.656l1.414-1.414a4 4 0 015.656 5.656l-1.414 1.414" /></svg>'
}

const currentYear = computed(() => new Date().getFullYear())

const brand = appConfig.brand
const footer = appConfig.footer

const iconMap: Record<string, { template: string }> = {
  github: GitHubIcon,
  twitter: TwitterIcon,
  link: LinkIcon,
}

const socialLinks = footer.socialLinks.map((link) => ({
  ...link,
  icon: iconMap[link.icon ?? ''] ?? LinkIcon,
}))

const productLinks = footer.productLinks
const resourceLinks = footer.resourceLinks
const legalLinks = footer.legalLinks
const techStack = footer.techStack
</script>

<style scoped>
/* Component uses Tailwind classes - no custom CSS needed */
</style>
