import { defineStore } from 'pinia'
import { settingsApi } from '@/api/settings'

export type ThemeId =
  | 'warm-editorial'
  | 'warm-editorial-dark'
  | 'solarpunk'
  | 'gothic-punk'
  | 'cyberpunk'
  | 'cassette-futurism'
  | 'cassette-futurism-dark'
  | 'vaporware'

export const themeOptions: Array<{ id: ThemeId; label: string; description: string }> = [
  {
    id: 'warm-editorial',
    label: 'Warm Editorial',
    description: 'Soft paper tones, rich ink accents, and serif-forward typography.',
  },
  {
    id: 'warm-editorial-dark',
    label: 'Warm Editorial Dark',
    description: 'Ink-on-walnut tones with warm highlights and moody contrast.',
  },
  {
    id: 'solarpunk',
    label: 'Solarpunk',
    description: 'Sunlit greens with organic warmth and optimistic contrast.',
  },
  {
    id: 'gothic-punk',
    label: 'Gothic Punk',
    description: 'Deep shadows, crimson highlights, and cathedral-inspired drama.',
  },
  {
    id: 'cyberpunk',
    label: 'Cyberpunk',
    description: 'Neon cyan/magenta on a high-contrast nightscape.',
  },
  {
    id: 'cassette-futurism',
    label: 'Cassette Futurism',
    description: 'Analog warmth, punchy neons, and chrome-era optimism.',
  },
  {
    id: 'cassette-futurism-dark',
    label: 'Cassette Futurism Dark',
    description: 'After-hours synth glow with warm tape highlights.',
  },
  {
    id: 'vaporware',
    label: 'Vaporware',
    description: 'Pastel glow, soft gradients, and nostalgic digital haze.',
  },
]

export const useThemeStore = defineStore('theme', {
  state: () => ({
    theme: 'warm-editorial' as ThemeId,
    isInitialized: false,
  }),
  actions: {
    async initTheme() {
      if (this.isInitialized) return

      if (!themeOptions.some((option) => option.id === this.theme)) {
        this.theme = 'warm-editorial'
      }
      this.applyTheme(this.theme)

      try {
        const serverTheme = await settingsApi.getTheme()
        if (themeOptions.some((option) => option.id === serverTheme)) {
          this.theme = serverTheme as ThemeId
          this.applyTheme(this.theme)
        }
      } catch (err) {
        // Keep local theme if server settings are unavailable.
      } finally {
        this.isInitialized = true
      }
    },
    async setThemeForAll(theme: ThemeId) {
      if (!themeOptions.some((option) => option.id === theme)) return

      const previousTheme = this.theme
      this.theme = theme
      this.applyTheme(theme)

      try {
        await settingsApi.setTheme(theme)
      } catch (err) {
        this.theme = previousTheme
        this.applyTheme(previousTheme)
        throw err
      }
    },
    applyTheme(theme: ThemeId) {
      if (typeof document === 'undefined') return
      const root = document.documentElement
      root.setAttribute('data-theme', theme)
    },
  },
  persist: {
    key: 'smoothweb-theme',
    paths: ['theme'],
  },
})
