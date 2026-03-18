import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useNodesStore } from '@/stores/nodes'

// Mock axios
vi.mock('axios', () => ({
  default: {
    create: () => ({
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
      interceptors: {
        request: { use: vi.fn() },
        response: { use: vi.fn() }
      }
    })
  }
}))

describe('Nodes Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('has correct initial state', () => {
    const store = useNodesStore()

    expect(store.nodes).toEqual([])
    expect(store.currentNode).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
    expect(store.total).toBe(0)
  })

  it('computes onlineNodes correctly', () => {
    const store = useNodesStore()

    store.nodes = [
      { id: '1', status: 'online' },
      { id: '2', status: 'offline' },
      { id: '3', status: 'online' }
    ]

    expect(store.onlineNodes).toHaveLength(2)
  })

  it('computes offlineNodes correctly', () => {
    const store = useNodesStore()

    store.nodes = [
      { id: '1', status: 'online' },
      { id: '2', status: 'offline' },
      { id: '3', status: 'online' }
    ]

    expect(store.offlineNodes).toHaveLength(1)
  })

  it('computes nodeStats correctly', () => {
    const store = useNodesStore()

    store.nodes = [
      { id: '1', status: 'online' },
      { id: '2', status: 'offline' },
      { id: '3', status: 'online' }
    ]

    expect(store.nodeStats).toEqual({
      total: 3,
      online: 2,
      offline: 1
    })
  })

  it('sets filters correctly', () => {
    const store = useNodesStore()

    store.setFilters({ search: 'test', status: 'online' })

    expect(store.filters.search).toBe('test')
    expect(store.filters.status).toBe('online')
  })

  it('clears current node', () => {
    const store = useNodesStore()

    store.currentNode = { id: '1', name: 'Test' }
    store.clearCurrentNode()

    expect(store.currentNode).toBeNull()
  })

  it('updates node status', () => {
    const store = useNodesStore()

    store.nodes = [
      { id: '1', status: 'online' },
      { id: '2', status: 'offline' }
    ]

    store.updateNodeStatus('1', 'offline')

    expect(store.nodes[0].status).toBe('offline')
  })

  it('updates current node status', () => {
    const store = useNodesStore()

    store.currentNode = { id: '1', status: 'online' }
    store.updateNodeStatus('1', 'offline')

    expect(store.currentNode.status).toBe('offline')
  })
})