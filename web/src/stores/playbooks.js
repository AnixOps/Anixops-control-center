import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { playbooksApi } from '@/api'

export const usePlaybooksStore = defineStore('playbooks', () => {
  // State
  const playbooks = ref([])
  const currentPlaybook = ref(null)
  const runningPlaybooks = ref(new Map())
  const loading = ref(false)
  const error = ref(null)
  const total = ref(0)

  // Getters
  const playbookStats = computed(() => ({
    total: playbooks.value.length,
    running: runningPlaybooks.value.size,
    success: playbooks.value.filter(p => p.lastStatus === 'success').length,
    failed: playbooks.value.filter(p => p.lastStatus === 'failed').length
  }))

  // Actions
  async function fetchPlaybooks(params = {}) {
    loading.value = true
    error.value = null

    try {
      const response = await playbooksApi.list(params)
      playbooks.value = response.data.data || response.data || []
      total.value = response.data.total || playbooks.value.length
    } catch (e) {
      error.value = e.message || 'Failed to fetch playbooks'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchPlaybook(id) {
    loading.value = true
    error.value = null

    try {
      const response = await playbooksApi.get(id)
      currentPlaybook.value = response.data
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch playbook'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function createPlaybook(data) {
    loading.value = true
    error.value = null

    try {
      const response = await playbooksApi.create(data)
      playbooks.value.unshift(response.data)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to create playbook'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updatePlaybook(id, data) {
    loading.value = true
    error.value = null

    try {
      const response = await playbooksApi.update(id, data)
      const index = playbooks.value.findIndex(p => p.id === id)
      if (index !== -1) {
        playbooks.value[index] = response.data
      }
      if (currentPlaybook.value?.id === id) {
        currentPlaybook.value = response.data
      }
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to update playbook'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deletePlaybook(id) {
    loading.value = true
    error.value = null

    try {
      await playbooksApi.delete(id)
      playbooks.value = playbooks.value.filter(p => p.id !== id)
      if (currentPlaybook.value?.id === id) {
        currentPlaybook.value = null
      }
    } catch (e) {
      error.value = e.message || 'Failed to delete playbook'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function runPlaybook(id, params = {}) {
    error.value = null

    try {
      const response = await playbooksApi.run(id, params)
      runningPlaybooks.value.set(id, {
        id,
        status: 'running',
        progress: 0,
        logs: [],
        startTime: new Date()
      })
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to run playbook'
      throw e
    }
  }

  async function stopPlaybook(id) {
    error.value = null

    try {
      await playbooksApi.stop(id)
      runningPlaybooks.value.delete(id)
    } catch (e) {
      error.value = e.message || 'Failed to stop playbook'
      throw e
    }
  }

  async function duplicatePlaybook(id) {
    error.value = null

    try {
      const response = await playbooksApi.duplicate(id)
      playbooks.value.unshift(response.data)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to duplicate playbook'
      throw e
    }
  }

  async function fetchTemplates() {
    try {
      const response = await playbooksApi.templates()
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function validatePlaybook(data) {
    try {
      const response = await playbooksApi.validate(data)
      return { valid: true, errors: [] }
    } catch (e) {
      return { valid: false, errors: e.errors || [e.message] }
    }
  }

  function updateRunProgress(id, progress) {
    const run = runningPlaybooks.value.get(id)
    if (run) {
      run.progress = progress
    }
  }

  function addRunLog(id, log) {
    const run = runningPlaybooks.value.get(id)
    if (run) {
      run.logs.push(log)
    }
  }

  function completeRun(id, status) {
    const run = runningPlaybooks.value.get(id)
    if (run) {
      run.status = status
      run.endTime = new Date()
    }

    // Update playbook last status
    const playbook = playbooks.value.find(p => p.id === id)
    if (playbook) {
      playbook.lastStatus = status
      playbook.lastRun = new Date().toISOString()
    }
  }

  function clearCurrentPlaybook() {
    currentPlaybook.value = null
  }

  return {
    // State
    playbooks,
    currentPlaybook,
    runningPlaybooks,
    loading,
    error,
    total,
    // Getters
    playbookStats,
    // Actions
    fetchPlaybooks,
    fetchPlaybook,
    createPlaybook,
    updatePlaybook,
    deletePlaybook,
    runPlaybook,
    stopPlaybook,
    duplicatePlaybook,
    fetchTemplates,
    validatePlaybook,
    updateRunProgress,
    addRunLog,
    completeRun,
    clearCurrentPlaybook
  }
})