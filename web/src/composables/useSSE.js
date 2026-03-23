import { ref, onUnmounted } from 'vue'

/**
 * SSE (Server-Sent Events) composable for real-time communication
 * with the Workers API
 */
export function useSSE() {
  const eventSource = ref(null)
  const connected = ref(false)
  const error = ref(null)
  const handlers = new Map()

  let reconnectTimer = null
  let reconnectAttempts = 0
  const maxReconnectAttempts = 10
  const baseReconnectDelay = 1000
  const maxReconnectDelay = 30000

  /**
   * Connect to SSE endpoint
   * @param {string} url - SSE endpoint URL
   * @param {string} token - Bearer token for authentication
   */
  function connect(url, token) {
    if (eventSource.value) {
      disconnect()
    }

    // EventSource doesn't support headers, so we need to pass token in URL
    const urlWithToken = token ? `${url}?token=${encodeURIComponent(token)}` : url

    try {
      eventSource.value = new EventSource(urlWithToken)
      error.value = null

      eventSource.value.onopen = () => {
        connected.value = true
        reconnectAttempts = 0
        console.log('SSE connected')
      }

      eventSource.value.onerror = (e) => {
        connected.value = false
        error.value = 'Connection error'

        // Attempt reconnection
        if (reconnectAttempts < maxReconnectAttempts) {
          scheduleReconnect(url, token)
        }
      }

      eventSource.value.onmessage = (e) => {
        handleMessage(e)
      }
    } catch (e) {
      error.value = e.message
      scheduleReconnect(url, token)
    }
  }

  /**
   * Handle incoming SSE message
   */
  function handleMessage(event) {
    try {
      const data = JSON.parse(event.data)

      // Handle Workers API format: { type: 'event_type', payload: {...}, timestamp: '...' }
      if (data.type) {
        dispatch(data.type, data.payload || data)
      }

      // Also dispatch to 'message' handlers
      dispatch('message', data)
    } catch (e) {
      // Plain text message
      dispatch('message', event.data)
    }
  }

  /**
   * Dispatch event to registered handlers
   */
  function dispatch(eventType, data) {
    if (handlers.has(eventType)) {
      for (const handler of handlers.get(eventType)) {
        try {
          handler(data)
        } catch (e) {
          console.error(`Error in SSE handler for ${eventType}:`, e)
        }
      }
    }
  }

  /**
   * Schedule reconnection with exponential backoff
   */
  function scheduleReconnect(url, token) {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
    }

    reconnectAttempts++
    const delay = Math.min(
      baseReconnectDelay * Math.pow(2, reconnectAttempts - 1),
      maxReconnectDelay
    )

    console.log(`SSE reconnecting in ${delay}ms (attempt ${reconnectAttempts})`)

    reconnectTimer = setTimeout(() => {
      connect(url, token)
    }, delay)
  }

  /**
   * Disconnect from SSE endpoint
   */
  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }

    if (eventSource.value) {
      eventSource.value.close()
      eventSource.value = null
    }

    connected.value = false
    handlers.clear()
  }

  /**
   * Register handler for event type
   * @param {string} eventType - Event type to listen for
   * @param {function} handler - Handler function
   */
  function on(eventType, handler) {
    if (!handlers.has(eventType)) {
      handlers.set(eventType, [])
    }
    handlers.get(eventType).push(handler)
  }

  /**
   * Remove handler for event type
   * @param {string} eventType - Event type
   * @param {function} handler - Handler to remove (optional, removes all if not provided)
   */
  function off(eventType, handler) {
    if (!handler) {
      handlers.delete(eventType)
    } else if (handlers.has(eventType)) {
      const eventHandlers = handlers.get(eventType)
      const index = eventHandlers.indexOf(handler)
      if (index > -1) {
        eventHandlers.splice(index, 1)
      }
    }
  }

  /**
   * Subscribe to a channel via REST API
   * @param {string} channel - Channel name
   * @param {string} token - Auth token
   * @param {string} baseUrl - API base URL
   */
  async function subscribe(channel, token, baseUrl = '/api/v1') {
    try {
      const response = await fetch(`${baseUrl}/sse/subscribe`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ channel })
      })
      return response.ok
    } catch (e) {
      console.error('Failed to subscribe to channel:', e)
      return false
    }
  }

  /**
   * Unsubscribe from a channel
   * @param {string} channel - Channel name
   * @param {string} token - Auth token
   * @param {string} baseUrl - API base URL
   */
  async function unsubscribe(channel, token, baseUrl = '/api/v1') {
    try {
      const response = await fetch(`${baseUrl}/sse/unsubscribe`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ channel })
      })
      return response.ok
    } catch (e) {
      console.error('Failed to unsubscribe from channel:', e)
      return false
    }
  }

  // Cleanup on unmount
  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    error,
    connect,
    disconnect,
    on,
    off,
    subscribe,
    unsubscribe
  }
}

/**
 * SSE event types from Workers API
 */
export const SSEEventTypes = {
  CONNECTED: 'connected',
  NODE_UPDATE: 'node_update',
  TASK_UPDATE: 'task_update',
  LOG: 'log',
  NOTIFICATION: 'notification'
}

/**
 * SSE channels
 */
export const SSEChannels = {
  GLOBAL: 'global',
  NODES: 'nodes',
  TASKS: 'tasks',
  LOGS: 'logs',
  USER: (userId) => `user:${userId}`,
  NODE: (nodeId) => `node:${nodeId}`,
  TASK: (taskId) => `task:${taskId}`
}