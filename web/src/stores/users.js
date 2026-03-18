import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { usersApi } from '@/api'

export const useUsersStore = defineStore('users', () => {
  // State
  const users = ref([])
  const currentUser = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const total = ref(0)
  const filters = ref({
    search: '',
    role: '',
    status: '',
    page: 1,
    limit: 20
  })

  // Getters
  const activeUsers = computed(() => users.value.filter(u => u.status === 'active'))
  const bannedUsers = computed(() => users.value.filter(u => u.status === 'banned'))
  const adminUsers = computed(() => users.value.filter(u => u.role === 'admin'))
  const userStats = computed(() => ({
    total: total.value,
    active: activeUsers.value.length,
    banned: bannedUsers.value.length,
    admins: adminUsers.value.length
  }))

  // Actions
  async function fetchUsers(params = {}) {
    loading.value = true
    error.value = null
    try {
      const response = await usersApi.list({ ...filters.value, ...params })
      users.value = response.data.data || response.data
      total.value = response.data.total || users.value.length
    } catch (e) {
      error.value = e.message || 'Failed to fetch users'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchUser(id) {
    loading.value = true
    error.value = null
    try {
      const response = await usersApi.get(id)
      currentUser.value = response.data
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch user'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function createUser(data) {
    loading.value = true
    error.value = null
    try {
      const response = await usersApi.create(data)
      users.value.unshift(response.data)
      total.value++
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to create user'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateUser(id, data) {
    loading.value = true
    error.value = null
    try {
      const response = await usersApi.update(id, data)
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        users.value[index] = response.data
      }
      if (currentUser.value?.id === id) {
        currentUser.value = response.data
      }
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to update user'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteUser(id) {
    loading.value = true
    error.value = null
    try {
      await usersApi.delete(id)
      users.value = users.value.filter(u => u.id !== id)
      total.value--
      if (currentUser.value?.id === id) {
        currentUser.value = null
      }
    } catch (e) {
      error.value = e.message || 'Failed to delete user'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function banUser(id) {
    try {
      const response = await usersApi.ban(id)
      updateUserStatus(id, 'banned')
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function unbanUser(id) {
    try {
      const response = await usersApi.unban(id)
      updateUserStatus(id, 'active')
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function resetPassword(id) {
    try {
      const response = await usersApi.resetPassword(id)
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function updateRole(id, role) {
    try {
      const response = await usersApi.updateRole(id, role)
      const index = users.value.findIndex(u => u.id === id)
      if (index !== -1) {
        users.value[index].role = role
      }
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function exportUsers(params = {}) {
    try {
      const response = await usersApi.export(params)
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `users-${new Date().toISOString().split('T')[0]}.csv`)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    } catch (e) {
      throw e
    }
  }

  function setFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function clearCurrentUser() {
    currentUser.value = null
  }

  function updateUserStatus(id, status) {
    const index = users.value.findIndex(u => u.id === id)
    if (index !== -1) {
      users.value[index].status = status
    }
    if (currentUser.value?.id === id) {
      currentUser.value.status = status
    }
  }

  return {
    // State
    users,
    currentUser,
    loading,
    error,
    total,
    filters,
    // Getters
    activeUsers,
    bannedUsers,
    adminUsers,
    userStats,
    // Actions
    fetchUsers,
    fetchUser,
    createUser,
    updateUser,
    deleteUser,
    banUser,
    unbanUser,
    resetPassword,
    updateRole,
    exportUsers,
    setFilters,
    clearCurrentUser,
    updateUserStatus
  }
})