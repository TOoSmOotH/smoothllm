import type { AppConfig } from '@/config/appConfig'

export const customAppConfig: Partial<AppConfig> = {
  brand: {
    name: 'Smooth LLM Proxy',
    shortName: 'SLP',
  },
  navigation: {
    header: [
      { path: '/', label: 'Home' },
      { path: '/dashboard', label: 'Dashboard' },
      { path: '/providers', label: 'Providers' },
      { path: '/keys', label: 'API Keys' },
      { path: '/usage', label: 'Usage' },
    ],
    sidebar: {
      main: [
        { path: '/', label: 'Home', icon: 'home' },
        { path: '/dashboard', label: 'Dashboard', icon: 'dashboard' },
        { path: '/providers', label: 'Providers', icon: 'settings' },
        { path: '/keys', label: 'API Keys', icon: 'lock' },
        { path: '/usage', label: 'Usage', icon: 'chart' },
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
    tagline: 'A LiteLLM-inspired API key management and LLM usage tracking proxy.',
    socialLinks: [
      { name: 'GitHub', url: 'https://github.com', icon: 'github' },
    ],
    productLinks: [
      { label: 'Providers', path: '/providers' },
      { label: 'API Keys', path: '/keys' },
      { label: 'Usage', path: '/usage' },
      { label: 'Documentation', path: '/docs' },
    ],
    resourceLinks: [
      { label: 'LiteLLM', url: 'https://docs.litellm.ai' },
      { label: 'OpenAI API', url: 'https://platform.openai.com/docs' },
      { label: 'Anthropic API', url: 'https://docs.anthropic.com' },
    ],
    legalLinks: [
      { label: 'Privacy Policy', path: '/privacy' },
      { label: 'Terms of Service', path: '/terms' },
    ],
    techStack: ['Vue 3', 'Go', 'Tailwind CSS', 'SQLite'],
  },
  home: {
    title: 'Smooth LLM Proxy',
    subtitle: 'Unified LLM Gateway',
    description:
      'Manage your LLM provider API keys, create proxy keys, and track usage across OpenAI, Anthropic, and local models with a single unified API endpoint.',
    ctas: [
      { label: 'Get Started', to: '/register', variant: 'primary' },
      { label: 'Sign In', to: '/login', variant: 'secondary' },
    ],
    features: [
      {
        title: 'Multi-Provider',
        description: 'Connect OpenAI, Anthropic, and local LLM endpoints',
        icon: 'sparkles',
        tone: 'primary',
      },
      {
        title: 'Proxy Keys',
        description: 'Create secure proxy API keys for your applications',
        icon: 'lock',
        tone: 'secondary',
      },
      {
        title: 'Usage Tracking',
        description: 'Monitor tokens, costs, and request statistics',
        icon: 'chart',
        tone: 'success',
      },
    ],
  },
}
