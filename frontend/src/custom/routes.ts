import type { RouteRecordRaw } from 'vue-router'

export const customRoutes: RouteRecordRaw[] = [
  {
    path: '/providers',
    name: 'providers',
    component: () => import('@/views/Providers.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/keys',
    name: 'api-keys',
    component: () => import('@/views/ApiKeys.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/usage',
    name: 'usage',
    component: () => import('@/views/Usage.vue'),
    meta: { requiresAuth: true },
  },
]
