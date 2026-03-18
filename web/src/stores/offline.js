import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface OfflineData {
  nodes: any[]
  users: any[]
  plugins: any[]
  settings: any
  lastSync: string
  version: number
}

const DB_NAME = 'anixops-offline'
const DB_VERSION = 1
const STORE_NAME = 'data'

// IndexedDB helper
function openDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const request = indexedDB.open(DB_NAME, DB_VERSION)

    request.onerror = () => reject(request.error)
    request.onsuccess = () => resolve(request.result)

    request.onupgradeneeded = (event) => {
      const db = (event.target as IDBOpenDBRequest).result
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME)
      }
    }
  })
}

async function saveToIndexedDB(key: string, data: any): Promise<void> {
  const db = await openDB()
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readwrite')
    const store = transaction.objectStore(STORE_NAME)
    const request = store.put(data, key)

    request.onerror = () => reject(request.error)
    request.onsuccess = () => resolve()
  })
}

async function getFromIndexedDB(key: string): Promise<any> {
  const db = await openDB()
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readonly')
    const store = transaction.objectStore(STORE_NAME)
    const request = store.get(key)

    request.onerror = () => reject(request.error)
    request.onsuccess = () => resolve(request.result)
  })
}

async function clearIndexedDB(): Promise<void> {
  const db = await openDB()
  return new Promise((resolve, reject) => {
    const transaction = db.transaction(STORE_NAME, 'readwrite')
    const store = transaction.objectStore(STORE_NAME)
    const request = store.clear()

    request.onerror = () => reject(request.error)
    request.onsuccess = () => resolve()
  })
}

export const useOfflineStore = defineStore('offline', () => {
  const isOnline = ref(navigator.onLine)
  const isOfflineMode = ref(false)
  const lastSyncTime = ref<string | null>(null)
  const pendingActions = ref<any[]>([])
  const syncInProgress = ref(false)

  // Computed
  const hasPendingActions = computed(() => pendingActions.value.length > 0)
  const pendingCount = computed(() => pendingActions.value.length)

  // Initialize
  async function initialize() {
    // Load pending actions from storage
    const stored = await getFromIndexedDB('pendingActions')
    if (stored) {
      pendingActions.value = stored
    }

    const lastSync = await getFromIndexedDB('lastSync')
    if (lastSync) {
      lastSyncTime.value = lastSync
    }

    // Listen for online/offline events
    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
  }

  function handleOnline() {
    isOnline.value = true
    if (hasPendingActions.value) {
      syncPendingActions()
    }
  }

  function handleOffline() {
    isOnline.value = false
  }

  // Save data for offline use
  async function saveOfflineData(data: Partial<OfflineData>) {
    const currentData = await getFromIndexedDB('offlineData') || {}
    const newData: OfflineData = {
      nodes: data.nodes || currentData.nodes || [],
      users: data.users || currentData.users || [],
      plugins: data.plugins || currentData.plugins || [],
      settings: data.settings || currentData.settings || {},
      lastSync: new Date().toISOString(),
      version: 1
    }

    await saveToIndexedDB('offlineData', newData)
    lastSyncTime.value = newData.lastSync
    await saveToIndexedDB('lastSync', newData.lastSync)
  }

  // Get offline data
  async function getOfflineData(): Promise<OfflineData | null> {
    return getFromIndexedDB('offlineData')
  }

  // Add pending action
  async function addPendingAction(action: {
    type: string
    endpoint: string
    method: string
    data?: any
    timestamp: string
  }) {
    pendingActions.value.push({
      ...action,
      id: Date.now().toString()
    })
    await saveToIndexedDB('pendingActions', pendingActions.value)
  }

  // Sync pending actions
  async function syncPendingActions(): Promise<{ success: number; failed: number }> {
    if (syncInProgress.value || !isOnline.value) {
      return { success: 0, failed: 0 }
    }

    syncInProgress.value = true
    let success = 0
    let failed = 0

    const actionsToSync = [...pendingActions.value]

    for (const action of actionsToSync) {
      try {
        // Execute the action
        const response = await fetch(action.endpoint, {
          method: action.method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: action.data ? JSON.stringify(action.data) : undefined
        })

        if (response.ok) {
          // Remove from pending
          const index = pendingActions.value.findIndex(a => a.id === action.id)
          if (index !== -1) {
            pendingActions.value.splice(index, 1)
          }
          success++
        } else {
          failed++
        }
      } catch {
        failed++
      }
    }

    await saveToIndexedDB('pendingActions', pendingActions.value)
    syncInProgress.value = false

    return { success, failed }
  }

  // Clear all offline data
  async function clearOfflineData() {
    await clearIndexedDB()
    pendingActions.value = []
    lastSyncTime.value = null
  }

  // Enable/disable offline mode
  function setOfflineMode(enabled: boolean) {
    isOfflineMode.value = enabled
  }

  return {
    isOnline,
    isOfflineMode,
    lastSyncTime,
    pendingActions,
    hasPendingActions,
    pendingCount,
    syncInProgress,
    initialize,
    saveOfflineData,
    getOfflineData,
    addPendingAction,
    syncPendingActions,
    clearOfflineData,
    setOfflineMode
  }
})