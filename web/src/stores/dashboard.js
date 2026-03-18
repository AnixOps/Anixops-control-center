import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { dashboardApi } from '@/api'

export const useDashboardStore = defineStore('dashboard', () => {
  // State
  const stats = ref({
    nodes: { total: 0, online: 0, offline: 0 },
    users: { total: 0, active: 0 },
    agents: { total: 0, online: 0 },
    traffic: { today: 0, month: 0 },
    plugins: { total: 0, active: 0 }
  })
  const activities = ref([])
  const alerts = ref([])
  const traffic = ref([])
  const loading = ref(false)
  const error = ref(null)
  const lastUpdate = ref(null)

  // Getters
  const hasAlerts = computed(() => alerts.value.length > 0)
  const criticalAlerts = computed(() => alerts.value.filter(a => a.level === 'critical'))
  const warningAlerts = computed(() => alerts.value.filter(a => a.level === 'warning'))

  // Actions
  async function fetchStats() {
    try {
      const response = await dashboardApi.stats()
      stats.value = response.data
      lastUpdate.value = new Date()
    } catch (e) {
      error.value = e.message || 'Failed to fetch stats'
      throw e
    }
  }

  async function fetchActivities(limit = 10) {
    try {
      const response = await dashboardApi.activities(limit)
      activities.value = response.data.data || response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch activities'
      throw e
    }
  }

  async function fetchAlerts() {
    try {
      const response = await dashboardApi.alerts()
      alerts.value = response.data.data || response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch alerts'
      throw e
    }
  }

  async function fetchTraffic(period = '24h') {
    try {
      const response = await dashboardApi.traffic(period)
      traffic.value = response.data.data || response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch traffic'
      throw e
    }
  }

  async function fetchAll() {
    loading.value = true
    error.value = null
    try {
      await Promise.all([
        fetchStats(),
        fetchActivities(),
        fetchAlerts(),
        fetchTraffic()
      ])
    } catch (e) {
      error.value = e.message || 'Failed to fetch dashboard data'
    } finally {
      loading.value = false
    }
  }

  function addActivity(activity) {
    activities.value.unshift(activity)
    if (activities.value.length > 50) {
      activities.value = activities.value.slice(0, 50)
    }
  }

  function addAlert(alert) {
    alerts.value.unshift(alert)
  }

  function removeAlert(id) {
    alerts.value = alerts.value.filter(a => a.id !== id)
  }

  function clearAlerts() {
    alerts.value = []
  }

  function updateStats(newStats) {
    stats.value = { ...stats.value, ...newStats }
    lastUpdate.value = new Date()
  }

  function updateTrafficPoint(point) {
    traffic.value.push(point)
    if (traffic.value.length > 100) {
      traffic.value = traffic.value.slice(-100)
    }
  }

  return {
    // State
    stats,
    activities,
    alerts,
    traffic,
    loading,
    error,
    lastUpdate,
    // Getters
    hasAlerts,
    criticalAlerts,
    warningAlerts,
    // Actions
    fetchStats,
    fetchActivities,
    fetchAlerts,
    fetchTraffic,
    fetchAll,
    addActivity,
    addAlert,
    removeAlert,
    clearAlerts,
    updateStats,
    updateTrafficPoint
  }
})