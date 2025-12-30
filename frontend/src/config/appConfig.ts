import { customAppConfig } from '@/custom/appConfig'

export type NavItem = {
  label: string
  path: string
}

export type SidebarNavItem = NavItem & {
  icon?: string
}

export type SocialLink = {
  name: string
  url: string
  icon?: string
}

export type FooterLink = {
  label: string
  path: string
}

export type ExternalLink = {
  label: string
  url: string
}

export type HomeCta = {
  label: string
  to: string
  variant?: 'primary' | 'secondary' | 'ghost'
}

export type HomeFeature = {
  title: string
  description: string
  icon: string
  tone?: 'primary' | 'secondary' | 'success'
}

export type AppConfig = {
  brand: {
    name: string
    shortName: string
  }
  navigation: {
    header: NavItem[]
    sidebar: {
      main: SidebarNavItem[]
      admin: SidebarNavItem[]
    }
  }
  footer: {
    tagline: string
    socialLinks: SocialLink[]
    productLinks: FooterLink[]
    resourceLinks: ExternalLink[]
    legalLinks: FooterLink[]
    techStack: string[]
  }
  home: {
    title: string
    subtitle: string
    description: string
    ctas: HomeCta[]
    features: HomeFeature[]
  }
}

const defaultConfig: AppConfig = {
  brand: {
    name: 'SmoothWeb',
    shortName: 'S',
  },
  navigation: {
    header: [
      { path: '/', label: 'Home' },
      { path: '/dashboard', label: 'Dashboard' },
      { path: '/profile', label: 'Profile' },
    ],
    sidebar: {
      main: [
        { path: '/', label: 'Home', icon: 'home' },
        { path: '/dashboard', label: 'Dashboard', icon: 'dashboard' },
        { path: '/profile', label: 'Profile', icon: 'profile' },
      ],
      admin: [
        { path: '/admin', label: 'Admin Dashboard', icon: 'admin' },
        { path: '/admin/users', label: 'Users', icon: 'users' },
        { path: '/admin/settings', label: 'Settings', icon: 'settings' },
      ],
    },
  },
  footer: {
    tagline: 'A modern, professional web application template with full-stack capabilities.',
    socialLinks: [
      { name: 'GitHub', url: 'https://github.com', icon: 'github' },
      { name: 'Twitter', url: 'https://twitter.com', icon: 'twitter' },
    ],
    productLinks: [
      { label: 'Features', path: '/features' },
      { label: 'Pricing', path: '/pricing' },
      { label: 'Documentation', path: '/docs' },
      { label: 'Changelog', path: '/changelog' },
    ],
    resourceLinks: [
      { label: 'Blog', url: '/blog' },
      { label: 'Community', url: '/community' },
      { label: 'Support', url: '/support' },
      { label: 'Status', url: '/status' },
    ],
    legalLinks: [
      { label: 'Privacy Policy', path: '/privacy' },
      { label: 'Terms of Service', path: '/terms' },
      { label: 'Cookie Policy', path: '/cookies' },
    ],
    techStack: ['Vue 3', 'Go', 'Tailwind CSS'],
  },
  home: {
    title: 'SmoothWeb',
    subtitle: 'Modern Web Application Template',
    description: 'A production-ready, full-stack web application featuring encrypted SQLite database, user management with RBAC, and a professional, modern design system.',
    ctas: [
      { label: 'Get Started', to: '/register', variant: 'primary' },
      { label: 'Sign In', to: '/login', variant: 'secondary' },
    ],
    features: [
      {
        title: 'Secure Auth',
        description: 'JWT authentication with encrypted SQLite database',
        icon: 'lock',
        tone: 'primary',
      },
      {
        title: 'Modern UI',
        description: 'Clean, professional design with accessibility in mind',
        icon: 'sparkles',
        tone: 'secondary',
      },
      {
        title: 'Modern Stack',
        description: 'Vue 3, Go, Tailwind CSS, and Docker ready',
        icon: 'bolt',
        tone: 'success',
      },
    ],
  },
}

const isPlainObject = (value: unknown): value is Record<string, unknown> => {
  if (!value || typeof value !== 'object') return false
  return Object.prototype.toString.call(value) === '[object Object]'
}

const mergeAppConfig = <T extends Record<string, unknown>>(base: T, override: Partial<T>): T => {
  const result: Record<string, unknown> = { ...base }
  for (const key of Object.keys(override)) {
    const overrideValue = override[key]
    const baseValue = base[key]

    if (Array.isArray(overrideValue)) {
      result[key] = overrideValue
      continue
    }

    if (isPlainObject(baseValue) && isPlainObject(overrideValue)) {
      result[key] = mergeAppConfig(baseValue, overrideValue)
      continue
    }

    if (overrideValue !== undefined) {
      result[key] = overrideValue
    }
  }
  return result as T
}

export const appConfig: AppConfig = mergeAppConfig(defaultConfig, customAppConfig)
