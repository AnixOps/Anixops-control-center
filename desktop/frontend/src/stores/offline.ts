import { ref, watch } from 'vue'
import { wailsAPI, type Node } from '../api/wails'

// IndexedDB wrapper for offline storage
class OfflineDB {
  private dbName = 'anixops-offline'
  private dbVersion = 1
  private db: IDBDatabase | null = null

  async init(): Promise<void> {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(this.dbName, this.dbVersion)

      request.onerror = () => reject(request.error)

      request.onsuccess = () => {
        this.db = request.result
        resolve()
      }

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result

        // Create object stores
        if (!db.objectStoreNames.contains('nodes')) {
          db.createObjectStore('nodes', { keyPath: 'id' })
        }
        if (!db.objectStoreNames.contains('users')) {
          db.createObjectStore('users', { keyPath: 'id' })
        }
        if (!db.objectStoreNames.contains('plugins')) {
          db.createObjectStore('plugins', { keyPath: 'name' })
        }
        if (!db.objectStoreNames.contains('settings')) {
          db.createObjectStore('settings', { keyPath: 'key' })
        }
        if (!db.objectStoreNames.contains('pendingActions')) {
          const store = db.createObjectStore('pendingActions', {
            keyPath: 'id',
            autoIncrement: true
          })
          store.createIndex('timestamp', 'timestamp', { unique: false })
        }
      }
    })
  }

  async get<T>(storeName: string, key: string): Promise<T | null> {
    if (!this.db) await this.init()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(storeName, 'readonly')
      const store = transaction.objectStore(storeName)
      const request = store.get(key)

      request.onsuccess = () => resolve(request.result || null)
      request.onerror = () => reject(request.error)
    })
  }

  async getAll<T>(storeName: string): Promise<T[]> {
    if (!this.db) await this.init()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(storeName, 'readonly')
      const store = transaction.objectStore(storeName)
      const request = store.getAll()

      request.onsuccess = () => resolve(request.result || [])
      request.onerror = () => reject(request.error)
    })
  }

  async put<T>(storeName: string, data: T): Promise<void> {
    if (!this.db) await this.init()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(storeName, 'readwrite')
      const store = transaction.objectStore(storeName)
      const request = store.put(data)

      request.onsuccess = () => resolve()
      request.onerror = () => reject(request.error)
    })
  }

  async putAll<T>(storeName: string, items: T[]): Promise<void> {
    if (!this.db) await this.init()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(storeName, 'readwrite')
      const store = transaction.objectStore(storeName)

      for (const item of items) {
        store.put(item)
      }

      transaction.oncomplete = () => resolve()
      transaction.onerror = () => reject(transaction.error)
    })
  }

  async delete(storeName: string, key: string): Promise<void> {
    if (!this.db) await this.init()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(storeName, 'readwrite')
      const store = transaction.objectStore(storeName)
      const request = store.delete(key)

      request.onsuccess = () => resolve()
      request.onerror = () => reject(request.error)
    })
  }

  async clear(storeName: string): Promise<void> {
    if (!this.db) await this.init()

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction(storeName, 'readwrite')
      const store = transaction.objectStore(storeName)
      const request = store.clear()

      request.onsuccess = () => resolve()
      request.onerror = () => reject(request.error)
    })
  }
}

// Pending action for offline operations
interface PendingAction {
  id?: number
  type: 'create' | 'update' | 'delete'
  entity: 'node' | 'user' | 'plugin'
  data: any
  timestamp: number
}

// Offline store
export const useOfflineStore = () => {
  const db = new OfflineDB()
  const isOnline = ref(navigator.onLine)
  const isInitialized = ref(false)
  const pendingActionsCount = ref(0)
  const lastSyncTime = ref<Date | null>(null)

  // Initialize offline storage
  async function initialize() {
    try {
      await db.init()
      await loadPendingActionsCount()
      await loadLastSyncTime()
      isInitialized.value = true

      // Setup online/offline listeners
      window.addEventListener('online', handleOnline)
      window.addEventListener('offline', handleOffline)
    } catch (error) {
      console.error('Failed to initialize offline storage:', error)
    }
  }

  // Handle online event
  async function handleOnline() {
    isOnline.value = true
    await syncPendingActions()
  }

  // Handle offline event
  function handleOffline() {
    isOnline.value = false
  }

  // Load pending actions count
  async function loadPendingActionsCount() {
    const actions = await db.getAll<PendingAction>('pendingActions')
    pendingActionsCount.value = actions.length
  }

  // Load last sync time
  async function loadLastSyncTime() {
    const record = await db.get<{ key: string; value: string }>('settings', 'lastSync')
    if (record) {
      lastSyncTime.value = new Date(record.value)
    }
  }

  // Save last sync time
  async function saveLastSyncTime() {
    lastSyncTime.value = new Date()
    await db.put('settings', {
      key: 'lastSync',
      value: lastSyncTime.value.toISOString()
    })
  }

  // Cache nodes for offline use
  async function cacheNodes(nodes: Node[]) {
    await db.clear('nodes')
    await db.putAll('nodes', nodes)
  }

  // Get cached nodes
  async function getCachedNodes(): Promise<Node[]> {
    return db.getAll<Node>('nodes')
  }

  // Cache a single node
  async function cacheNode(node: Node) {
    await db.put('nodes', node)
  }

  // Remove cached node
  async function removeCachedNode(id: string) {
    await db.delete('nodes', id)
  }

  // Add pending action
  async function addPendingAction(
    type: PendingAction['type'],
    entity: PendingAction['entity'],
    data: any
  ) {
    const action: PendingAction = {
      type,
      entity,
      data,
      timestamp: Date.now()
    }
    await db.put('pendingActions', action)
    await loadPendingActionsCount()
  }

  // Sync pending actions
  async function syncPendingActions(): Promise<{ success: number; failed: number }> {
    if (isOnline.value) {
      return { success: 0, failed: 0 }
    }

    const actions = await db.getAll<PendingAction>('pendingActions')
    let success = 0
    let failed = 0

    for (const action of actions) {
      try {
        // Attempt to execute the action
        // This would be implemented based on the actual API
        success++

        // Remove successful action
        if (action.id) {
          await db.delete('pendingActions', action.id.toString())
        }
      } catch (error) {
        failed++
        console.error('Failed to sync action:', action, error)
      }
    }

    await loadPendingActionsCount()
    await saveLastSyncTime()

    return { success, failed }
  }

  // Clear all cached data
  async function clearCache() {
    await db.clear('nodes')
    await db.clear('users')
    await db.clear('plugins')
    await db.clear('pendingActions')
    pendingActionsCount.value = 0
  }

  // Cleanup
  function cleanup() {
    window.removeEventListener('online', handleOnline)
    window.removeEventListener('offline', handleOffline)
  }

  return {
    isOnline,
    isInitialized,
    pendingActionsCount,
    lastSyncTime,
    initialize,
    cacheNodes,
    getCachedNodes,
    cacheNode,
    removeCachedNode,
    addPendingAction,
    syncPendingActions,
    clearCache,
    cleanup
  }
}