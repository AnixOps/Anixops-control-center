import { createI18n } from 'vue-i18n'
import enUS from './locales/en-US'
import zhCN from './locales/zh-CN'
import zhTW from './locales/zh-TW'
import jaJP from './locales/ja-JP'
import arSA from './locales/ar-SA'

// Supported locale codes
const supportedLocales = ['en-US', 'zh-CN', 'zh-TW', 'ja-JP', 'ar-SA']

// Get saved language or detect from browser
function getDefaultLocale(): string {
  const saved = localStorage.getItem('locale')
  if (saved && supportedLocales.includes(saved)) {
    return saved
  }

  // Detect from browser
  const browserLang = navigator.language || (navigator as any).userLanguage

  // Map browser languages to supported locales
  if (browserLang.startsWith('ar')) {
    return 'ar-SA'
  }
  if (browserLang.startsWith('zh-TW') || browserLang.startsWith('zh-HK') || browserLang.startsWith('zh-Hant')) {
    return 'zh-TW'
  }
  if (browserLang.startsWith('zh')) {
    return 'zh-CN'
  }
  if (browserLang.startsWith('ja')) {
    return 'ja-JP'
  }
  return 'en-US'
}

const i18n = createI18n({
  legacy: false,
  locale: getDefaultLocale(),
  fallbackLocale: 'en-US',
  messages: {
    'en-US': enUS,
    'zh-CN': zhCN,
    'zh-TW': zhTW,
    'ja-JP': jaJP,
    'ar-SA': arSA
  },
  datetimeFormats: {
    'en-US': {
      short: {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      },
      relative: {
        // Custom format for relative time
      }
    },
    'zh-CN': {
      short: {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }
    },
    'zh-TW': {
      short: {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }
    },
    'ja-JP': {
      short: {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }
    },
    'ar-SA': {
      short: {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      },
      long: {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      }
    }
  },
  numberFormats: {
    'en-US': {
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 0
      },
      bytes: {
        style: 'unit',
        unit: 'byte',
        unitDisplay: 'short'
      }
    },
    'zh-CN': {
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 0
      }
    },
    'zh-TW': {
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 0
      }
    },
    'ja-JP': {
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 0
      }
    },
    'ar-SA': {
      decimal: {
        style: 'decimal',
        minimumFractionDigits: 0,
        maximumFractionDigits: 2
      },
      percent: {
        style: 'percent',
        minimumFractionDigits: 0
      }
    }
  }
})

export default i18n

// Helper to change locale
export function setLocale(locale: string): void {
  i18n.global.locale.value = locale
  localStorage.setItem('locale', locale)
  document.documentElement.setAttribute('lang', locale)
  applyDirection(locale)
}

// Get current locale
export function getLocale(): string {
  return i18n.global.locale.value
}

// Available locales
export const availableLocales = [
  { code: 'en-US', name: 'English', nativeName: 'English' },
  { code: 'zh-CN', name: 'Chinese (Simplified)', nativeName: '简体中文' },
  { code: 'zh-TW', name: 'Chinese (Traditional)', nativeName: '繁體中文' },
  { code: 'ja-JP', name: 'Japanese', nativeName: '日本語' },
  { code: 'ar-SA', name: 'Arabic', nativeName: 'العربية' }
]

// RTL locale detection
export function isRTL(locale: string): boolean {
  const rtlLocales = ['ar', 'he', 'fa', 'ur']
  return rtlLocales.some(rtl => locale.startsWith(rtl))
}

// Apply RTL direction
export function applyDirection(locale: string): void {
  const dir = isRTL(locale) ? 'rtl' : 'ltr'
  document.documentElement.setAttribute('dir', dir)
}