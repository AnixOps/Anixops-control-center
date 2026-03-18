import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { pluginsApi } from '@/api'

export const usePluginsStore = defineStore('plugins', () => {
  // State
  const plugins = ref([])
  const currentPlugin = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const executing = ref(false)

  // Getters
  const activePlugins = computed(() => plugins.value.filter(p => p.status === 'running'))
  const inactivePlugins = computed(() => plugins.value.filter(p => p.status !== 'running'))
  const pluginStats = computed(() => ({
    total: plugins.value.length,
    active: activePlugins.value.length,
    inactive: inactivePlugins.value.length
  }))

  // Actions
  async function fetchPlugins() {
    loading.value = true
    error.value = null
    try {
      const response = await pluginsApi.list()
      plugins.value = response.data.data || response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch plugins'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchPlugin(name) {
    loading.value = true
    error.value = null
    try {
      const response = await pluginsApi.get(name)
      currentPlugin.value = response.data
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch plugin'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function executePlugin(name, action, params = {}) {
    executing.value = true
    error.value = null
    try {
      const response = await pluginsApi.execute(name, action, params)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to execute plugin'
      throw e
    } finally {
      executing.value = false
    }
  }

  async function fetchPluginStatus(name) {
    try {
      const response = await pluginsApi.status(name)
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function fetchPluginConfig(name) {
    try {
      const response = await pluginsApi.config(name)
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function updatePluginConfig(name, config) {
    loading.value = true
    error.value = null
    try {
      const response = await pluginsApi.updateConfig(name, config)
      if (currentPlugin.value?.name === name) {
        currentPlugin.value.config = config
      }
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to update plugin config'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function startPlugin(name) {
    try {
      const response = await pluginsApi.start(name)
      updatePluginStatus(name, 'running')
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function stopPlugin(name) {
    try {
      const response = await pluginsApi.stop(name)
      updatePluginStatus(name, 'stopped')
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function restartPlugin(name) {
    try {
      const response = await pluginsApi.restart(name)
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function enablePlugin(name) {
    try {
      const response = await pluginsApi.enable(name)
      updatePluginStatus(name, 'running')
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function disablePlugin(name) {
    try {
      const response = await pluginsApi.disable(name)
      updatePluginStatus(name, 'disabled')
      return response.data
    } catch (e) {
      throw e
    }
  }

  function updatePluginStatus(name, status) {
    const index = plugins.value.findIndex(p => p.name === name)
    if (index !== -1) {
      plugins.value[index].status = status
    }
    if (currentPlugin.value?.name === name) {
      currentPlugin.value.status = status
    }
  }

  function clearCurrentPlugin() {
    currentPlugin.value = null
  }

  return {
    // State
    plugins,
    currentPlugin,
    loading,
    error,
    executing,
    // Getters
    activePlugins,
    inactivePlugins,
    pluginStats,
    // Actions
    fetchPlugins,
    fetchPlugin,
    executePlugin,
    fetchPluginStatus,
    fetchPluginConfig,
    updatePluginConfig,
    startPlugin,
    stopPlugin,
    restartPlugin,
    enablePlugin,
    disablePlugin,
    updatePluginStatus,
    clearCurrentPlugin
  }
})