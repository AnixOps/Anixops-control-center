import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { nodesApi } from '@/api'

export const useNodesStore = defineStore('nodes', () => {
  // State
  const nodes = ref([])
  const currentNode = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const total = ref(0)
  const filters = ref({
    search: '',
    status: '',
    page: 1,
    limit: 20
  })

  // Getters
  const onlineNodes = computed(() => nodes.value.filter(n => n.status === 'online'))
  const offlineNodes = computed(() => nodes.value.filter(n => n.status === 'offline'))
  const nodeStats = computed(() => ({
    total: nodes.value.length,
    online: onlineNodes.value.length,
    offline: offlineNodes.value.length
  }))

  // Actions
  async function fetchNodes(params = {}) {
    loading.value = true
    error.value = null
    try {
      const response = await nodesApi.list({ ...filters.value, ...params })
      nodes.value = response.data.data || response.data
      total.value = response.data.total || nodes.value.length
    } catch (e) {
      error.value = e.message || 'Failed to fetch nodes'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchNode(id) {
    loading.value = true
    error.value = null
    try {
      const response = await nodesApi.get(id)
      currentNode.value = response.data
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to fetch node'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function createNode(data) {
    loading.value = true
    error.value = null
    try {
      const response = await nodesApi.create(data)
      nodes.value.unshift(response.data)
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to create node'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateNode(id, data) {
    loading.value = true
    error.value = null
    try {
      const response = await nodesApi.update(id, data)
      const index = nodes.value.findIndex(n => n.id === id)
      if (index !== -1) {
        nodes.value[index] = response.data
      }
      if (currentNode.value?.id === id) {
        currentNode.value = response.data
      }
      return response.data
    } catch (e) {
      error.value = e.message || 'Failed to update node'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteNode(id) {
    loading.value = true
    error.value = null
    try {
      await nodesApi.delete(id)
      nodes.value = nodes.value.filter(n => n.id !== id)
      if (currentNode.value?.id === id) {
        currentNode.value = null
      }
    } catch (e) {
      error.value = e.message || 'Failed to delete node'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function startNode(id) {
    try {
      const response = await nodesApi.start(id)
      const index = nodes.value.findIndex(n => n.id === id)
      if (index !== -1) {
        nodes.value[index].status = 'starting'
      }
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function stopNode(id) {
    try {
      const response = await nodesApi.stop(id)
      const index = nodes.value.findIndex(n => n.id === id)
      if (index !== -1) {
        nodes.value[index].status = 'stopping'
      }
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function restartNode(id) {
    try {
      const response = await nodesApi.restart(id)
      return response.data
    } catch (e) {
      throw e
    }
  }

  async function fetchNodeStats(id) {
    try {
      const response = await nodesApi.stats(id)
      return response.data
    } catch (e) {
      throw e
    }
  }

  function setFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function clearCurrentNode() {
    currentNode.value = null
  }

  // Real-time updates
  function updateNodeStatus(id, status) {
    const index = nodes.value.findIndex(n => n.id === id)
    if (index !== -1) {
      nodes.value[index].status = status
    }
    if (currentNode.value?.id === id) {
      currentNode.value.status = status
    }
  }

  return {
    // State
    nodes,
    currentNode,
    loading,
    error,
    total,
    filters,
    // Getters
    onlineNodes,
    offlineNodes,
    nodeStats,
    // Actions
    fetchNodes,
    fetchNode,
    createNode,
    updateNode,
    deleteNode,
    startNode,
    stopNode,
    restartNode,
    fetchNodeStats,
    setFilters,
    clearCurrentNode,
    updateNodeStatus
  }
})