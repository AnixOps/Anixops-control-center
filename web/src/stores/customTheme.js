import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export interface CustomTheme {
  id: string
  name: string
  colors: ThemeColors
  typography: ThemeTypography
  spacing: ThemeSpacing
  borderRadius: ThemeBorderRadius
  shadows: ThemeShadows
  isCustom: boolean
  createdAt: string
  updatedAt: string
}

export interface ThemeColors {
  primary: string
  primaryHover: string
  primaryActive: string
  secondary: string
  secondaryHover: string
  accent: string
  background: string
  backgroundSecondary: string
  surface: string
  surfaceHover: string
  text: string
  textSecondary: string
  textMuted: string
  border: string
  borderLight: string
  success: string
  warning: string
  error: string
  info: string
}

export interface ThemeTypography {
  fontFamily: string
  fontSize: {
    xs: string
    sm: string
    base: string
    lg: string
    xl: string
    '2xl': string
    '3xl': string
  }
  fontWeight: {
    normal: number
    medium: number
    semibold: number
    bold: number
  }
  lineHeight: {
    tight: number
    normal: number
    relaxed: number
  }
}

export interface ThemeSpacing {
  xs: string
  sm: string
  md: string
  lg: string
  xl: string
  '2xl': string
}

export interface ThemeBorderRadius {
  none: string
  sm: string
  md: string
  lg: string
  xl: string
  full: string
}

export interface ThemeShadows {
  sm: string
  md: string
  lg: string
  xl: string
}

