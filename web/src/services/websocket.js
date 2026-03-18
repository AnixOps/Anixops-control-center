import { ref, onUnmounted } from 'vue'

class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 1000
    this.listeners = new Map()
    this.connected = ref(false)
    this.connecting = ref(false)
  }

  connect(url) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return
    }

    this.connecting.value = true
    this.ws = new WebSocket(url)

    this.ws.onopen = () => {
      this.connected.value = true
      this.connecting.value = false
      this.reconnectAttempts = 0
      console.log('[WebSocket] Connected')
      this.emit('connected')
    }

    this.ws.onclose = (event) => {
      this.connected.value = false
      this.connecting.value = false
      console.log('[WebSocket] Disconnected:', event.code, event.reason)
      this.emit('disconnected', { code: event.code, reason: event.reason })

      // Attempt reconnect
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++
        const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1)
        console.log(`[WebSocket] Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`)
        setTimeout(() => this.connect(url), delay)
      }
    }

    this.ws.onerror = (error) => {
      console.error('[WebSocket] Error:', error)
      this.emit('error', error)
    }

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        this.handleMessage(data)
      } catch (e) {
        console.error('[WebSocket] Failed to parse message:', event.data)
      }
    }
  }

  handleMessage(data) {
    const { type, payload } = data

    // Emit to specific event type
    if (type && this.listeners.has(type)) {
      this.listeners.get(type).forEach(callback => callback(payload))
    }

    // Emit to all listeners
    if (this.listeners.has('*')) {
      this.listeners.get('*').forEach(callback => callback(data))
    }
  }

  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, new Set())
    }
    this.listeners.get(event).add(callback)

    // Return unsubscribe function
    return () => this.off(event, callback)
  }

  off(event, callback) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).delete(callback)
    }
  }

  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => callback(data))
    }
  }

  send(type, payload) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, payload }))
      return true
    }
    return false
  }

  disconnect() {
    if (this.ws) {
      this.reconnectAttempts = this.maxReconnectAttempts // Prevent reconnect
      this.ws.close()
      this.ws = null
    }
    this.connected.value = false
  }

  isConnected() {
    return this.connected.value
  }
}

// Singleton instance
const wsService = new WebSocketService()

// Composable for Vue components
export function useWebSocket() {
  const subscribe = (event, callback) => {
    const unsubscribe = wsService.on(event, callback)

    onUnmounted(() => {
      unsubscribe()
    })

    return unsubscribe
  }

  const send = (type, payload) => wsService.send(type, payload)
  const connect = (url) => wsService.connect(url)
  const disconnect = () => wsService.disconnect()

  return {
    connected: wsService.connected,
    connecting: wsService.connecting,
    subscribe,
    send,
    connect,
    disconnect
  }
}

// Event types
export const WSEvents = {
  // Connection
  CONNECTED: 'connected',
  DISCONNECTED: 'disconnected',
  ERROR: 'error',

  // System
  SYSTEM_STATUS: 'system:status',
  SYSTEM_METRICS: 'system:metrics',

  // Nodes
  NODE_CREATED: 'node:created',
  NODE_UPDATED: 'node:updated',
  NODE_DELETED: 'node:deleted',
  NODE_STATUS: 'node:status',
  NODE_STATS: 'node:stats',

  // Plugins
  PLUGIN_STATUS: 'plugin:status',
  PLUGIN_EVENT: 'plugin:event',

  // Users
  USER_CREATED: 'user:created',
  USER_UPDATED: 'user:updated',
  USER_DELETED: 'user.deleted',
  USER_LOGIN: 'user:login',
  USER_LOGOUT: 'user:logout',

  // Agents
  AGENT_CONNECTED: 'agent:connected',
  AGENT_DISCONNECTED: 'agent:disconnected',
  AGENT_STATUS: 'agent:status',

  // Logs
  LOG_ENTRY: 'log:entry',

  // Alerts
  ALERT_CREATED: 'alert:created',
  ALERT_RESOLVED: 'alert:resolved',

  // Activity
  ACTIVITY: 'activity',

  // Traffic
  TRAFFIC_UPDATE: 'traffic:update'
}

export default wsService