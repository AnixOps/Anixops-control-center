import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useThemeStore } from '@/stores/theme'

// Mock localStorage
const localStorageMock = {
  store: {} as Record<string, string>,
  getItem: vi.fn((key: string) => localStorageMock.store[key] || null),
  setItem: vi.fn((key: string, value: string) => {
    localStorageMock.store[key] = value
  }),
  removeItem: vi.fn((key: string) => {
    delete localStorageMock.store[key]
  }),
  clear: vi.fn(() => {
    localStorageMock.store = {}
  })
}

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  value: vi.fn((query: string) => ({
    matches: query === '(prefers-color-scheme: dark)',
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
})

describe('Theme Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorageMock.clear()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  it('has correct initial state', () => {
    const store = useThemeStore()

    expect(store.mode).toBe('system')
    expect(store.accentColor).toBe('blue')
    expect(store.compactMode).toBe(false)
    expect(store.fontSize).toBe('medium')
  })

  it('sets theme mode correctly', () => {
    const store = useThemeStore()

    store.setMode('dark')
    expect(store.mode).toBe('dark')

    store.setMode('light')
    expect(store.mode).toBe('light')

    store.setMode('system')
    expect(store.mode).toBe('system')
  })

  it('sets accent color correctly', () => {
    const store = useThemeStore()

    const colors = ['blue', 'purple', 'green', 'orange', 'pink']
    colors.forEach(color => {
      store.setAccentColor(color)
      expect(store.accentColor).toBe(color)
    })
  })

  it('sets compact mode correctly', () => {
    const store = useThemeStore()

    store.setCompactMode(true)
    expect(store.compactMode).toBe(true)

    store.setCompactMode(false)
    expect(store.compactMode).toBe(false)
  })

  it('sets font size correctly', () => {
    const store = useThemeStore()

    const sizes = ['small', 'medium', 'large']
    sizes.forEach(size => {
      store.setFontSize(size)
      expect(store.fontSize).toBe(size)
    })
  })

  it('computes effectiveTheme correctly', () => {
    const store = useThemeStore()

    // System mode - should match prefers-color-scheme
    store.setMode('system')
    expect(store.effectiveTheme).toBe('dark')

    // Dark mode
    store.setMode('dark')
    expect(store.effectiveTheme).toBe('dark')

    // Light mode
    store.setMode('light')
    expect(store.effectiveTheme).toBe('light')
  })

  it('computes isDark correctly', () => {
    const store = useThemeStore()

    store.setMode('dark')
    expect(store.isDark).toBe(true)

    store.setMode('light')
    expect(store.isDark).toBe(false)
  })

  it('toggles theme correctly', () => {
    const store = useThemeStore()

    store.setMode('light')
    store.toggleTheme()
    expect(store.mode).toBe('dark')

    store.toggleTheme()
    expect(store.mode).toBe('light')
  })

  it('persists settings to localStorage', () => {
    const store = useThemeStore()

    store.setMode('dark')
    store.setAccentColor('purple')
    store.setCompactMode(true)
    store.setFontSize('large')

    expect(localStorageMock.setItem).toHaveBeenCalled()
  })

  it('loads settings from localStorage', () => {
    localStorageMock.store['theme-settings'] = JSON.stringify({
      mode: 'dark',
      accentColor: 'green',
      compactMode: true,
      fontSize: 'large'
    })

    const store = useThemeStore()

    expect(store.mode).toBe('dark')
    expect(store.accentColor).toBe('green')
    expect(store.compactMode).toBe(true)
    expect(store.fontSize).toBe('large')
  })

  it('resets to defaults', () => {
    const store = useThemeStore()

    store.setMode('dark')
    store.setAccentColor('purple')
    store.setCompactMode(true)
    store.setFontSize('large')

    store.resetToDefaults()

    expect(store.mode).toBe('system')
    expect(store.accentColor).toBe('blue')
    expect(store.compactMode).toBe(false)
    expect(store.fontSize).toBe('medium')
  })
})