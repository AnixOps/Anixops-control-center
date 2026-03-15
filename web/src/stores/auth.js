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
      token.value = response.data.access_token
      user.value = {
        id: response.data.user_id || '1',
        email: email,
        role: response.data.role || 'admin'
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
      token.value = response.data.access_token
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