// Preset themes
const presetThemes: CustomTheme[] = [
  {
    id: 'default-dark',
    name: 'Default Dark',
    isCustom: false,
    colors: {
      primary: '#3B82F6',
      primaryHover: '#2563EB',
      primaryActive: '#1D4ED8',
      secondary: '#6366F1',
      secondaryHover: '#4F46E5',
      accent: '#8B5CF6',
      background: '#0F172A',
      backgroundSecondary: '#1E293B',
      surface: '#334155',
      surfaceHover: '#475569',
      text: '#F8FAFC',
      textSecondary: '#CBD5E1',
      textMuted: '#94A3B8',
      border: '#475569',
      borderLight: '#334155',
      success: '#22C55E',
      warning: '#F59E0B',
      error: '#EF4444',
      info: '#0EA5E9'
    },
    typography: {
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: { xs: '0.75rem', sm: '0.875rem', base: '1rem', lg: '1.125rem', xl: '1.25rem', '2xl': '1.5rem', '3xl': '1.875rem' },
      fontWeight: { normal: 400, medium: 500, semibold: 600, bold: 700 },
      lineHeight: { tight: 1.25, normal: 1.5, relaxed: 1.75 }
    },
    spacing: { xs: '0.25rem', sm: '0.5rem', md: '1rem', lg: '1.5rem', xl: '2rem', '2xl': '3rem' },
    borderRadius: { none: '0', sm: '0.125rem', md: '0.375rem', lg: '0.5rem', xl: '0.75rem', full: '9999px' },
    shadows: {
      sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
      md: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
      lg: '0 10px 15px -3px rgb(0 0 0 / 0.1)',
      xl: '0 20px 25px -5px rgb(0 0 0 / 0.1)'
    },
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  {
    id: 'ocean',
    name: 'Ocean Blue',
    isCustom: false,
    colors: {
      primary: '#0EA5E9',
      primaryHover: '#0284C7',
      primaryActive: '#0369A1',
      secondary: '#06B6D4',
      secondaryHover: '#0891B2',
      accent: '#14B8A6',
      background: '#0C4A6E',
      backgroundSecondary: '#075985',
      surface: '#0369A1',
      surfaceHover: '#0284C7',
      text: '#F0F9FF',
      textSecondary: '#BAE6FD',
      textMuted: '#7DD3FC',
      border: '#0284C7',
      borderLight: '#0369A1',
      success: '#10B981',
      warning: '#F59E0B',
      error: '#EF4444',
      info: '#06B6D4'
    },
    typography: {
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: { xs: '0.75rem', sm: '0.875rem', base: '1rem', lg: '1.125rem', xl: '1.25rem', '2xl': '1.5rem', '3xl': '1.875rem' },
      fontWeight: { normal: 400, medium: 500, semibold: 600, bold: 700 },
      lineHeight: { tight: 1.25, normal: 1.5, relaxed: 1.75 }
    },
    spacing: { xs: '0.25rem', sm: '0.5rem', md: '1rem', lg: '1.5rem', xl: '2rem', '2xl': '3rem' },
    borderRadius: { none: '0', sm: '0.125rem', md: '0.375rem', lg: '0.5rem', xl: '0.75rem', full: '9999px' },
    shadows: {
      sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
      md: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
      lg: '0 10px 15px -3px rgb(0 0 0 / 0.1)',
      xl: '0 20px 25px -5px rgb(0 0 0 / 0.1)'
    },
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  {
    id: 'forest',
    name: 'Forest Green',
    isCustom: false,
    colors: {
      primary: '#22C55E',
      primaryHover: '#16A34A',
      primaryActive: '#15803D',
      secondary: '#84CC16',
      secondaryHover: '#65A30D',
      accent: '#10B981',
      background: '#14532D',
      backgroundSecondary: '#166534',
      surface: '#15803D',
      surfaceHover: '#16A34A',
      text: '#F0FDF4',
      textSecondary: '#BBF7D0',
      textMuted: '#86EFAC',
      border: '#16A34A',
      borderLight: '#15803D',
      success: '#22C55E',
      warning: '#EAB308',
      error: '#DC2626',
      info: '#14B8A6'
    },
    typography: {
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: { xs: '0.75rem', sm: '0.875rem', base: '1rem', lg: '1.125rem', xl: '1.25rem', '2xl': '1.5rem', '3xl': '1.875rem' },
      fontWeight: { normal: 400, medium: 500, semibold: 600, bold: 700 },
      lineHeight: { tight: 1.25, normal: 1.5, relaxed: 1.75 }
    },
    spacing: { xs: '0.25rem', sm: '0.5rem', md: '1rem', lg: '1.5rem', xl: '2rem', '2xl': '3rem' },
    borderRadius: { none: '0', sm: '0.125rem', md: '0.375rem', lg: '0.5rem', xl: '0.75rem', full: '9999px' },
    shadows: {
      sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
      md: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
      lg: '0 10px 15px -3px rgb(0 0 0 / 0.1)',
      xl: '0 20px 25px -5px rgb(0 0 0 / 0.1)'
    },
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  {
    id: 'sunset',
    name: 'Sunset Orange',
    isCustom: false,
    colors: {
      primary: '#F97316',
      primaryHover: '#EA580C',
      primaryActive: '#C2410C',
      secondary: '#FB923C',
      secondaryHover: '#F97316',
      accent: '#FBBF24',
      background: '#431407',
      backgroundSecondary: '#7C2D12',
      surface: '#9A3412',
      surfaceHover: '#C2410C',
      text: '#FFF7ED',
      textSecondary: '#FED7AA',
      textMuted: '#FDBA74',
      border: '#EA580C',
      borderLight: '#9A3412',
      success: '#22C55E',
      warning: '#FBBF24',
      error: '#DC2626',
      info: '#0EA5E9'
    },
    typography: {
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: { xs: '0.75rem', sm: '0.875rem', base: '1rem', lg: '1.125rem', xl: '1.25rem', '2xl': '1.5rem', '3xl': '1.875rem' },
      fontWeight: { normal: 400, medium: 500, semibold: 600, bold: 700 },
      lineHeight: { tight: 1.25, normal: 1.5, relaxed: 1.75 }
    },
    spacing: { xs: '0.25rem', sm: '0.5rem', md: '1rem', lg: '1.5rem', xl: '2rem', '2xl': '3rem' },
    borderRadius: { none: '0', sm: '0.125rem', md: '0.375rem', lg: '0.5rem', xl: '0.75rem', full: '9999px' },
    shadows: {
      sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
      md: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
      lg: '0 10px 15px -3px rgb(0 0 0 / 0.1)',
      xl: '0 20px 25px -5px rgb(0 0 0 / 0.1)'
    },
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  },
  {
    id: 'purple-haze',
    name: 'Purple Haze',
    isCustom: false,
    colors: {
      primary: '#8B5CF6',
      primaryHover: '#7C3AED',
      primaryActive: '#6D28D9',
      secondary: '#A78BFA',
      secondaryHover: '#8B5CF6',
      accent: '#C084FC',
      background: '#1E1B4B',
      backgroundSecondary: '#312E81',
      surface: '#3730A3',
      surfaceHover: '#4338CA',
      text: '#F5F3FF',
      textSecondary: '#DDD6FE',
      textMuted: '#C4B5FD',
      border: '#6D28D9',
      borderLight: '#4C1D95',
      success: '#22C55E',
      warning: '#F59E0B',
      error: '#EF4444',
      info: '#06B6D4'
    },
    typography: {
      fontFamily: 'Inter, system-ui, sans-serif',
      fontSize: { xs: '0.75rem', sm: '0.875rem', base: '1rem', lg: '1.125rem', xl: '1.25rem', '2xl': '1.5rem', '3xl': '1.875rem' },
      fontWeight: { normal: 400, medium: 500, semibold: 600, bold: 700 },
      lineHeight: { tight: 1.25, normal: 1.5, relaxed: 1.75 }
    },
    spacing: { xs: '0.25rem', sm: '0.5rem', md: '1rem', lg: '1.5rem', xl: '2rem', '2xl': '3rem' },
    borderRadius: { none: '0', sm: '0.125rem', md: '0.375rem', lg: '0.5rem', xl: '0.75rem', full: '9999px' },
    shadows: {
      sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
      md: '0 4px 6px -1px rgb(0 0 0 / 0.1)',
      lg: '0 10px 15px -3px rgb(0 0 0 / 0.1)',
      xl: '0 20px 25px -5px rgb(0 0 0 / 0.1)'
    },
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString()
  }
]

const STORAGE_KEY = 'custom-themes'
const ACTIVE_THEME_KEY = 'active-theme'

export const useCustomThemeStore = defineStore('customTheme', () => {
  const themes = ref<CustomTheme[]>([...presetThemes])
  const activeThemeId = ref('default-dark')
  const isEditing = ref(false)
  const editingTheme = ref<CustomTheme | null>(null)

  // Computed
  const activeTheme = computed(() => {
    return themes.value.find(t => t.id === activeThemeId.value) || themes.value[0]
  })

  const customThemes = computed(() => {
    return themes.value.filter(t => t.isCustom)
  })

  const presetThemesList = computed(() => {
    return themes.value.filter(t => !t.isCustom)
  })

  // Load themes from storage
  function loadThemes() {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      try {
        const customThemes = JSON.parse(stored)
        // Add custom themes to presets
        themes.value = [...presetThemes, ...customThemes]
      } catch (e) {
        console.error('Failed to load custom themes:', e)
      }
    }

    const activeId = localStorage.getItem(ACTIVE_THEME_KEY)
    if (activeId && themes.value.find(t => t.id === activeId)) {
      activeThemeId.value = activeId
    }
  }

  // Save custom themes
  function saveCustomThemes() {
    const custom = themes.value.filter(t => t.isCustom)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(custom))
  }

  // Set active theme
  function setActiveTheme(id: string) {
    const theme = themes.value.find(t => t.id === id)
    if (theme) {
      activeThemeId.value = id
      localStorage.setItem(ACTIVE_THEME_KEY, id)
      applyTheme(theme)
    }
  }

  // Apply theme to CSS variables
  function applyTheme(theme: CustomTheme) {
    const root = document.documentElement

    // Apply colors
    Object.entries(theme.colors).forEach(([key, value]) => {
      const cssVar = `--color-${key.replace(/([A-Z])/g, '-$1').toLowerCase()}`
      root.style.setProperty(cssVar, value)
    })

    // Apply typography
    root.style.setProperty('--font-family', theme.typography.fontFamily)
    Object.entries(theme.typography.fontSize).forEach(([key, value]) => {
      root.style.setProperty(`--font-size-${key}`, value)
    })

    // Apply spacing
    Object.entries(theme.spacing).forEach(([key, value]) => {
      root.style.setProperty(`--spacing-${key}`, value)
    })

    // Apply border radius
    Object.entries(theme.borderRadius).forEach(([key, value]) => {
      root.style.setProperty(`--radius-${key}`, value)
    })
  }

  // Create new custom theme
  function createTheme(name: string, baseThemeId?: string): CustomTheme {
    const base = baseThemeId
      ? themes.value.find(t => t.id === baseThemeId)
      : themes.value[0]

    const newTheme: CustomTheme = {
      id: `custom-${Date.now()}`,
      name,
      isCustom: true,
      colors: { ...base?.colors || {} } as ThemeColors,
      typography: { ...base?.typography || {} } as ThemeTypography,
      spacing: { ...base?.spacing || {} } as ThemeSpacing,
      borderRadius: { ...base?.borderRadius || {} } as ThemeBorderRadius,
      shadows: { ...base?.shadows || {} } as ThemeShadows,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }

    themes.value.push(newTheme)
    saveCustomThemes()
    return newTheme
  }

  // Update theme
  function updateTheme(id: string, updates: Partial<CustomTheme>) {
    const theme = themes.value.find(t => t.id === id)
    if (theme && theme.isCustom) {
      Object.assign(theme, updates, { updatedAt: new Date().toISOString() })
      saveCustomThemes()

      if (activeThemeId.value === id) {
        applyTheme(theme)
      }
    }
  }

  // Delete custom theme
  function deleteTheme(id: string) {
    const index = themes.value.findIndex(t => t.id === id)
    if (index !== -1 && themes.value[index].isCustom) {
      themes.value.splice(index, 1)
      saveCustomThemes()

      // Switch to default if deleted theme was active
      if (activeThemeId.value === id) {
        setActiveTheme('default-dark')
      }
    }
  }

  // Duplicate theme
  function duplicateTheme(id: string): CustomTheme | null {
    const original = themes.value.find(t => t.id === id)
    if (!original) return null

    const duplicate: CustomTheme = {
      ...JSON.parse(JSON.stringify(original)),
      id: `custom-${Date.now()}`,
      name: `${original.name} (Copy)`,
      isCustom: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }

    themes.value.push(duplicate)
    saveCustomThemes()
    return duplicate
  }

  // Export theme
  function exportTheme(id: string): string {
    const theme = themes.value.find(t => t.id === id)
    if (!theme) return ''
    return JSON.stringify(theme, null, 2)
  }

  // Import theme
  function importTheme(json: string): CustomTheme | null {
    try {
      const theme = JSON.parse(json) as CustomTheme
      theme.id = `custom-${Date.now()}`
      theme.isCustom = true
      theme.createdAt = new Date().toISOString()
      theme.updatedAt = new Date().toISOString()

      themes.value.push(theme)
      saveCustomThemes()
      return theme
    } catch (e) {
      console.error('Failed to import theme:', e)
      return null
    }
  }

  // Initialize
  loadThemes()
  applyTheme(activeTheme.value)

  // Watch for changes
  watch(activeThemeId, () => {
    applyTheme(activeTheme.value)
  })

  return {
    themes,
    activeThemeId,
    activeTheme,
    customThemes,
    presetThemesList,
    isEditing,
    editingTheme,
    loadThemes,
    saveCustomThemes,
    setActiveTheme,
    applyTheme,
    createTheme,
    updateTheme,
    deleteTheme,
    duplicateTheme,
    exportTheme,
    importTheme
  }
})