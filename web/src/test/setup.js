import { config } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { vi } from 'vitest'

// Setup global mocks
config.global.mocks = {
  $t: (key) => key,
  $i18n: {
    locale: 'en-US'
  }
}

// Auto-setup Pinia for each test
beforeEach(() => {
  setActivePinia(createPinia())
})

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
})