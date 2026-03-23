import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || null)
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  async function login(email, password) {
    try {
      const response = await api.post('/auth/login', { email, password })

      // Handle Workers API response format: { success: true, data: { access_token, user } }
      const data = response.data.data || response.data
      token.value = data.access_token
      user.value = {
        id: data.user?.id || '1',
        email: data.user?.email || email,
        role: data.user?.role || 'admin'
      }

      localStorage.setItem('token', token.value)
      localStorage.setItem('user', JSON.stringify(user.value))

      return { success: true }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Login failed'
      }
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function refreshToken() {
    try {
      const response = await api.post('/auth/refresh')
      const data = response.data.data || response.data
      token.value = data.access_token
      localStorage.setItem('token', token.value)
      return true
    } catch {
      logout()
      return false
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    isAdmin,
    login,
    logout,
    refreshToken
  }
})