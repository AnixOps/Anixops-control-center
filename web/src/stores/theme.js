import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type Theme = 'light' | 'dark' | 'system'
export type ThemeColor = 'blue' | 'purple' | 'green' | 'orange' | 'red'

export interface ThemeConfig {
  theme: Theme
  color: ThemeColor
  compactMode: boolean
  fontSize: 'small' | 'medium' | 'large'
}

const DEFAULT_CONFIG: ThemeConfig = {
  theme: 'dark',
  color: 'blue',
  compactMode: false,
  fontSize: 'medium'
}

// Color palettes for different theme colors
export const colorPalettes: Record<ThemeColor, { primary: string; secondary: string; accent: string }> = {
  blue: { primary: '#3B82F6', secondary: '#1D4ED8', accent: '#60A5FA' },
  purple: { primary: '#8B5CF6', secondary: '#6D28D9', accent: '#A78BFA' },
  green: { primary: '#22C55E', secondary: '#16A34A', accent: '#4ADE80' },
  orange: { primary: '#F97316', secondary: '#EA580C', accent: '#FB923C' },
  red: { primary: '#EF4444', secondary: '#DC2626', accent: '#F87171' }
}

export const useThemeStore = defineStore('theme', () => {
  // State
  const config = ref<ThemeConfig>(loadConfig())

  // Apply theme on load
  applyTheme(config.value)

  // Watch for changes and persist
  watch(config, (newConfig) => {
    saveConfig(newConfig)
    applyTheme(newConfig)
  }, { deep: true })

  // Getters
  const isDark = ref(true)
  const currentColor = ref(config.value.color)

  // Actions
  function setTheme(theme: Theme) {
    config.value.theme = theme
    updateDarkMode()
  }

  function setColor(color: ThemeColor) {
    config.value.color = color
    currentColor.value = color
    applyColor(color)
  }

  function toggleCompactMode() {
    config.value.compactMode = !config.value.compactMode
    applyCompactMode(config.value.compactMode)
  }

  function setFontSize(size: 'small' | 'medium' | 'large') {
    config.value.fontSize = size
    applyFontSize(size)
  }

  function toggleTheme() {
    if (config.value.theme === 'system') {
      setTheme(isSystemDark() ? 'light' : 'dark')
    } else {
      setTheme(config.value.theme === 'dark' ? 'light' : 'dark')
    }
  }

  function reset() {
    config.value = { ...DEFAULT_CONFIG }
    applyTheme(config.value)
  }

  // Helper functions
  function loadConfig(): ThemeConfig {
    try {
      const saved = localStorage.getItem('theme-config')
      if (saved) {
        return { ...DEFAULT_CONFIG, ...JSON.parse(saved) }
      }
    } catch {
      // Ignore errors
    }
    return { ...DEFAULT_CONFIG }
  }

  function saveConfig(config: ThemeConfig) {
    try {
      localStorage.setItem('theme-config', JSON.stringify(config))
    } catch {
      // Ignore errors
    }
  }

  function isSystemDark(): boolean {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  function updateDarkMode() {
    if (config.value.theme === 'system') {
      isDark.value = isSystemDark()
    } else {
      isDark.value = config.value.theme === 'dark'
    }

    if (isDark.value) {
      document.documentElement.classList.add('dark')
      document.documentElement.classList.remove('light')
    } else {
      document.documentElement.classList.add('light')
      document.documentElement.classList.remove('dark')
    }
  }

  function applyTheme(config: ThemeConfig) {
    updateDarkMode()
    applyColor(config.color)
    applyCompactMode(config.compactMode)
    applyFontSize(config.fontSize)
  }

  function applyColor(color: ThemeColor) {
    const palette = colorPalettes[color]
    document.documentElement.style.setProperty('--color-primary', palette.primary)
    document.documentElement.style.setProperty('--color-secondary', palette.secondary)
    document.documentElement.style.setProperty('--color-accent', palette.accent)
  }

  function applyCompactMode(compact: boolean) {
    if (compact) {
      document.documentElement.classList.add('compact')
    } else {
      document.documentElement.classList.remove('compact')
    }
  }

  function applyFontSize(size: 'small' | 'medium' | 'large') {
    const sizes = {
      small: '14px',
      medium: '16px',
      large: '18px'
    }
    document.documentElement.style.fontSize = sizes[size]
  }

  // Listen for system theme changes
  function setupSystemThemeListener() {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', () => {
      if (config.value.theme === 'system') {
        updateDarkMode()
      }
    })
  }

  // Initialize
  setupSystemThemeListener()

  return {
    config,
    isDark,
    currentColor,
    setTheme,
    setColor,
    toggleCompactMode,
    setFontSize,
    toggleTheme,
    reset
  }
})