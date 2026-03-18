import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { wailsAPI, type AppConfig, type UserInfo } from '../api/wails'

export const useAppStore = defineStore('app', () => {
  // State
  const config = ref<AppConfig | null>(null)
  const user = ref<UserInfo | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const initialized = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const theme = computed(() => config.value?.theme || 'dark')
  const apiUrl = computed(() => config.value?.api_url || 'http://localhost:8080/api/v1')
  const language = computed(() => config.value?.language || 'en')

  // Actions
  async function initialize() {
    if (initialized.value) return

    loading.value = true
    try {
      // Load config
      config.value = await wailsAPI.getConfig()

      // Check authentication
      const authenticated = await wailsAPI.isAuthenticated()
      if (authenticated) {
        token.value = await wailsAPI.getAuthToken()
        user.value = await wailsAPI.getUser()
      }

      initialized.value = true
    } catch (e) {
      error.value = String(e)
    } finally {
      loading.value = false
    }
  }

  async function setTheme(newTheme: string) {
    try {
      await wailsAPI.setTheme(newTheme)
      if (config.value) {
        config.value.theme = newTheme
      }
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function setApiUrl(url: string) {
    try {
      await wailsAPI.setAPIUrl(url)
      if (config.value) {
        config.value.api_url = url
      }
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function setLanguage(lang: string) {
    try {
      await wailsAPI.setLanguage(lang)
      if (config.value) {
        config.value.language = lang
      }
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function login(email: string, password: string): Promise<boolean> {
    loading.value = true
    error.value = null

    try {
      // Call backend API for login
      const response = await fetch(`${apiUrl.value}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      })

      if (!response.ok) {
        throw new Error('Login failed')
      }

      const data = await response.json()
      token.value = data.access_token

      const userInfo: UserInfo = {
        id: data.user_id || '1',
        email: email,
        name: data.name || email.split('@')[0],
        role: data.role || 'admin'
      }
      user.value = userInfo

      // Store in Wails backend
      await wailsAPI.setAuthTokens(data.access_token, data.refresh_token || '')
      await wailsAPI.setUser(userInfo)

      return true
    } catch (e) {
      error.value = String(e)
      return false
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      await wailsAPI.logout()
      token.value = null
      user.value = null
    } catch (e) {
      error.value = String(e)
    }
  }

  async function checkForUpdates() {
    try {
      return await wailsAPI.checkForUpdates()
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function getNetworkInterfaces() {
    try {
      return await wailsAPI.getNetworkInterfaces()
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function ping(host: string) {
    try {
      return await wailsAPI.ping(host)
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function dnsLookup(domain: string) {
    try {
      return await wailsAPI.dnsLookup(domain)
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  // Window controls
  function minimizeWindow() {
    wailsAPI.minimize()
  }

  function maximizeWindow() {
    wailsAPI.toggleMaximize()
  }

  function closeWindow() {
    wailsAPI.close()
  }

  function openExternalURL(url: string) {
    wailsAPI.openURL(url)
  }

  return {
    // State
    config,
    user,
    token,
    loading,
    error,
    initialized,
    // Getters
    isAuthenticated,
    theme,
    apiUrl,
    language,
    // Actions
    initialize,
    setTheme,
    setApiUrl,
    setLanguage,
    login,
    logout,
    checkForUpdates,
    getNetworkInterfaces,
    ping,
    dnsLookup,
    // Window controls
    minimizeWindow,
    maximizeWindow,
    closeWindow,
    openExternalURL
  }
})