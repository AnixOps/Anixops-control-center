import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  getDefaultLocale,
  setLocale,
  getLocale,
  availableLocales,
  isRTL,
  applyDirection
} from '@/i18n'

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

// Mock navigator
Object.defineProperty(window, 'navigator', {
  value: {
    language: 'en-US',
    userLanguage: 'en-US'
  }
})

// Mock document.documentElement
const documentElementMock = {
  setAttribute: vi.fn(),
  getAttribute: vi.fn(),
  lang: 'en-US',
  dir: 'ltr'
}

Object.defineProperty(document, 'documentElement', {
  value: documentElementMock
})

describe('i18n', () => {
  beforeEach(() => {
    localStorageMock.clear()
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.clearAllMocks()
  })

  describe('availableLocales', () => {
    it('contains all expected locales', () => {
      const localeCodes = availableLocales.map(l => l.code)

      expect(localeCodes).toContain('en-US')
      expect(localeCodes).toContain('zh-CN')
      expect(localeCodes).toContain('zh-TW')
      expect(localeCodes).toContain('ja-JP')
      expect(localeCodes).toContain('ar-SA')
    })

    it('has correct native names for each locale', () => {
      const enUS = availableLocales.find(l => l.code === 'en-US')
      const zhCN = availableLocales.find(l => l.code === 'zh-CN')
      const zhTW = availableLocales.find(l => l.code === 'zh-TW')
      const jaJP = availableLocales.find(l => l.code === 'ja-JP')
      const arSA = availableLocales.find(l => l.code === 'ar-SA')

      expect(enUS?.nativeName).toBe('English')
      expect(zhCN?.nativeName).toBe('简体中文')
      expect(zhTW?.nativeName).toBe('繁體中文')
      expect(jaJP?.nativeName).toBe('日本語')
      expect(arSA?.nativeName).toBe('العربية')
    })
  })

  describe('isRTL', () => {
    it('returns true for RTL locales', () => {
      expect(isRTL('ar-SA')).toBe(true)
      expect(isRTL('ar')).toBe(true)
      expect(isRTL('he-IL')).toBe(true)
      expect(isRTL('fa-IR')).toBe(true)
    })

    it('returns false for LTR locales', () => {
      expect(isRTL('en-US')).toBe(false)
      expect(isRTL('zh-CN')).toBe(false)
      expect(isRTL('ja-JP')).toBe(false)
    })
  })

  describe('applyDirection', () => {
    it('sets RTL for Arabic locale', () => {
      applyDirection('ar-SA')
      expect(documentElementMock.setAttribute).toHaveBeenCalledWith('dir', 'rtl')
    })

    it('sets LTR for English locale', () => {
      applyDirection('en-US')
      expect(documentElementMock.setAttribute).toHaveBeenCalledWith('dir', 'ltr')
    })
  })

  describe('setLocale', () => {
    it('sets locale and persists to localStorage', () => {
      setLocale('zh-CN')

      expect(localStorageMock.setItem).toHaveBeenCalledWith('locale', 'zh-CN')
      expect(documentElementMock.setAttribute).toHaveBeenCalledWith('lang', 'zh-CN')
    })

    it('applies RTL direction for Arabic locale', () => {
      setLocale('ar-SA')

      expect(documentElementMock.setAttribute).toHaveBeenCalledWith('lang', 'ar-SA')
      expect(documentElementMock.setAttribute).toHaveBeenCalledWith('dir', 'rtl')
    })
  })
})