import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAppStore } from './app'

export interface Node {
  id: string
  name: string
  host: string
  port: number
  status: 'online' | 'offline' | 'starting' | 'stopping'
  type: string
  users: number
  traffic: number
  last_seen?: string
  cpu_usage?: number
  memory_usage?: number
  version?: string
}

export const useNodesStore = defineStore('nodes', () => {
  const appStore = useAppStore()

  // State
  const nodes = ref<Node[]>([])
  const currentNode = ref<Node | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const total = ref(0)
  const filters = ref({
    search: '',
    status: '',
    type: '',
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
  async function fetchNodes(params: Partial<typeof filters.value> = {}) {
    loading.value = true
    error.value = null

    try {
      const queryParams = new URLSearchParams()
      Object.entries({ ...filters.value, ...params }).forEach(([key, value]) => {
        if (value) queryParams.append(key, String(value))
      })

      const response = await fetch(`${appStore.apiUrl}/nodes?${queryParams}`, {
        headers: {
          'Authorization': `Bearer ${appStore.token}`
        }
      })

      if (!response.ok) throw new Error('Failed to fetch nodes')

      const data = await response.json()
      nodes.value = data.data || data
      total.value = data.total || nodes.value.length
    } catch (e) {
      error.value = String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchNode(id: string) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${appStore.apiUrl}/nodes/${id}`, {
        headers: {
          'Authorization': `Bearer ${appStore.token}`
        }
      })

      if (!response.ok) throw new Error('Failed to fetch node')

      currentNode.value = await response.json()
      return currentNode.value
    } catch (e) {
      error.value = String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function createNode(data: Partial<Node>) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${appStore.apiUrl}/nodes`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${appStore.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      })

      if (!response.ok) throw new Error('Failed to create node')

      const newNode = await response.json()
      nodes.value.unshift(newNode)
      return newNode
    } catch (e) {
      error.value = String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateNode(id: string, data: Partial<Node>) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${appStore.apiUrl}/nodes/${id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${appStore.token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
      })

      if (!response.ok) throw new Error('Failed to update node')

      const updatedNode = await response.json()
      const index = nodes.value.findIndex(n => n.id === id)
      if (index !== -1) {
        nodes.value[index] = updatedNode
      }
      if (currentNode.value?.id === id) {
        currentNode.value = updatedNode
      }
      return updatedNode
    } catch (e) {
      error.value = String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteNode(id: string) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${appStore.apiUrl}/nodes/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${appStore.token}`
        }
      })

      if (!response.ok) throw new Error('Failed to delete node')

      nodes.value = nodes.value.filter(n => n.id !== id)
      if (currentNode.value?.id === id) {
        currentNode.value = null
      }
    } catch (e) {
      error.value = String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function startNode(id: string) {
    try {
      const response = await fetch(`${appStore.apiUrl}/nodes/${id}/start`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${appStore.token}`
        }
      })

      if (!response.ok) throw new Error('Failed to start node')

      updateNodeStatus(id, 'starting')
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function stopNode(id: string) {
    try {
      const response = await fetch(`${appStore.apiUrl}/nodes/${id}/stop`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${appStore.token}`
        }
      })

      if (!response.ok) throw new Error('Failed to stop node')

      updateNodeStatus(id, 'stopping')
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  async function restartNode(id: string) {
    try {
      const response = await fetch(`${appStore.apiUrl}/nodes/${id}/restart`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${appStore.token}`
        }
      })

      if (!response.ok) throw new Error('Failed to restart node')

      updateNodeStatus(id, 'starting')
    } catch (e) {
      error.value = String(e)
      throw e
    }
  }

  function updateNodeStatus(id: string, status: Node['status']) {
    const node = nodes.value.find(n => n.id === id)
    if (node) {
      node.status = status
    }
    if (currentNode.value?.id === id) {
      currentNode.value.status = status
    }
  }

  function setFilters(newFilters: Partial<typeof filters.value>) {
    filters.value = { ...filters.value, ...newFilters }
  }

  function setNodes(newNodes: Node[]) {
    nodes.value = newNodes
  }

  function clearCurrentNode() {
    currentNode.value = null
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
    updateNodeStatus,
    setFilters,
    setNodes,
    clearCurrentNode
  }
})