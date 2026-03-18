import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { agentsApi } from '@/api'

export const useAgentsStore = defineStore('agents', () => {
  // State
  const agents = ref([])
  const currentAgent = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const onlineAgents = computed(() => agents.value.filter(a => a.status === 'online'))
  const offlineAgents = computed(() => agents.value.filter(a => a.status === 'offline'))
  const agentStats = computed(() => ({
    total: agents.value.length,
    online: onlineAgents.value.length,
    offline: offlineAgents.value.length
  }))

  // Actions
  async function fetchAgents(params = {}) {
    loading.value = true
    error.value = null

    try {
      const response = await agentsApi.list(params)
      agents.value = response.data.data || response.data || []
    } catch (e) {
      error.value = e.message || 'Failed to fetch agents'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchAgent(id) {
    loading.value = true
    error.value = null

    try {
      const response = await agentsApi.get(id)
      currentAgent.value = response.data
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch agent'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function createAgent(data) {
    loading.value = true
    error.value = null

    try {
      const response = await agentsApi.create(data)
      agents.value.unshift(response.data)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to create agent'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateAgent(id, data) {
    loading.value = true
    error.value = null

    try {
      const response = await agentsApi.update(id, data)
      const index = agents.value.findIndex(a => a.id === id)
      if (index !== -1) {
        agents.value[index] = response.data
      }
      if (currentAgent.value?.id === id) {
        currentAgent.value = response.data
      }
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to update agent'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteAgent(id) {
    loading.value = true
    error.value = null

    try {
      await agentsApi.delete(id)
      agents.value = agents.value.filter(a => a.id !== id)
      if (currentAgent.value?.id === id) {
        currentAgent.value = null
      }
    } catch (e) {
      error.value = e.message || 'Failed to delete agent'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function connectAgent(id) {
    try {
      const response = await agentsApi.connect(id)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to connect to agent'
      throw e
    }
  }

  async function disconnectAgent(id) {
    try {
      await agentsApi.disconnect(id)
    } catch (e) {
      error.value = e.message || 'Failed to disconnect from agent'
      throw e
    }
  }

  async function executeCommand(id, command) {
    try {
      const response = await agentsApi.execute(id, command)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to execute command'
      throw e
    }
  }

  async function fetchAgentStats(id) {
    try {
      const response = await agentsApi.stats(id)
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function restartAgent(id) {
    try {
      await agentsApi.restart(id)
      const agent = agents.value.find(a => a.id === id)
      if (agent) {
        agent.status = 'restarting'
      }
    } catch (e) {
      error.value = e.message || 'Failed to restart agent'
      throw e
    }
  }

  async function updateAgentFirmware(id) {
    try {
      await agentsApi.update(id)
    } catch (e) {
      error.value = e.message || 'Failed to update agent'
      throw e
    }
  }

  // Real-time updates
  function updateAgentStatus(id, status) {
    const agent = agents.value.find(a => a.id === id)
    if (agent) {
      agent.status = status
    }
  }

  function updateAgentStats(id, stats) {
    const agent = agents.value.find(a => a.id === id)
    if (agent) {
      Object.assign(agent, stats)
    }
  }

  function clearCurrentAgent() {
    currentAgent.value = null
  }

  return {
    // State
    agents,
    currentAgent,
    loading,
    error,
    // Getters
    onlineAgents,
    offlineAgents,
    agentStats,
    // Actions
    fetchAgents,
    fetchAgent,
    createAgent,
    updateAgent,
    deleteAgent,
    connectAgent,
    disconnectAgent,
    executeCommand,
    fetchAgentStats,
    restartAgent,
    updateAgentFirmware,
    updateAgentStatus,
    updateAgentStats,
    clearCurrentAgent
  }
})